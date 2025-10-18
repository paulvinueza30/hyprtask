package processlist

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
)

type ProcessList struct {
	stateManager *stateManager
}

func NewProcessList(procs []taskmanager.TaskProcess) *ProcessList {
	return &ProcessList{
		stateManager: newStateManager(procs),
	}
}

func (p *ProcessList) Init() tea.Cmd {
	return nil
}

func (p *ProcessList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case messages.ProcessListMsg:
		p.stateManager.setState(typedMsg)
		return p, nil
	case tea.KeyMsg:
		return p, p.stateManager.handleKeyMsg(typedMsg)
	}

	return p, nil
}

func (p *ProcessList) View() string {
	wsName := p.stateManager.getWorkspaceName()
	wsNameStr := "all"
	if wsName != nil {
		wsNameStr = *wsName
	}

	title := fmt.Sprintf("Process List for workspace %s", wsNameStr)
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render(title)

	header = lipgloss.PlaceHorizontal(80, lipgloss.Center, header)

	var content strings.Builder
	for i, process := range p.stateManager.getProcs() {
		processLine := fmt.Sprintf("%d. Process (PID: %d)", i+1, process.PID)
		content.WriteString(processLine + "\n")
	}

	instructions := "Press " + keymap.Get().ChangeToWorkspaceSelectorScreen.Help().Key + " to change to workspace view"

	return lipgloss.JoinVertical(lipgloss.Center, header, content.String(), instructions)
}
