package loggy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

type Model struct {
	List          list.Model
	fetcher       Fetcher[LogEntries]
	offset, limit int
}
type LogEntryMsg struct {
	Entries LogEntries
	Error   error
}

type KeyMap struct {
	Help       key.Binding
	Pagination key.Binding
	Spinner    key.Binding
	TitleBar   key.Binding
	StatusBar  key.Binding
	NextPage   key.Binding
	PrevPage   key.Binding
}

var (
	keyMap = KeyMap{
		Help: key.NewBinding(
			key.WithKeys("H", "h"),
			key.WithHelp("h", "toggle help"),
		),
		Pagination: key.NewBinding(
			key.WithKeys("p", "P"),
			key.WithHelp("p", "toggle pagination"),
		),
		Spinner: key.NewBinding(
			key.WithKeys("s", "S"),
			key.WithHelp("s", "toggle pagination"),
		),
		TitleBar: key.NewBinding(
			key.WithKeys("t", "T"),
			key.WithHelp("t", "toggle title bar"),
		),
		StatusBar: key.NewBinding(
			key.WithKeys("b", "B"),
			key.WithHelp("b", "toggle status bar"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("n", "N"),
			key.WithHelp("n", "next page"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("p", "P"),
			key.WithHelp("p", "previous page"),
		),
	}
)

// New creates a new instance of the logger view.
func NewModel(fetcher Fetcher[LogEntries]) *Model {
	entries := LogEntries{}

	d := list.NewDefaultDelegate()

	l := list.New(entries.ToItems(), d, appStyle.GetWidth(), appStyle.GetHeight())
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.Help,
			keyMap.Pagination,
			keyMap.Spinner,
		}
	}
	return &Model{
		List:    l,
		fetcher: fetcher,
		offset:  0,
		limit:   100,
	}
}
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.List.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, keyMap.NextPage):
			m.offset += m.limit
			return m, m.fetchNewPage
		case key.Matches(msg, keyMap.PrevPage):
			if m.offset > 0 {
				m.offset -= m.limit
				if m.offset < 0 {
					m.offset = 0
				}
				return m, m.fetchNewPage
			}
		}

		newModel, cmd := m.List.Update(msg)
		m.List = newModel
		return m, cmd

	case LogEntryMsg:
		if msg.Error != nil {
			m.List.SetItems([]list.Item{}) // Clear list on error
			return m, nil
		}

		// Safe replacement of items based on current offset and limit
		currentItems := m.List.Items()
		newItems := msg.Entries.ToItems()

		// Calculate the replacement range
		startIndex := m.offset

		if startIndex > len(currentItems) {
			// If the startIndex is beyond the current list, pad the list with nil or placeholders
			for i := len(currentItems); i < startIndex; i++ {
				currentItems = append(currentItems, nil) // Append nil or a placeholder
			}
		}

		// Replace or append new items safely
		for i, item := range newItems {
			if startIndex+i < len(currentItems) {
				currentItems[startIndex+i] = item
			} else {
				currentItems = append(currentItems, item)
			}
		}

		m.List.SetItems(currentItems)
	}

	newList, cmd := m.List.Update(msg)
	m.List = newList
	return m, tea.Batch(cmd)
}

func (m *Model) fetchNewPage() tea.Msg {
	entries, err := m.fetcher(m.offset, m.limit)
	return LogEntryMsg{Entries: entries, Error: err}
}

func (m *Model) Init() tea.Cmd {
	return m.fetchNewPage
}

// View returns a string representation of the code bubble.
func (m *Model) View() string {
	return appStyle.Render(m.List.View())
}
