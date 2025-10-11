package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

type Model struct {
	displayDataChan <-chan viewmodel.DisplayData
	viewActionChan  chan<- viewmodel.ViewAction

	cursor      int // idk how im gonna get it but we'll see
	displayData viewmodel.DisplayData
}

func NewModel(ddChan chan viewmodel.DisplayData, viewActChan chan viewmodel.ViewAction) *Model {
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
	}
	return m, nil
}

func (m *Model) View() string {
	return ""
}
func (m *Model) listenToDisplayDataChan() tea.Cmd {
	return func() tea.Msg {
		displayData := <-m.displayDataChan
		return displayData
	}
}
