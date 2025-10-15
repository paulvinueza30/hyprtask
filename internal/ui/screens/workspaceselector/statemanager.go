package workspaceselector

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/ui/components/workspacebox"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

type selected struct {
	row int
	col int
}

type workspaceState struct {
	count              int
	workspaces         []tea.Model // preserves order for model
	workspaceIDToIndex map[int]int // workspace ID -> slice index (O(1) lookup)
	selected           selected
	scrollOffset       int // current scroll position (row offset)
}

type stateManager struct {
	state *workspaceState
}

func (sm *stateManager) getWorkspaceIndex(sel selected) int {
	return sel.row*2 + sel.col
}

func (sm *stateManager) getMaxRows() int {
	return (len(sm.state.workspaces) + 1) / 2
}

func (sm *stateManager) isValidPosition(sel selected) bool {
	index := sm.getWorkspaceIndex(sel)
	return index < len(sm.state.workspaces)
}

func newStateManager() *stateManager {
	return &stateManager{
		state: &workspaceState{
			count:              0,
			workspaces:         make([]tea.Model, 0),
			workspaceIDToIndex: make(map[int]int),
			selected:           selected{row: 0, col: 0},
		},
	}
}

func (sm *stateManager) createWorkspaceBoxes(workspaceData []*viewmodel.WorkspaceData) {
	sm.state.count = len(workspaceData)

	// Track which workspaces exist in new data
	existingWorkspaceIDs := make(map[int]bool)
	newWorkspaces := make([]tea.Model, 0, len(workspaceData))
	newWorkspaceIDToIndex := make(map[int]int)

	for i, workspaceInfo := range workspaceData {
		if workspaceInfo == nil {
			continue
		}

		workspaceID := workspaceInfo.WorkspaceID
		existingWorkspaceIDs[workspaceID] = true
		newWorkspaceIDToIndex[workspaceID] = i

		// Reuse existing workspace box if it exists, otherwise create new one
		if oldIndex, exists := sm.state.workspaceIDToIndex[workspaceID]; exists && oldIndex < len(sm.state.workspaces) {
			// Reuse existing workspace box and update data
			existingBox := sm.state.workspaces[oldIndex]
			newWorkspaces = append(newWorkspaces, existingBox)

			if workspaceBox, ok := existingBox.(*workspacebox.WorkspaceBox); ok {
				workspaceBox.WindowCount = workspaceInfo.ActiveProcsCount
				workspaceBox.CPUUsage = workspaceInfo.TotalCPU
				workspaceBox.MemUsage = workspaceInfo.TotalMEM
			}
		} else {
			// Create new workspace box
			workspaceBox := workspacebox.NewWorkspaceBox(
				workspaceInfo.WorkspaceID,
				workspaceInfo.WorkspaceName,
			)
			workspaceBox.WindowCount = workspaceInfo.ActiveProcsCount
			workspaceBox.CPUUsage = workspaceInfo.TotalCPU
			workspaceBox.MemUsage = workspaceInfo.TotalMEM

			newWorkspaces = append(newWorkspaces, workspaceBox)
		}
	}

	// Update state
	sm.state.workspaces = newWorkspaces
	sm.state.workspaceIDToIndex = newWorkspaceIDToIndex

	sm.ensureValidPosition()
	sm.updateWorkspaceSelection()
}

func (sm *stateManager) getWorkspaceCount() int {
	return sm.state.count
}

func (sm *stateManager) getWorkspaces() []tea.Model {
	return sm.state.workspaces
}

func (sm *stateManager) updateWorkspaces(updatedWorkspaces []tea.Model) {
	// Update the slice with new workspace instances
	sm.state.workspaces = updatedWorkspaces
}

func (sm *stateManager) ensureValidPosition() {
	maxRows := sm.getMaxRows()

	if !sm.isValidPosition(sm.state.selected) {
		if sm.state.selected.row < maxRows && sm.isValidPosition(selected{row: sm.state.selected.row, col: 0}) {
			sm.state.selected.col = 0
		} else {
			for row := maxRows - 1; row >= 0; row-- {
				for col := 1; col >= 0; col-- {
					testSelected := selected{row: row, col: col}
					if sm.isValidPosition(testSelected) {
						sm.state.selected.row = row
						sm.state.selected.col = col
						return
					}
				}
			}
		}
	}
}

func (sm *stateManager) findNextRowInColumn(col, currentRow, direction int) int {
	if direction > 0 {
		for row := currentRow + 1; row < sm.getMaxRows(); row++ {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return row
			}
		}
		for row := 0; row < currentRow; row++ {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return row
			}
		}
	} else {
		for row := currentRow - 1; row >= 0; row-- {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return row
			}
		}
		for row := sm.getMaxRows() - 1; row > currentRow; row-- {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return row
			}
		}
	}
	return currentRow
}

func (sm *stateManager) findNextColInRow(row, currentCol, direction int) int {
	if direction > 0 {
		for col := currentCol + 1; col <= 1; col++ {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return col
			}
		}
		for col := 0; col < currentCol; col++ {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return col
			}
		}
	} else {
		for col := currentCol - 1; col >= 0; col-- {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return col
			}
		}
		for col := 1; col > currentCol; col-- {
			if sm.isValidPosition(selected{row: row, col: col}) {
				return col
			}
		}
	}
	return currentCol
}

func (sm *stateManager) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	action, handled := keymap.Get().HandleKeyMsg(msg)
	if !handled {
		return nil
	}

	switch action {
	case "quit":
		return tea.Quit
	case "navigate_left":
		sm.handleLeft()
		sm.updateWorkspaceSelection()
	case "navigate_right":
		sm.handleRight()
		sm.updateWorkspaceSelection()
	case "navigate_up":
		sm.handleUp()
		sm.updateWorkspaceSelection()
	case "navigate_down":
		sm.handleDown()
		sm.updateWorkspaceSelection()
	case "scroll_up":
		sm.scrollUp()
	case "scroll_down":
		sm.scrollDown()
	}

	return nil
}

func (sm *stateManager) updateWorkspaceSelection() {
	selectedIndex := sm.getWorkspaceIndex(sm.state.selected)

	// O(1) lookup for selected workspace
	if selectedIndex < len(sm.state.workspaces) {
		selectedWorkspace := sm.state.workspaces[selectedIndex]
		if workspaceBox, ok := selectedWorkspace.(*workspacebox.WorkspaceBox); ok {
			workspaceBox.SetSelected(true)
		}
	}

	// Unselect all other workspaces using slice iteration (maintains order)
	for i, workspace := range sm.state.workspaces {
		if i != selectedIndex {
			if workspaceBox, ok := workspace.(*workspacebox.WorkspaceBox); ok {
				if workspaceBox.IsSelected {
					workspaceBox.SetSelected(false)
				}
			}
		}
	}
}

func (sm *stateManager) handleLeft() {
	sm.state.selected.col = sm.findNextColInRow(sm.state.selected.row, sm.state.selected.col, -1)
}

func (sm *stateManager) handleRight() {
	sm.state.selected.col = sm.findNextColInRow(sm.state.selected.row, sm.state.selected.col, 1)
}

func (sm *stateManager) handleUp() {
	sm.state.selected.row = sm.findNextRowInColumn(sm.state.selected.col, sm.state.selected.row, -1)
	sm.ensureSelectionVisible()
}

func (sm *stateManager) handleDown() {
	sm.state.selected.row = sm.findNextRowInColumn(sm.state.selected.col, sm.state.selected.row, 1)
	sm.ensureSelectionVisible()
}

func (sm *stateManager) getScrollOffset() int {
	return sm.state.scrollOffset
}

func (sm *stateManager) scrollUp() {
	if sm.state.scrollOffset > 0 {
		sm.state.scrollOffset--
		sm.ensureSelectionAtTop()
	}
}

func (sm *stateManager) scrollDown() {
	maxRows := sm.getMaxRows()
	if sm.state.scrollOffset < maxRows-2 {
		sm.state.scrollOffset++
		sm.ensureSelectionAtTop()
	}
}

func (sm *stateManager) ensureSelectionAtTop() {
	sm.state.selected.row = sm.state.scrollOffset
	sm.ensureValidPosition()
	sm.updateWorkspaceSelection()
}

func (sm *stateManager) ensureSelectionVisible() {
	visibleRows := 2
	selectedRow := sm.state.selected.row

	if selectedRow < sm.state.scrollOffset {
		sm.state.scrollOffset = selectedRow
	} else if selectedRow >= sm.state.scrollOffset+visibleRows {
		sm.state.scrollOffset = selectedRow - visibleRows + 1
	}
}
