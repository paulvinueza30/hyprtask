package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)
type RenderContext struct {
	Width int
	Height int
	HasFocus bool
}
type Model struct {
	displayDataChan <-chan viewmodel.DisplayData
	viewActionChan  chan<- viewmodel.ViewAction

	displayData viewmodel.DisplayData
	
	RenderContext RenderContext
	
}

func NewModel(ddChan chan viewmodel.DisplayData, viewActChan chan viewmodel.ViewAction) *Model {
	theme.Init()
	return &Model{displayDataChan: ddChan, viewActionChan: viewActChan}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.listenToDisplayDataChan(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewmodel.DisplayData:
		m.displayData = msg
		return m, m.listenToDisplayDataChan()
	case tea.WindowSizeMsg:
		m.RenderContext.Height = msg.Height
		m.RenderContext.Width = msg.Width
		return m, nil
	}

	return m, nil
}

func (m *Model) View() string {
  header := theme.Get().Header.
        Width(m.RenderContext.Width).
        Align(lipgloss.Center).Render("HyprTask")

    return header 
}
func (m *Model) listenToDisplayDataChan() tea.Cmd {
	return func() tea.Msg {
		displayData := <-m.displayDataChan
		return displayData
	}
}
