package processlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens"
)

type stateManager struct {
	workspaceID *int
	processList []taskmanager.TaskProcess
}

func newStateManager() *stateManager {
	return &stateManager{
		workspaceID: nil,
		processList: []taskmanager.TaskProcess{},
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
