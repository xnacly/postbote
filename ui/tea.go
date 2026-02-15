package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xnacly/postbote"
	"github.com/xnacly/postbote/mail"
)

type state struct {
	vi        *vi
	activeIdx uint
}

type model struct {
	state state
}

func newModel(folders []mail.Folder) model {
	return model{
		state: state{
			vi: &vi{},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		if msg, some := m.state.vi.update(typed); some {
			switch msg.command {
			case "q":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	return m.state.vi.pending()
}

func Run(config postbote.Flags) error {
	p := tea.NewProgram(newModel([]mail.Folder{
		{
			Name:     "Inbox",
			Path:     "Inbox",
			Messages: []mail.Message{},
		},
	}))
	_, err := p.Run()
	return err
}
