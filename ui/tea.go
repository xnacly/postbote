package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/xnacly/postbote"
	"github.com/xnacly/postbote/mail"
)

var (
	base         = lipgloss.Color("#eff1f5")
	mauve        = lipgloss.Color("#8839ef")
	crust        = lipgloss.Color("#4c4f69")
	paneStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	folderStyle  = lipgloss.NewStyle().Bold(true).Foreground(mauve)
	messageStyle = lipgloss.NewStyle().Foreground(crust)

	statusStyle = lipgloss.NewStyle().Background(crust).Foreground(base).Padding(0, 2)
)

type childPane struct {
	Id      string
	Name    string
	Mime    string
	Content string
}

type model struct {
	vi        vi
	activeIdx uint
	width     int
	height    int
	ready     bool

	viewportWidth  int
	parent         mail.Folder
	parentViewPort viewport.Model

	current         mail.Folder
	currentViewPort viewport.Model

	child         childPane
	childViewPort viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		m.viewportWidth = (m.width / 3) - 2
		viewportHeights := m.height - 4

		if !m.ready {
			m.currentViewPort = viewport.New(m.viewportWidth, viewportHeights)
			m.parentViewPort = viewport.New(m.viewportWidth, viewportHeights)
			m.childViewPort = viewport.New(m.viewportWidth, viewportHeights)
			m.ready = true
		} else {
			m.currentViewPort.Height = viewportHeights
			m.currentViewPort.Width = m.viewportWidth
			m.parentViewPort.Height = viewportHeights
			m.parentViewPort.Width = m.viewportWidth
			m.childViewPort.Height = viewportHeights
			m.childViewPort.Width = m.viewportWidth
		}
	case tea.KeyMsg:
		if msg, some := m.vi.update(typed); some {
			switch msg.command {
			case "q":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) status(width int) string {
	left := fmt.Sprintf("%s (d:%d;m:%d)",
		m.current.Path,
		len(m.current.Folders),
		len(m.current.Messages),
	)

	viPending := m.vi.pending()

	left = lipgloss.NewStyle().
		Width(width - lipgloss.Width(viPending) - 4).
		Align(lipgloss.Left).
		Render(left)

	return left + viPending
}

func (m model) pane(vp *viewport.Model, folder mail.Folder) string {
	lines := []string{}
	for i, f := range folder.Folders {
		lines = append(lines, fmt.Sprint(
			" [",
			i,
			"] ",
			folderStyle.Render(f.Name),
		))
	}

	for i, msg := range folder.Messages {
		lines = append(lines, fmt.Sprint(
			" [",
			i+len(folder.Folders),
			"] ",
			messageStyle.Render(msg.From),
			"; ",
			messageStyle.Render(msg.Subject),
		))
	}

	vp.SetContent(strings.Join(lines, "\n"))

	return paneStyle.Render(vp.View())
}

func (m model) childpane() string {
	wrapped := lipgloss.NewStyle().
		Width(m.viewportWidth).
		Render(m.child.Content)

	switch m.child.Mime {
	case "text/plain":
		m.childViewPort.SetContent(wrapped)
	default:
		panic("execute command defined in PostboteConfig.Mime.<type>")
	}

	return paneStyle.Render(m.childViewPort.View())
}

func (m model) View() string {
	layout := lipgloss.JoinHorizontal(lipgloss.Top,
		m.pane(&m.parentViewPort, m.parent),
		m.pane(&m.currentViewPort, m.current),
		m.childpane(),
	)

	status := statusStyle.Render(m.status(m.width))
	return lipgloss.JoinVertical(lipgloss.Left, layout, status)
}

func Run(config postbote.Flags) error {
	exmplmsg := `Hi João,

I'm glad you enjoyed my article. Always fun to share some newfound knowledge :).
Thanks for reaching out!

Best
Matteo

On Wednesday, July 16th, 2025 at 09:08, John (Joao) Batalha <john@amplemarket.com> wrote:

> hey - just emailing you to thank you for writing the post on lexers. really enjoyed it!
>
> joão batalha
> CEO @ Amplemarket
`
	m := model{
		parent: mail.Folder{
			Name:     "MAIL",
			Path:     "MAIL/",
			Messages: []mail.Message{},
			Folders: []mail.Folder{
				{Name: "Inbox"},
				{Name: "Sent"},
				{Name: "Spam"},
			},
		},
		current: mail.Folder{
			Name:    "Inbox",
			Path:    "MAIL/Inbox",
			Folders: []mail.Folder{},
			Messages: []mail.Message{
				{
					UID:     0,
					Subject: "Re: post on lexers",
					From:    "contact@xnacly.me",
					Date:    time.Now(),
					Attachments: []mail.Attachment{
						{
							ID:       "0",
							Name:     "text",
							MimeType: "text/plain",
							Size:     int64(len(exmplmsg)),
						},
					},
				},
			},
		},
		child: childPane{
			Id:      "0",
			Name:    "Re: post on lexers",
			Mime:    "text/plain",
			Content: exmplmsg,
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
