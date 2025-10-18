package processlist

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
)

type ProcessList struct {
	processes []taskmanager.TaskProcess
}

func NewProcessList(procs []taskmanager.TaskProcess) *ProcessList {
	return &ProcessList{
		processes: procs,
	}
}

func (p *ProcessList) Init() tea.Cmd {
	return nil
}

func (p *ProcessList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case messages.ProcessListMsg:
		p.processes = typedMsg.Processes
		return p, nil
	}

	return p, nil
}

func (p *ProcessList) View() string {
	if len(p.processes) == 0 {
		return "No processes available"
	}

	title := "Process List"
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render(title)

	header = lipgloss.PlaceHorizontal(80, lipgloss.Center, header)

	var content strings.Builder
	content.WriteString(header + "\n\n")

	for i, process := range p.processes {
		processLine := fmt.Sprintf("%d. Process (PID: %d)", i+1, process.PID)
		content.WriteString(processLine + "\n")
	}

	return content.String()
}
