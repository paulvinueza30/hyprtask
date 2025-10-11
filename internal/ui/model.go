package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"

	"github.com/paulvinueza30/hyprtask/internal/ui/components"
)

type Model struct {
	displayDataChan <-chan viewmodel.DisplayData
	viewActionChan  chan<- viewmodel.ViewAction

	cursor      int // idk how im gonna get it but we'll see
	displayData viewmodel.DisplayData
	
	width  int
	height int
}

func NewModel(ddChan chan viewmodel.DisplayData, viewActChan chan viewmodel.ViewAction) *Model {
	return &Model{displayDataChan: ddChan, viewActionChan: viewActChan }
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
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *Model) View() string {
  styledHeader := components.HeaderStyle.
        Width(m.width).
        Align(lipgloss.Center)

    // 2. Render the actual text using the configured style.
    headerText := styledHeader.Render("HyprTask")

    // For now, you can just return the header string.
    // Later, you will use lipgloss.JoinVertical() to add body content.
    return headerText
}
func (m *Model) listenToDisplayDataChan() tea.Cmd {
	return func() tea.Msg {
		displayData := <-m.displayDataChan
		return displayData
	}
}
