package gui

import (
	"log"
	"time"

	"github.com/DucNg/fip-player/dbus"
	"github.com/DucNg/fip-player/metadata"
	"github.com/DucNg/fip-player/player"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var radios = []list.Item{
	item{
		title:       "FIP",
		desc:        "I have ’em all over my house",
		streamUrl:   "https://stream.radiofrance.fr/fip/fip.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip",
	},
	item{
		title:       "FIP Jazz",
		desc:        "It's good on toast",
		streamUrl:   "https://stream.radiofrance.fr/fipjazz/fipjazz_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_jazz",
	},
}

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
}

func UpdateMetadataLoop(m *model, delayToRefresh time.Duration, url string) {
	m.metadataLoopChan = make(chan struct{})

	for {
		t := time.NewTimer(delayToRefresh)
		select {
		case <-t.C:
			delayToRefresh = setMetadata(m, url)
		case <-m.metadataLoopChan:
			return
		}
	}
}

func (m *model) Init() tea.Cmd {
	initialMetadataUrl := m.list.SelectedItem().(item).metadataUrl

	delayToRefresh := setMetadata(m, initialMetadataUrl)

	go UpdateMetadataLoop(m, delayToRefresh, initialMetadataUrl)

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			// Change the mpv stream
			item := m.list.SelectedItem().(item)
			m.mpv.SendCommand([]string{"loadfile", item.streamUrl})

			// Change the metadata loop url
			m.metadataLoopChan <- struct{}{}              // Stop the existing loop
			go UpdateMetadataLoop(m, 0, item.metadataUrl) // Start the new one

			// Change the description
			// TODO
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Render(ins *dbus.Instance, mpv *player.MPV) {
	m := model{
		list: list.New(radios, list.NewDefaultDelegate(), 0, 0),
		mpv:  mpv,
		ins:  ins,
	}
	m.list.Title = "FIP Radios"

	p := tea.NewProgram(&m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}

func setMetadata(m *model, url string) time.Duration {
	fm := metadata.FetchMetadata(url)

	dbus.UpdateMetadata(m.ins, fm)

	return fm.Delay()
}
