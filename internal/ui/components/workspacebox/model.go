package workspacebox

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	workspaceName := strings.TrimPrefix(wb.Name, "special:")
	title := theme.Get().WorkspaceView.Title.Render("WS:" + workspaceName)

	var processText string
	if wb.WindowCount == 1 {
		processText = theme.Get().WorkspaceView.Details.Render(fmt.Sprintf("%d process", wb.WindowCount))
	} else {
		processText = theme.Get().WorkspaceView.Details.Render(fmt.Sprintf("%d processes", wb.WindowCount))
	}

	stats := theme.Get().WorkspaceView.Details.Render(fmt.Sprintf("CPU:%.1f%% | MEM %.1f%%", wb.CPUUsage, wb.MemUsage))

	content := lipgloss.JoinVertical(lipgloss.Center, title, processText, stats)

	return boxStyle.Render(content)
}

func (wb *WorkspaceBox) SetSelected(selected bool) {
	wb.IsSelected = selected
}
