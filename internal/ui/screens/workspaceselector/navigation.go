package workspaceselector

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
)

// NavigationHandler handles all keyboard navigation for workspace selector
type NavigationHandler struct {
	ws *WorkspaceSelectorView
}

// Helper methods for 2D navigation
func (nh *NavigationHandler) getWorkspaceIndex(selected Selected) int {
	return selected.Row*2 + selected.Col
}

func (nh *NavigationHandler) getMaxRows() int {
	return (len(nh.ws.Workspaces) + 1) / 2
}

func (nh *NavigationHandler) isValidPosition(selected Selected) bool {
	index := nh.getWorkspaceIndex(selected)
	return index < len(nh.ws.Workspaces)
}

// NewNavigationHandler creates a new navigation handler
func NewNavigationHandler(ws *WorkspaceSelectorView) *NavigationHandler {
	return &NavigationHandler{ws: ws}
}

// EnsureValidPosition ensures the current selection is valid after workspace data changes
func (nh *NavigationHandler) EnsureValidPosition() {
	maxRows := nh.getMaxRows()
	if nh.ws.Selected.Row >= maxRows {
		nh.ws.Selected.Row = maxRows - 1
		if nh.ws.Selected.Row < 0 {
			nh.ws.Selected.Row = 0
		}
	}
	if nh.ws.Selected.Col >= 2 {
		nh.ws.Selected.Col = 1
	}

	// Make sure the selected position is valid
	if !nh.isValidPosition(nh.ws.Selected) {
		// Find the last valid position
		for row := maxRows - 1; row >= 0; row-- {
			for col := 1; col >= 0; col-- {
				testSelected := Selected{Row: row, Col: col}
				if nh.isValidPosition(testSelected) {
					nh.ws.Selected.Row = row
					nh.ws.Selected.Col = col
					break
				}
			}
		}
	}
}

// HandleKeyMsg processes all key messages and returns the appropriate command
func (nh *NavigationHandler) HandleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action, handled := keymap.Get().HandleKeyMsg(msg)
	if !handled {
		return nh.ws, nil
	}

	switch action {
	case "quit":
		return nh.ws, tea.Quit
	case "navigate_left":
		nh.handleLeft()
	case "navigate_right":
		nh.handleRight()
	case "navigate_up":
		nh.handleUp()
	case "navigate_down":
		nh.handleDown()
	}

	return nh.ws, nil
}

// handleLeft moves left in the same row
func (nh *NavigationHandler) handleLeft() {
	if nh.ws.Selected.Col > 0 {
		nh.ws.Selected.Col--
	} else {
		// Wrap to right column of same row
		nh.ws.Selected.Col = 1
	}
}

// handleRight moves right in the same row
func (nh *NavigationHandler) handleRight() {
	if nh.ws.Selected.Col < 1 {
		testSelected := Selected{Row: nh.ws.Selected.Row, Col: nh.ws.Selected.Col + 1}
		if nh.isValidPosition(testSelected) {
			nh.ws.Selected.Col++
		} else {
			// Wrap to left column of same row
			nh.ws.Selected.Col = 0
		}
	} else {
		// Wrap to left column of same row
		nh.ws.Selected.Col = 0
	}
}

// handleUp moves up to same column in previous row
func (nh *NavigationHandler) handleUp() {
	if nh.ws.Selected.Row > 0 {
		nh.ws.Selected.Row--
	} else {
		// Wrap to last row, same column
		nh.ws.Selected.Row = nh.getMaxRows() - 1
	}
	// Ensure the position is valid
	if !nh.isValidPosition(nh.ws.Selected) {
		nh.ws.Selected.Col = 0
	}
}

// handleDown moves down to same column in next row
func (nh *NavigationHandler) handleDown() {
	if nh.ws.Selected.Row < nh.getMaxRows()-1 {
		nh.ws.Selected.Row++
	} else {
		// Wrap to first row, same column
		nh.ws.Selected.Row = 0
	}
	// Ensure the position is valid
	if !nh.isValidPosition(nh.ws.Selected) {
		nh.ws.Selected.Col = 0
	}
}
