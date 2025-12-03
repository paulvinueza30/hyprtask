package processlist

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)
type sortOptions struct {
	key viewmodel.SortKey
	order viewmodel.SortOrder
}

type state struct {
	workspaceID   *int
	workspaceName *string
	processList   []taskmanager.TaskProcess
	sortOptions   sortOptions
}
type stateManager struct {
	state *state
	table *table.Model
}

func newStateManager(procs []taskmanager.TaskProcess, table *table.Model) *stateManager {
	sortOptions := sortOptions{
		key: viewmodel.SortByNone,
		order: viewmodel.OrderNone,
	}
	return &stateManager{
		state: &state{
			workspaceID:   nil,
			workspaceName: nil,
			processList:   procs,
			sortOptions:   sortOptions,
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
	case "sort_key_left":
		return sm.sortKeyLeft()
	case "sort_key_right":
		return sm.sortKeyRight()
	case "toggle_sort_order":
		return sm.toggleSortOrder()
	case "kill_process":
		return sm.killProcess(false)
	case "kill_process_force":
		return sm.killProcess(true)
	}

	return nil
}

func (sm *stateManager) changeToWorkspaceSelectorView() tea.Cmd {
	return func() tea.Msg {
		return messages.NewChangeScreenMsg(screens.WorkspaceSelector, messages.WorkspaceListMsg{})
	}
}

func (sm *stateManager) setState(msg messages.ProcessListMsg) {
	currentSortOptions := sm.state.sortOptions
	sm.state = &state{
		workspaceID:   msg.WorkspaceID,
		workspaceName: msg.WorkspaceName,
		processList:   msg.Processes,
		sortOptions:   currentSortOptions,
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

var visualSortOrder = []viewmodel.SortKey{
	viewmodel.SortByPID,
	viewmodel.SortByProgramName,
	viewmodel.SortByUser,
	viewmodel.SortByCPU,
	viewmodel.SortByMEM,
}

func (sm *stateManager) getCurrentVisualIndex() int {
	currentKey := sm.state.sortOptions.key
	if currentKey == viewmodel.SortByNone {
		return -1
	}
	for i, key := range visualSortOrder {
		if key == currentKey {
			return i
		}
	}
	return -1
}

func (sm *stateManager) sortKeyLeft() tea.Cmd {
	currentIndex := sm.getCurrentVisualIndex()
	if currentIndex == -1 {
		sm.state.sortOptions.key = visualSortOrder[len(visualSortOrder)-1]
	} else if currentIndex > 0 {
		sm.state.sortOptions.key = visualSortOrder[currentIndex-1]
	} else {
		sm.state.sortOptions.key = viewmodel.SortByNone
		sm.state.sortOptions.order = viewmodel.OrderNone
		return func() tea.Msg {
			return messages.NewChangeSortOptionMsg(sm.state.sortOptions.key, sm.state.sortOptions.order)
		}
	}
	
	if sm.state.sortOptions.order == viewmodel.OrderNone {
		sm.state.sortOptions.order = viewmodel.OrderDESC
	}
	
	return func() tea.Msg {
		return messages.NewChangeSortOptionMsg(sm.state.sortOptions.key, sm.state.sortOptions.order)
	}
}

func (sm *stateManager) sortKeyRight() tea.Cmd {
	currentIndex := sm.getCurrentVisualIndex()
	if currentIndex == -1 {
		sm.state.sortOptions.key = visualSortOrder[0]
	} else if currentIndex < len(visualSortOrder)-1 {
		sm.state.sortOptions.key = visualSortOrder[currentIndex+1]
	} else {
		sm.state.sortOptions.key = viewmodel.SortByNone
		sm.state.sortOptions.order = viewmodel.OrderNone
		return func() tea.Msg {
			return messages.NewChangeSortOptionMsg(sm.state.sortOptions.key, sm.state.sortOptions.order)
		}
	}
	
	if sm.state.sortOptions.order == viewmodel.OrderNone {
		sm.state.sortOptions.order = viewmodel.OrderDESC
	}
	
	return func() tea.Msg {
		return messages.NewChangeSortOptionMsg(sm.state.sortOptions.key, sm.state.sortOptions.order)
	}
}
func (sm *stateManager) toggleSortOrder() tea.Cmd {
	switch sm.state.sortOptions.order {
	case viewmodel.OrderASC:
		sm.state.sortOptions.order = viewmodel.OrderDESC
	case viewmodel.OrderDESC:
		sm.state.sortOptions.order = viewmodel.OrderASC
	}
	return func() tea.Msg {
		return messages.NewChangeSortOptionMsg(sm.state.sortOptions.key, sm.state.sortOptions.order)
	}
}

func (sm *stateManager) updateTable(table *table.Model) {
	sm.table = table
}

func (sm *stateManager) killProcess(force bool) tea.Cmd {
	if sm.table == nil {
		return nil
	}
	selectedRow := sm.table.Cursor()
	if selectedRow >= 0 && selectedRow < len(sm.state.processList) {
		proc := sm.state.processList[selectedRow]
		return func() tea.Msg {
			return ShowConfirmationMsg{
				PID:         proc.PID,
				ProcessName: proc.ProgramName,
				Command:     proc.CommandLine,
				Force:       force,
			}
		}
	}
	return nil
}