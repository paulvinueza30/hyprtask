package processlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
)

type ProcessList struct {
	Processes []taskmanager.TaskProcess
}

func NewProcessList() *ProcessList {
	return &ProcessList{}
}

func (p *ProcessList) Init() tea.Cmd {
	return nil
}

func (p *ProcessList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, nil
}

func (p *ProcessList) View() string {
	return ""
}
