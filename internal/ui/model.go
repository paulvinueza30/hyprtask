package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens/processlist"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens/workspaceselector"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

type Model struct {
	displayDataChan <-chan viewmodel.DisplayData
	viewActionChan  chan<- viewmodel.ViewAction
	taskActionChan  chan<- taskmanager.TaskAction

	displayData  viewmodel.DisplayData
	windowWidth  int
	windowHeight int

	screens      map[screens.ScreenType]tea.Model
	activeScreen screens.ScreenType
	
	processListWorkspaceID *int // nil = all processes, &workspaceID = specific workspace
}

func NewModel(ddChan chan viewmodel.DisplayData, viewActChan chan viewmodel.ViewAction, taskActChan chan taskmanager.TaskAction) *Model {
	theme.Init()
	keymap.Init()

	model := &Model{
		displayDataChan: ddChan,
		viewActionChan:  viewActChan,
		taskActionChan:  taskActChan,
		// TODO: Make this dynamic based on if theres workspaces or not
		activeScreen: screens.WorkspaceSelector,
		screens: map[screens.ScreenType]tea.Model{
			screens.WorkspaceSelector: workspaceselector.NewWorkspaceSelectorView(),
			screens.ProcessList:       processlist.NewProcessList([]taskmanager.TaskProcess{}),
		},
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
	var broadcastMsg tea.Msg

	switch msg := msg.(type) {
	case viewmodel.DisplayData:
		m.displayData = msg
		cmds = append(cmds, m.listenToDisplayDataChan())
		cmds = append(cmds, m.updateWorkspaceSelectorWithDisplayData()...)
		cmds = append(cmds, m.updateProcessListWithDisplayData()...)

	case messages.ChangeScreenMsg[messages.ProcessListMsg]:
		processes := m.getProcsForWorkspace(msg.ScreenMsg.WorkspaceID)
		msg.ScreenMsg.Processes = processes
		m.SetActiveScreen(msg.ScreenType)
		
		// Store the workspace context
		m.processListWorkspaceID = msg.ScreenMsg.WorkspaceID
		
		broadcastMsg = msg.ScreenMsg
	case messages.ChangeScreenMsg[messages.WorkspaceListMsg]:
		m.SetActiveScreen(msg.ScreenType)
		broadcastMsg = msg.ScreenMsg

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// Create a copy for broadcasting (space for the header)
		broadcastMsg = tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - 3,
		}
	case tea.KeyMsg:
		if activeScreen, exists := m.screens[m.activeScreen]; exists {
			updatedScreen, cmd := activeScreen.Update(msg)
			m.screens[m.activeScreen] = updatedScreen
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case messages.ChangeSortOptionMsg:
		m.sendSortActionToViewModel(msg)
	case messages.KillProcessMsg:
		m.sendKillActionToTaskManager(msg)
	default:
	if activeScreen, exists := m.screens[m.activeScreen]; exists {
			updatedScreen, cmd := activeScreen.Update(msg)
			m.screens[m.activeScreen] = updatedScreen
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	if broadcastMsg != nil {
		cmds = append(cmds, m.broadcastToScreens(broadcastMsg)...)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	header := theme.Get().Header.
		Width(m.windowWidth).Height(m.windowHeight / 10).
		Align(lipgloss.Center).Render("HyprTask")

	content := m.screens[m.activeScreen].View()

	return lipgloss.JoinVertical(lipgloss.Center, header, content)
}

func (m *Model) SetActiveScreen(st screens.ScreenType) {
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

func (m *Model) getProcsForWorkspace(workspaceID *int) []taskmanager.TaskProcess {
	if workspaceID == nil {
		return m.displayData.All
	} else {
		if workspaceData, exists := m.displayData.Hypr.WorkspaceToProcs[*workspaceID]; exists {
			return workspaceData.ActiveProcs
		}
		return []taskmanager.TaskProcess{}
	}
}

func (m *Model) broadcastToScreens(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	for screenType, screen := range m.screens {
		var cmd tea.Cmd
		m.screens[screenType], cmd = screen.Update(msg)
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (m *Model) updateWorkspaceSelectorWithDisplayData() []tea.Cmd {
	var cmds []tea.Cmd

	workspaceMsg := messages.NewWorkspaceDataMsg(m.displayData.Hypr)
	if screen, exists := m.screens[screens.WorkspaceSelector]; exists {
		updatedScreen, cmd := screen.Update(workspaceMsg)
		m.screens[screens.WorkspaceSelector] = updatedScreen
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}

func (m *Model) updateProcessListWithDisplayData() []tea.Cmd {
	var cmds []tea.Cmd


	var processMsg messages.ProcessListMsg
	
	if m.processListWorkspaceID == nil {
		// Currently viewing all processes
		processMsg = messages.NewAllProcessesMsg()
		processMsg.Processes = m.displayData.All
	} else {
		// Currently viewing a specific workspace
		workspaceID := *m.processListWorkspaceID
		processMsg = messages.ProcessListMsg{
			WorkspaceID:   m.processListWorkspaceID,
			WorkspaceName: m.getWorkspaceNameByID(workspaceID),
			Processes:     m.getProcsForWorkspace(m.processListWorkspaceID),
		}
	}
	
	if screen, exists := m.screens[screens.ProcessList]; exists {
		updatedScreen, cmd := screen.Update(processMsg)
		m.screens[screens.ProcessList] = updatedScreen
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}
func (m *Model) sendSortActionToViewModel(msg messages.ChangeSortOptionMsg){
	m.viewActionChan <- viewmodel.ViewAction{
		NewSortKey: msg.Key,
		NewSortOrder: msg.Order,
	}
	logger.Log.Info("Sending sort action to viewmodel", "action", msg)
}
func (m *Model) sendKillActionToTaskManager(msg messages.KillProcessMsg){
	m.taskActionChan <- taskmanager.TaskAction{
		Type:    taskmanager.TaskActionKill,
		Payload: taskmanager.KillProcessPayload{
			PID: msg.PID,
			Force: msg.Force,
		},
	}
	logger.Log.Info("Sending kill action to taskmanager", "action", msg)
}

func (m *Model) getWorkspaceNameByID(workspaceID int) *string {
	if workspaceData, exists := m.displayData.Hypr.WorkspaceToProcs[workspaceID]; exists {
		return &workspaceData.WorkspaceName
	}
	return nil
}