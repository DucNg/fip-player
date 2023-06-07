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
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var width int

var program *tea.Program

type item struct {
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
	playingItemIndex int
	volume           float64
	trackName        string
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
	m.playingItemIndex = m.list.Index()
	m.volume = 100

	m.mpv.SendCommand([]string{"loadfile", m.list.Items()[m.playingItemIndex].(item).streamUrl})

	delayToRefresh := setMetadata(m)
	go UpdateMetadataLoop(m, delayToRefresh)

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "m":
			m.mpv.ToggleMute()
			return m, nil
		case "enter":
			if m.playingItemIndex == m.list.Index() {
				break
			}

			// Reset desc
			previousItem := m.list.Items()[m.playingItemIndex].(item)
			previousItem.desc = radios[m.playingItemIndex].(item).desc
			m.list.SetItem(m.playingItemIndex, previousItem)

			// Get new selection
			item := m.list.SelectedItem().(item)
			m.playingItemIndex = m.list.Index()

			// Change the mpv stream
			m.mpv.SendCommand([]string{"loadfile", item.streamUrl})

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
		m.list.SetSize(msg.Width-h, msg.Height-v-1)
	case descriptionUpdate:
		item := m.list.Items()[m.playingItemIndex].(item)
		item.desc = string(msg)
		cmd := m.list.SetItem(m.playingItemIndex, item)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	bar := topBar(
		m.list.Items()[m.playingItemIndex].(item).title,
		m.trackName,
		int(m.volume),
		m.mpv.IsMute(),
	)

	out := bar + "\n" + m.list.View()
	return docStyle.Render(out)
}

// Render creates the GUI and returns last selected radio index on close
func Render(ins *dbus.Instance, mpv *player.MPV, lastRadioIndex int) int {
	guiList := make([]list.Item, len(radios))
	copy(guiList, radios)

	m := model{
		list: list.New(guiList, list.NewDefaultDelegate(), 0, 0),
		mpv:  mpv,
		ins:  ins,
	}

	m.list.Select(lastRadioIndex)
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

	return m.playingItemIndex
}

func setMetadata(m *model) time.Duration {
	playingItem := m.list.Items()[m.playingItemIndex].(item)
	fm := metadata.FetchMetadata(playingItem.metadataUrl)
	dbus.UpdateMetadata(m.ins, fm)

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
	centerStr := header_center_s.Copy().
		Width(width - lipgloss.Width(statusStr) - lipgloss.Width(volumeStr)).
		Render(trackName)
	s := lipgloss.JoinHorizontal(lipgloss.Top, statusStr, centerStr, volumeStr)
	return s
}
