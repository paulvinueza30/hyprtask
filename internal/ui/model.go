package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens/processlist"
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

	screens      map[ScreenType]tea.Model
	activeScreen ScreenType
}

func NewModel(ddChan chan viewmodel.DisplayData, viewActChan chan viewmodel.ViewAction) *Model {
	theme.Init()
	keymap.Init()

	model := &Model{
		displayDataChan: ddChan,
		viewActionChan:  viewActChan,
		screens: map[ScreenType]tea.Model{
			WorkspaceSelector: workspaceselector.NewWorkspaceSelectorView(),
			ProcessList:       processlist.NewProcessList(),
		},
		activeScreen: WorkspaceSelector,
	}

	return model
}

func (m *Model) Init() tea.Cmd {
	listenCmd := m.listenToDisplayDataChan()
	var screenCmds []tea.Cmd

	for _, screen := range m.screens {
		screenCmds = append(screenCmds, screen.Init())
	}

	return tea.Batch(listenCmd, tea.Batch(screenCmds...))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var listenCmd tea.Cmd

	switch typedMsg := msg.(type) {
	case viewmodel.DisplayData:
		m.displayData = typedMsg
		listenCmd = m.listenToDisplayDataChan()

		msg = messages.NewWorkspaceDataMsg(typedMsg.Hypr)

	case tea.WindowSizeMsg:
		m.windowWidth = typedMsg.Width
		m.windowHeight = typedMsg.Height
		m.windowHeight = typedMsg.Height - 3
	}

	for screenType, screen := range m.screens {
		var cmd tea.Cmd
		m.screens[screenType], cmd = screen.Update(msg)
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, listenCmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	header := theme.Get().Header.
		Width(m.windowWidth).Height(m.windowHeight / 10).
		Align(lipgloss.Center).Render("HyprTask")

	content := m.screens[m.activeScreen].View()

	return lipgloss.JoinVertical(lipgloss.Center, header, content)
}

func (m *Model) SetActiveScreen(st ScreenType) {
	if _, exists := m.screens[st]; exists {
		m.activeScreen = st
	}
}

func (m *Model) listenToDisplayDataChan() tea.Cmd {
	return func() tea.Msg {
		displayData := <-m.displayDataChan
		return displayData
	}
}
