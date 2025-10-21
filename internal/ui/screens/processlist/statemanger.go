package processlist

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens"
)

type state struct {
	workspaceID   *int
	workspaceName *string
	processList   []taskmanager.TaskProcess
}
type stateManager struct {
	state *state
	table *table.Model
}

func newStateManager(procs []taskmanager.TaskProcess, table *table.Model) *stateManager {
	return &stateManager{
		state: &state{
			workspaceID:   nil,
			workspaceName: nil,
			processList:   procs,
		},
		table: table,
	}
}

func (sm *stateManager) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	action, handled := keymap.Get().HandleKeyMsg(msg)
	if !handled {
		return nil
	}

	switch action {
	case "change_to_workspace_view":
		return sm.changeToWorkspaceSelectorView()
	case "quit":
		return tea.Quit
	}

	return nil
}

func (sm *stateManager) changeToWorkspaceSelectorView() tea.Cmd {
	return func() tea.Msg {
		return messages.NewChangeScreenMsg(screens.WorkspaceSelector, messages.WorkspaceListMsg{})
	}
}

func (sm *stateManager) setState(msg messages.ProcessListMsg) {
	sm.state = &state{
		workspaceID:   msg.WorkspaceID,
		workspaceName: msg.WorkspaceName,
		processList:   msg.Processes,
	}
}
func (sm *stateManager) getProcs() []taskmanager.TaskProcess {
	return sm.state.processList
}
func (sm *stateManager) getWorkspaceID() *int {
	return sm.state.workspaceID
}
func (sm *stateManager) getWorkspaceName() *string {
	return sm.state.workspaceName
}
