package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens/workspaceselector"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

type Model struct {
	displayDataChan <-chan viewmodel.DisplayData
	viewActionChan  chan<- viewmodel.ViewAction

	displayData  viewmodel.DisplayData
	windowWidth  int
	windowHeight int

	ActiveView tea.Model
}

func NewModel(ddChan chan viewmodel.DisplayData, viewActChan chan viewmodel.ViewAction) *Model {
	theme.Init()
	keymap.Init()
	model := &Model{displayDataChan: ddChan, viewActionChan: viewActChan}

	model.ActiveView = workspaceselector.NewWorkspaceSelectorView()

	return model
}

func (m *Model) Init() tea.Cmd {
	listenCmd := m.listenToDisplayDataChan()
	viewCmd := m.ActiveView.Init()

	return tea.Batch(listenCmd, viewCmd)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var listenCmd tea.Cmd

	switch msg := msg.(type) {
	case viewmodel.DisplayData:
		m.displayData = msg
		listenCmd = m.listenToDisplayDataChan()

		// Always send workspace data to workspace screen (even if count is 0 to clear display)
		logger.Log.Info("Processing DisplayData",
			"workspaceCount", msg.Hypr.WorkspaceCount,
			"processCount", len(msg.All))

		workspaceMsg := messages.NewWorkspaceDataMsg(msg.Hypr)
		// Send this message to the active view
		m.ActiveView, _ = m.ActiveView.Update(workspaceMsg)
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}

	var viewCmd tea.Cmd

	// TODO: propagate msg to all screens not just the active view
	m.ActiveView, viewCmd = m.ActiveView.Update(msg)

	return m, tea.Batch(listenCmd, viewCmd)
}

func (m *Model) View() string {
	// TODO: HEADER LOGO AND ASCII ART
	header := theme.Get().Header.
		Width(m.windowWidth).Height(m.windowHeight / 10).
		Align(lipgloss.Center).Render("HyprTask")

	content := m.ActiveView.View()

	return lipgloss.JoinVertical(lipgloss.Center, header, content)
}
func (m *Model) listenToDisplayDataChan() tea.Cmd {
	return func() tea.Msg {
		displayData := <-m.displayDataChan
		return displayData
	}
}
