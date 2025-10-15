package workspaceselector

import (
	"fmt"

	"github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/ui/components/viewtitle"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
)

type WorkspaceSelectorView struct {
	FlexBox      *flexbox.FlexBox
	stateManager *stateManager

	Title  tea.Model
	width  int
	height int
}

func NewWorkspaceSelectorView() *WorkspaceSelectorView {
	flexbox := flexbox.New(0, 0)

	ws := &WorkspaceSelectorView{
		FlexBox:      flexbox,
		Title:        viewtitle.NewViewTitle("Select A Workspace"),
		stateManager: newStateManager(),
	}

	return ws
}

func (ws *WorkspaceSelectorView) Init() tea.Cmd {
	titleCmd := ws.Title.Init()

	var workspaceCmds []tea.Cmd
	for _, workspace := range ws.stateManager.getWorkspaces() {
		workspaceCmds = append(workspaceCmds, workspace.Init())
	}

	allCmds := append([]tea.Cmd{titleCmd}, workspaceCmds...)
	return tea.Batch(allCmds...)
}

func (ws *WorkspaceSelectorView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case messages.WorkspaceDataMsg:
		ws.stateManager.createWorkspaceBoxes(msg.Workspaces)
	case tea.KeyMsg:
		cmd := ws.stateManager.handleKeyMsg(msg)
		return ws, cmd
	case tea.WindowSizeMsg:
		widthPadding := ws.calculateWidthPadding(msg.Width)
		heightPadding := ws.calculateHeightPadding(msg.Height)

		ws.FlexBox.SetWidth(msg.Width - widthPadding).SetHeight(msg.Height - heightPadding)

		ws.width = msg.Width
		ws.height = msg.Height
	}

	var titleCmd tea.Cmd
	ws.Title, titleCmd = ws.Title.Update(msg)
	cmds = append(cmds, titleCmd)

	workspaces := ws.stateManager.getWorkspaces()
	updatedWorkspaces := make([]tea.Model, len(workspaces))

	for i, workspace := range workspaces {
		var workspaceCmd tea.Cmd
		updatedWorkspaces[i], workspaceCmd = workspace.Update(msg)
		cmds = append(cmds, workspaceCmd)
	}

	ws.stateManager.updateWorkspaces(updatedWorkspaces)

	return ws, tea.Batch(cmds...)
}

func (ws *WorkspaceSelectorView) View() string {
	workspaceCountHeader := theme.Get().WorkspaceView.Title.Render(fmt.Sprintf("%d Workspaces", ws.stateManager.getWorkspaceCount()))
	title := ws.Title.View()
	workspaceGrid := ws.createWorkspaceGrid()
	instructions := theme.Get().WorkspaceView.Details.Render(keymap.Get().GetHelpText())

	centeredHeader := lipgloss.PlaceHorizontal(ws.width, lipgloss.Center, workspaceCountHeader)
	var marginTop, marginBottom int
	if ws.height > 20 {
		marginTop = 1
		marginBottom = 1
	}

	header := lipgloss.JoinVertical(lipgloss.Center, centeredHeader, title)

	headerStyled := lipgloss.NewStyle().MarginTop(marginTop).Render(header)
	gridStyled := lipgloss.NewStyle().MarginTop(marginTop).MarginBottom(marginBottom).Render(workspaceGrid)
	instructionsStyled := lipgloss.NewStyle().Render(instructions)

	return lipgloss.JoinVertical(lipgloss.Center, headerStyled, gridStyled, instructionsStyled)
}

func (ws *WorkspaceSelectorView) createWorkspaceGrid() string {
	workspaces := ws.stateManager.getWorkspaces()
	if len(workspaces) == 0 {
		return "No workspaces available"
	}

	scrollOffset := ws.stateManager.getScrollOffset()
	var flexRows []*flexbox.Row

	cols := 2
	if len(workspaces) < 2 {
		cols = len(workspaces)
	}

	maxVisibleRows := 2
	if ws.height <= 20 {
		maxVisibleRows = 1
	}

	startIndex := scrollOffset * cols
	row := ws.FlexBox.NewRow()
	rowCount := 0

	for i := startIndex; i < len(workspaces) && rowCount < maxVisibleRows; i++ {
		workspace := workspaces[i]
		cell := flexbox.NewCell(1, 1).SetContent(workspace.View())
		row.AddCells(cell)

		workspaceIndexInRow := i - startIndex
		if (workspaceIndexInRow+1)%cols == 0 || i == len(workspaces)-1 {
			flexRows = append(flexRows, row)
			rowCount++

			if i < len(workspaces)-1 && rowCount < maxVisibleRows {
				row = ws.FlexBox.NewRow()
			}
		}
	}

	ws.FlexBox.SetRows(flexRows)
	gridContent := ws.FlexBox.Render()

	// Add down arrow indicator if there are more workspaces below
	if ws.hasMoreWorkspacesBelow() {
		remainingCount := ws.getRemainingWorkspacesCount()
		indicator := fmt.Sprintf("↓ %d more below ↓", remainingCount)
		centeredIndicator := lipgloss.PlaceHorizontal(ws.width, lipgloss.Center, indicator)
		return lipgloss.JoinVertical(lipgloss.Center, gridContent, centeredIndicator)
	}

	return gridContent
}

func (ws *WorkspaceSelectorView) hasMoreWorkspacesBelow() bool {
	workspaces := ws.stateManager.getWorkspaces()
	if len(workspaces) == 0 {
		return false
	}

	scrollOffset := ws.stateManager.getScrollOffset()
	cols := 2
	if len(workspaces) < 2 {
		cols = len(workspaces)
	}

	maxVisibleRows := 2
	if ws.height <= 20 {
		maxVisibleRows = 1
	}

	totalVisibleWorkspaces := (scrollOffset + maxVisibleRows) * cols
	return totalVisibleWorkspaces < len(workspaces)
}

func (ws *WorkspaceSelectorView) getRemainingWorkspacesCount() int {
	workspaces := ws.stateManager.getWorkspaces()
	if len(workspaces) == 0 {
		return 0
	}

	scrollOffset := ws.stateManager.getScrollOffset()
	cols := 2
	if len(workspaces) < 2 {
		cols = len(workspaces)
	}

	maxVisibleRows := 2
	if ws.height <= 20 {
		maxVisibleRows = 1
	}

	totalVisibleWorkspaces := (scrollOffset + maxVisibleRows) * cols
	remaining := len(workspaces) - totalVisibleWorkspaces

	if remaining < 0 {
		return 0
	}
	return remaining
}

// calculateWidthPadding calculates padding based on width only
func (ws *WorkspaceSelectorView) calculateWidthPadding(width int) int {
	const (
		widthThreshold     = 55
		widthBreakPoint    = 110
		standardMultiplier = 8.7
		highMultiplier     = 9.0
	)

	if width <= widthThreshold {
		return 0
	}

	if width <= widthBreakPoint {
		padding := int(float64(width-widthThreshold) * standardMultiplier)
		return padding
	} else {
		padding := int(float64(width-widthThreshold) * highMultiplier)
		return padding
	}
}

// calculateHeightPadding calculates padding based on height only
func (ws *WorkspaceSelectorView) calculateHeightPadding(height int) int {
	const (
		heightThreshold    = 10
		heightBreakPoint   = 20
		standardMultiplier = .9
		highMultiplier     = .8
	)

	if height <= heightBreakPoint {
		padding := int(float64(height-heightThreshold) * standardMultiplier)
		return padding
	} else {
		padding := int(float64(height-heightThreshold) * highMultiplier)
		return padding
	}
}
