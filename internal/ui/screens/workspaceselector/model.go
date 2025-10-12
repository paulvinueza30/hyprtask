package workspaceselector

import (
	"fmt"

	"github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/ui/components/viewtitle"
	"github.com/paulvinueza30/hyprtask/internal/ui/components/workspacebox"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
)

type Selected struct {
	Row int
	Col int
}

type WorkspaceSelectorView struct {
	Title          *viewtitle.ViewTitle
	Workspaces     []*workspacebox.WorkspaceBox
	Selected       Selected
	FlexBox        *flexbox.FlexBox
	WorkspaceCount int
	Navigation     *NavigationHandler
}

func NewWorkspaceSelectorView() *WorkspaceSelectorView {
	flexbox := flexbox.New(80, 20) // Default size, will be updated by window size
	ws := &WorkspaceSelectorView{
		Title:          viewtitle.NewViewTitle("Select Workspace"),
		Workspaces:     []*workspacebox.WorkspaceBox{}, // Start empty, populated by messages
		Selected:       Selected{Row: 0, Col: 0},
		FlexBox:        flexbox,
		WorkspaceCount: 0, // Start with 0, updated by messages
	}
	ws.Navigation = NewNavigationHandler(ws)
	return ws
}

func (ws *WorkspaceSelectorView) Init() tea.Cmd {
	titleCmd := ws.Title.Init()

	var workspaceCmds []tea.Cmd
	for _, workspace := range ws.Workspaces {
		workspaceCmds = append(workspaceCmds, workspace.Init())
	}

	allCmds := append([]tea.Cmd{titleCmd}, workspaceCmds...)
	return tea.Batch(allCmds...)
}

func (ws *WorkspaceSelectorView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case messages.WorkspaceDataMsg:
		// Store workspace count
		ws.WorkspaceCount = msg.Count

		// Log workspace data
		logger.Log.Info("Received workspace data",
			"count", len(msg.Workspaces),
			"workspaces", msg.Workspaces)

		// Create new workspace boxes with real data
		ws.Workspaces = make([]*workspacebox.WorkspaceBox, len(msg.Workspaces))
		for i, workspaceData := range msg.Workspaces {
			ws.Workspaces[i] = workspacebox.NewWorkspaceBox(
				workspaceData.WorkspaceID,
				workspaceData.WorkspaceName,
			)
			// Set the data directly on the new box
			ws.Workspaces[i].WindowCount = workspaceData.ActiveProcsCount
			ws.Workspaces[i].CPUUsage = workspaceData.TotalCPU
			ws.Workspaces[i].MemUsage = workspaceData.TotalMEM
		}

		// Ensure selected position is valid after workspace data changes
		ws.Navigation.EnsureValidPosition()

		logger.Log.Info("Created workspace boxes",
			"count", len(ws.Workspaces),
			"selectedRow", ws.Selected.Row,
			"selectedCol", ws.Selected.Col)
	case tea.KeyMsg:
		// Handle navigation and update workspace boxes
		updatedModel, cmd := ws.Navigation.HandleKeyMsg(msg)
		ws = updatedModel.(*WorkspaceSelectorView)

		// Update workspace boxes to reflect new selection
		selectedIndex := ws.Navigation.getWorkspaceIndex(ws.Selected)
		for i, workspace := range ws.Workspaces {
			workspace.SetSelected(i == selectedIndex)
		}

		return ws, cmd
	case tea.WindowSizeMsg:
		// Update FlexBox size (subtract space for title and padding)
		ws.FlexBox.SetWidth(msg.Width).SetHeight(msg.Height - 10)
	}

	// Update title
	updatedTitle, titleCmd := ws.Title.Update(msg)
	ws.Title = updatedTitle.(*viewtitle.ViewTitle)
	cmds = append(cmds, titleCmd)

	// Update all workspace boxes
	selectedIndex := ws.Navigation.getWorkspaceIndex(ws.Selected)
	for i, workspace := range ws.Workspaces {
		workspace.SetSelected(i == selectedIndex)
		var workspaceCmd tea.Cmd
		updatedWorkspace, workspaceCmd := workspace.Update(msg)
		ws.Workspaces[i] = updatedWorkspace.(*workspacebox.WorkspaceBox)
		cmds = append(cmds, workspaceCmd)
	}

	return ws, tea.Batch(cmds...)
}

func (ws *WorkspaceSelectorView) View() string {
	// Workspace count header
	workspaceCountHeader := theme.Get().WorkspaceView.Title.Render(fmt.Sprintf("%d Workspaces", ws.WorkspaceCount))
	title := ws.Title.View()

	// Skip rendering if no workspaces
	if len(ws.Workspaces) == 0 {
		return workspaceCountHeader + "\n\n" + title + "\n\nNo workspaces available"
	}

	// Create dynamic grid layout (2 workspaces per row)
	var rows []*flexbox.Row

	for i := 0; i < len(ws.Workspaces); i += 2 {
		row := ws.FlexBox.NewRow()

		// Add first workspace in this row
		if i < len(ws.Workspaces) {
			workspace := ws.Workspaces[i]
			boxContent := workspace.View()
			cell := flexbox.NewCell(1, 1).SetContent(boxContent)
			row.AddCells(cell)
		}

		// Add second workspace in this row (if it exists)
		if i+1 < len(ws.Workspaces) {
			workspace := ws.Workspaces[i+1]
			boxContent := workspace.View()
			cell := flexbox.NewCell(1, 1).SetContent(boxContent)
			row.AddCells(cell)
		}

		rows = append(rows, row)
	}

	// Set the rows and render
	ws.FlexBox.SetRows(rows)
	workspaceGrid := ws.FlexBox.Render()

	// Add navigation instructions
	instructions := theme.Get().WorkspaceView.Details.Render(keymap.Get().GetHelpText())

	return workspaceCountHeader + "\n\n" + title + "\n\n" + workspaceGrid + "\n\n" + instructions
}
