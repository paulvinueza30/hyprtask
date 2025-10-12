package viewtitle

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
)

type ViewTitle struct {
	Text string
}

func NewViewTitle(text string) *ViewTitle {
	return &ViewTitle{
		Text: text,
	}
}

func (vt *ViewTitle) Init() tea.Cmd {
	return nil
}

func (vt *ViewTitle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return vt, nil
}

func (vt *ViewTitle) View() string {
	return theme.Get().ViewModel.Title.Render(vt.Text)
}
