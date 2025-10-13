package workspaceselector

import (
	"fmt"

	"github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/ui/components/viewtitle"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
)

// Grid responsiveness constants - modify these to adjust behavior
const (
	// Screen size percentages for grid dimensions (dynamic scaling)
	GRID_WIDTH_PERCENTAGE_SMALL  = 0.9 // 90% for small screens
	GRID_WIDTH_PERCENTAGE_MEDIUM = 0.6 // 60% for medium screens
	GRID_WIDTH_PERCENTAGE_LARGE  = 0.5 // 50% for large screens

	GRID_HEIGHT_PERCENTAGE_SMALL  = 0.9 // 90% for small screens
	GRID_HEIGHT_PERCENTAGE_MEDIUM = 0.7 // 70% for medium screens
	GRID_HEIGHT_PERCENTAGE_LARGE  = 0.6 // 60% for large screens

	// Screen size thresholds for dynamic scaling
	SMALL_SCREEN_THRESHOLD = 80  // Below this: use small screen percentages
	LARGE_SCREEN_THRESHOLD = 120 // Above this: use large screen percentages

	// Minimum grid dimensions
	MIN_GRID_WIDTH  = 30 // Minimum characters wide
	MIN_GRID_HEIGHT = 10 // Minimum lines tall

	// Screen size thresholds for responsive behavior
	SMALL_SCREEN_WIDTH_THRESHOLD  = 30 // Below this: single column layout
	SMALL_SCREEN_HEIGHT_THRESHOLD = 16 // Below this: single row layout

	// Workspace box minimum dimensions
	MIN_BOX_WIDTH  = 25 // Minimum characters per workspace box
	MIN_BOX_HEIGHT = 6  // Minimum lines per workspace box

	// Default grid layout preferences
	DEFAULT_COLUMNS = 2 // Preferred number of columns
	DEFAULT_ROWS    = 2 // Preferred number of rows
)

type WorkspaceSelectorView struct {
	FlexBox      *flexbox.FlexBox
	stateManager *stateManager

	Title tea.Model
}

func NewWorkspaceSelectorView() *WorkspaceSelectorView {
	flexbox := flexbox.New(0, 0)

	ws := &WorkspaceSelectorView{
		FlexBox:      flexbox,
		Title:        viewtitle.NewViewTitle("Paul is gay and so is Gabe"),
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
		logger.Log.Info("Received workspace data",
			"count", len(msg.Workspaces),
			"workspaces", msg.Workspaces)

		ws.stateManager.createWorkspaceBoxes(msg.Workspaces)

		logger.Log.Info("Created workspace boxes",
			"count", ws.stateManager.getWorkspaceCount())
	case tea.KeyMsg:
		cmd := ws.stateManager.handleKeyMsg(msg)
		return ws, cmd
	case tea.WindowSizeMsg:
		// Calculate grid dimensions based on screen size with responsive formula
		gridWidth, gridHeight := ws.calculateGridDimensions(msg.Width, msg.Height)
		ws.FlexBox.SetWidth(gridWidth).SetHeight(gridHeight)

		// Add border to the FlexBox for debugging
		borderStyle := theme.Get().WorkspaceView.Box.Copy().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("3"))
		ws.FlexBox.SetStyle(borderStyle)

		// Force recalculation of the flexbox layout
		ws.FlexBox.ForceRecalculate()
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

	return workspaceCountHeader + "\n\n" + title + "\n\n" + workspaceGrid + "\n\n" + instructions
}

func (ws *WorkspaceSelectorView) createWorkspaceGrid() string {
	workspaces := ws.stateManager.getWorkspaces()
	if len(workspaces) == 0 {
		return "No workspaces available"
	}

	// Get actual FlexBox dimensions and calculate responsive grid
	width := ws.FlexBox.GetWidth()
	height := ws.FlexBox.GetHeight()
	cols, rows := ws.calculateResponsiveGridLayout(width, height, len(workspaces))

	scrollOffset := ws.stateManager.getScrollOffset()
	var flexRows []*flexbox.Row

	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		row := ws.FlexBox.NewRow()

		for colIndex := 0; colIndex < cols; colIndex++ {
			workspaceIndex := (scrollOffset+rowIndex)*cols + colIndex

			if workspaceIndex < len(workspaces) {
				workspace := workspaces[workspaceIndex]
				boxContent := workspace.View()
				cell := flexbox.NewCell(1, 1).SetContent(boxContent)
				row.AddCells(cell)
			}
		}

		flexRows = append(flexRows, row)
	}

	ws.FlexBox.SetRows(flexRows)
	return ws.FlexBox.Render()
}

func (ws *WorkspaceSelectorView) calculateGridDimensions(screenWidth, screenHeight int) (gridWidth, gridHeight int) {
	// Dynamic scaling based on screen size
	var widthPercentage, heightPercentage float64

	// Determine width percentage based on screen size
	if screenWidth <= SMALL_SCREEN_THRESHOLD {
		widthPercentage = GRID_WIDTH_PERCENTAGE_SMALL // 90% for small screens
	} else if screenWidth >= LARGE_SCREEN_THRESHOLD {
		widthPercentage = GRID_WIDTH_PERCENTAGE_LARGE // 50% for large screens
	} else {
		widthPercentage = GRID_WIDTH_PERCENTAGE_MEDIUM // 60% for medium screens
	}

	// Determine height percentage based on screen size
	if screenHeight <= SMALL_SCREEN_THRESHOLD {
		heightPercentage = GRID_HEIGHT_PERCENTAGE_SMALL // 90% for small screens
	} else if screenHeight >= LARGE_SCREEN_THRESHOLD {
		heightPercentage = GRID_HEIGHT_PERCENTAGE_LARGE // 60% for large screens
	} else {
		heightPercentage = GRID_HEIGHT_PERCENTAGE_MEDIUM // 70% for medium screens
	}

	gridWidth = int(float64(screenWidth) * widthPercentage)
	gridHeight = int(float64(screenHeight) * heightPercentage)

	// Ensure minimum dimensions
	if gridWidth < MIN_GRID_WIDTH {
		gridWidth = MIN_GRID_WIDTH
	}
	if gridHeight < MIN_GRID_HEIGHT {
		gridHeight = MIN_GRID_HEIGHT
	}

	// Ensure maximum dimensions (don't exceed screen with minimal padding)
	if gridWidth > screenWidth-2 {
		gridWidth = screenWidth - 2
	}
	if gridHeight > screenHeight-5 {
		gridHeight = screenHeight - 5
	}

	return gridWidth, gridHeight
}

func (ws *WorkspaceSelectorView) calculateResponsiveGridLayout(width, height, workspaceCount int) (cols, rows int) {
	// Calculate maximum possible columns and rows based on grid dimensions
	maxCols := width / MIN_BOX_WIDTH
	maxRows := height / MIN_BOX_HEIGHT

	// Ensure we have at least 1 column and 1 row
	if maxCols < 1 {
		maxCols = 1
	}
	if maxRows < 1 {
		maxRows = 1
	}

	// If grid is small, use single column or row
	if width < SMALL_SCREEN_WIDTH_THRESHOLD {
		cols = 1
		rows = workspaceCount
	} else if height < SMALL_SCREEN_HEIGHT_THRESHOLD {
		cols = workspaceCount
		rows = 1
	} else {
		// For larger grids, prefer default layout but adapt as needed
		cols = DEFAULT_COLUMNS
		if maxCols > DEFAULT_COLUMNS {
			cols = maxCols
		}

		rows = DEFAULT_ROWS
		if maxRows > DEFAULT_ROWS {
			rows = maxRows
		}

		// Adjust if we have fewer workspaces than grid capacity
		if workspaceCount < cols {
			cols = workspaceCount
		}
		if workspaceCount < rows {
			rows = workspaceCount
		}
	}

	// Ensure we don't exceed workspace count
	if cols > workspaceCount {
		cols = workspaceCount
	}
	if rows > workspaceCount {
		rows = workspaceCount
	}

	// Calculate actual rows needed based on workspace count and columns
	actualRows := (workspaceCount + cols - 1) / cols // Ceiling division
	if actualRows < rows {
		rows = actualRows
	}

	return cols, rows
}
