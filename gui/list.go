package gui

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DucNg/fip-player/dbus"
	"github.com/DucNg/fip-player/metadata"
	"github.com/DucNg/fip-player/player"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Usable width in terminal
var width int

const (
	topBarHeight      = 1
	progressBarHeight = 3
	ellipsis          = "…"
)

var program *tea.Program

type item struct {
	id                                  int
	title, desc, streamUrl, metadataUrl string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list             list.Model
	mpv              *player.MPV
	ins              *dbus.Instance
	metadataLoopChan chan struct{}
	playingItem      item
	volume           float64
	trackName        string
	ValueOfOneSecond float64
	progress         progress.Model
}

func UpdateMetadataLoop(m *model, delayToRefresh time.Duration) {
	m.metadataLoopChan = make(chan struct{})

	for {
		t := time.NewTimer(delayToRefresh)
		select {
		case <-t.C:
			delayToRefresh = setMetadata(m)
		case <-m.metadataLoopChan:
			return
		}
	}
}

func (m *model) Init() tea.Cmd {
	m.playingItem = m.list.SelectedItem().(item)
	m.volume = 100

	m.mpv.SendCommand([]string{"loadfile", m.playingItem.streamUrl})

	delayToRefresh := setMetadata(m)
	go UpdateMetadataLoop(m, delayToRefresh)

	return tickCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Disable keybinds when filtering
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "m":
			m.mpv.ToggleMute()
		case "enter":
			newSelectedItem := m.list.SelectedItem().(item)
			if m.playingItem.id == newSelectedItem.id {
				break
			}

			// Reset desc
			previousItem := m.playingItem
			previousItem.desc = radios[m.playingItem.id].(item).desc
			m.list.SetItem(m.getPlayingItemIndex(), previousItem)

			// Get new selection
			m.playingItem = m.list.SelectedItem().(item)

			// Change the mpv stream
			m.mpv.SendCommand([]string{"loadfile", m.playingItem.streamUrl})

			// Change the metadata loop url
			m.metadataLoopChan <- struct{}{} // Stop the existing loop
			go UpdateMetadataLoop(m, 0)      // Start the new one
		case "+", "-":
			if msg.String() == "+" && m.volume < 100 {
				m.volume += 5
			} else if msg.String() == "-" && m.volume > 0 {
				m.volume -= 5
			}
			m.mpv.SetVolume(m.volume)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()

		width = msg.Width - h
		height := msg.Height - v - topBarHeight - progressBarHeight

		m.list.SetSize(width, height)
		m.progress.Width = width
	case descriptionUpdate:
		m.playingItem.desc = string(msg)
		cmd := m.list.SetItem(m.getPlayingItemIndex(), m.playingItem)
		return m, cmd

	// Progress bar
	case tickMsg:
		// Sometimes delayToRefresh is not perfectly accurate
		// it's better to set progress bar to 0% in this case
		if m.progress.Percent() >= 1.0 {
			m.progress.SetPercent(0)
		}

		cmd := m.progress.IncrPercent(m.ValueOfOneSecond)
		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	bar := topBar(
		m.playingItem.title,
		m.trackName,
		int(m.volume),
		m.mpv.IsMute(),
	)

	out := bar + "\n" +
		m.list.View() + "\n\n" +
		m.progress.View()

	return docStyle.Render(out)
}

// Render creates the GUI and returns last selected radio index on close
func Render(ins *dbus.Instance, mpv *player.MPV, lastRadioID int) int {
	m := model{
		list:     list.New(getRadiosWithIDs(), list.NewDefaultDelegate(), 0, 0),
		mpv:      mpv,
		ins:      ins,
		progress: progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage()),
	}

	m.list.Select(getIndexFromID(m.list, lastRadioID))
	m.list.SetShowTitle(false)
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "volume up")),
			key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "volume down")),
			key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "toggle mute")),
		}
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())

	program = p

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-c
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}

	return m.playingItem.id
}

func setMetadata(m *model) time.Duration {
	fm := metadata.FetchMetadata(m.playingItem.metadataUrl)
	dbus.UpdateMetadata(m.ins, fm)

	m.ValueOfOneSecond = fm.ValueOfOneSecond()
	m.progress.SetPercent(fm.ProgressPercent())

	go program.Send(updateDesc(m, fm))

	return fm.Delay()
}

type descriptionUpdate string

func updateDesc(m *model, fm *metadata.FipMetadata) descriptionUpdate {
	m.trackName = fmt.Sprintf("▶ %v - %v", fm.Now.FirstLine, fm.Now.SecondLine)
	return descriptionUpdate(m.trackName)
}

// render topbar
func topBar(currentStation string, trackName string, volume int, muted bool) string {
	var mutedStr string
	if muted {
		mutedStr = fmt.Sprintf("Muted(%d%%)", volume)
	} else {
		mutedStr = fmt.Sprintf("Volume %d%%", volume)
	}
	statusStr := header_status_s.Render(currentStation)
	volumeStr := header_volume_s.Render(mutedStr)

	maxTrackNameWidth := width - lipgloss.Width(statusStr) - lipgloss.Width(volumeStr)
	trackName = truncate.StringWithTail(trackName, uint(maxTrackNameWidth), ellipsis)

	centerStr := header_center_s.Copy().
		Width(width - lipgloss.Width(statusStr) - lipgloss.Width(volumeStr)).
		Render(trackName)
	s := lipgloss.JoinHorizontal(lipgloss.Top, statusStr, centerStr, volumeStr)
	return s
}

func (m model) getPlayingItemIndex() int {
	return getIndexFromID(m.list, m.playingItem.id)
}

func getIndexFromID(list list.Model, id int) int {
	for i := 0; i < len(list.Items()); i++ {
		if list.Items()[i].(item).id == id {
			return i
		}
	}

	return -1
}

func getRadiosWithIDs() []list.Item {
	radiosWithIDs := make([]list.Item, len(radios))
	for i := 0; i < len(radios); i++ {
		itemWithID := radios[i].(item)
		itemWithID.id = i
		radiosWithIDs[i] = itemWithID
	}

	return radiosWithIDs
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
