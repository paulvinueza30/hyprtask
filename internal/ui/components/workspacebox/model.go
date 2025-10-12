package workspacebox

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/ui/theme"
)

type WorkspaceBox struct {
	ID          int
	Name        string
	WindowCount int
	CPUUsage    float64
	MemUsage    float64
	IsSelected  bool
}

type UpdateStatsMsg struct {
	ID          int
	WindowCount int
	CPUUsage    float64
	MemUsage    float64
}

func NewWorkspaceBox(id int, name string) *WorkspaceBox {
	return &WorkspaceBox{
		ID:          id,
		Name:        name,
		WindowCount: 0,
		CPUUsage:    0.0,
		MemUsage:    0.0,
		IsSelected:  false,
	}
}

func (wb *WorkspaceBox) Init() tea.Cmd {
	return nil
}

func (wb *WorkspaceBox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateStatsMsg:
		if msg.ID == wb.ID {
			wb.WindowCount = msg.WindowCount
			wb.CPUUsage = msg.CPUUsage
			wb.MemUsage = msg.MemUsage
		}
	}
	return wb, nil
}

func (wb *WorkspaceBox) View() string {
	boxStyle := theme.Get().WorkspaceView.Box
	if wb.IsSelected {
		boxStyle = theme.Get().WorkspaceView.SelectedBox
	}

	content := theme.Get().WorkspaceView.Title.Render(wb.Name)
	content += "\n"

	if wb.WindowCount == 1 {
		content += theme.Get().WorkspaceView.Details.Render(fmt.Sprintf("%d process", wb.WindowCount))
	} else {
		content += theme.Get().WorkspaceView.Details.Render(fmt.Sprintf("%d processes", wb.WindowCount))
	}

	content += "\n"
	content += theme.Get().WorkspaceView.Details.Render(fmt.Sprintf("CPU:%.1f%% | MEM %.1f%%", wb.CPUUsage, wb.MemUsage))

	return boxStyle.Render(content)
}

func (wb *WorkspaceBox) SetSelected(selected bool) {
	wb.IsSelected = selected
}
