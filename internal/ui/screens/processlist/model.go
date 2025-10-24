package processlist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/keymap"
	"github.com/paulvinueza30/hyprtask/internal/ui/messages"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

type ProcessList struct {
	stateManager *stateManager
	table        table.Model
	width        int
	height       int
}

func NewProcessList(procs []taskmanager.TaskProcess) *ProcessList {
	columns := []table.Column{
		{Title: "PID", Width: 8},
		{Title: "Program", Width: 20},
		{Title: "User", Width: 12},
		{Title: "Command", Width: 30},
		{Title: "CPU%", Width: 8},
		{Title: "Mem%", Width: 8},
	}
	
	rows := make([]table.Row, len(procs))
	for i, proc := range procs {
		rows[i] = table.Row{
			fmt.Sprintf("%d", proc.PID),
			proc.ProgramName,
			proc.User,
			proc.CommandLine,
			fmt.Sprintf("%.1f", proc.Metrics.CPU),
			fmt.Sprintf("%.1f", proc.Metrics.MEM),
		}
	}
	
	styles := table.DefaultStyles()
	styles.Header = styles.Header.Align(lipgloss.Center)
	styles.Cell = styles.Cell.Align(lipgloss.Center)
	
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithStyles(styles),
	)
	
	return &ProcessList{
		stateManager: newStateManager(procs, &t),
		table:        t,
	}
}

func (p *ProcessList) Init() tea.Cmd {
	return nil
}

func (p *ProcessList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case messages.ProcessListMsg:
		p.stateManager.setState(typedMsg)
		p.updateTableWithProcesses(p.stateManager.getProcs())
		return p, nil
	case tea.WindowSizeMsg:
		p.handleWindowSize(typedMsg)
		return p, nil
	case tea.KeyMsg:
		updatedTable, cmd := p.table.Update(msg)
		p.table = updatedTable
		if cmd != nil {
			return p, cmd
		}
		return p, p.stateManager.handleKeyMsg(typedMsg)
	}
	return p, nil
}

func (p *ProcessList) View() string {
	wsName := p.stateManager.getWorkspaceName()
	wsNameStr := "all"
	if wsName != nil {
		wsNameStr = *wsName
	}

	title := fmt.Sprintf("Process List for workspace %s", wsNameStr)
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render(title)

	tableView := p.table.View()
	// TODO: Add scratchpad for full help text
	tableHelp := "Table Help: " + p.table.HelpView()
	
	// Debug: Show current sort options
	sortKey := p.stateManager.state.sortOptions.key
	sortOrder := p.stateManager.state.sortOptions.order
	sortKeyStr := getSortKeyString(sortKey)
	sortOrderStr := getSortOrderString(sortOrder)
	debugInfo := fmt.Sprintf(" | Sort: %s %s", sortKeyStr, sortOrderStr)
	tableHelp += debugInfo
	
	instructions := keymap.Get().GetHelpText(screens.ProcessList)

	centeredHeader := lipgloss.PlaceHorizontal(p.width, lipgloss.Center, header)
	centeredTable := lipgloss.PlaceHorizontal(p.width, lipgloss.Center, tableView)
	centeredHelp := lipgloss.PlaceHorizontal(p.width, lipgloss.Center, tableHelp)
	centeredInstructions := lipgloss.PlaceHorizontal(p.width, lipgloss.Center, instructions)

	var marginTop, marginBottom int
	if p.height > 20 {
		marginTop = 2
		marginBottom = 1
	}

	headerStyled := lipgloss.NewStyle().MarginTop(marginTop).Render(centeredHeader)
	tableStyled := lipgloss.NewStyle().MarginTop(marginTop).MarginBottom(marginBottom).Render(centeredTable)
	helpStyled := lipgloss.NewStyle().MarginBottom(marginBottom).Render(centeredHelp)
	instructionsStyled := lipgloss.NewStyle().Render(centeredInstructions)

	return lipgloss.JoinVertical(lipgloss.Center, headerStyled, tableStyled, helpStyled, instructionsStyled)
}

func (p *ProcessList) updateTableWithProcesses(procs []taskmanager.TaskProcess) {
	rows := make([]table.Row, len(procs))
	for i, proc := range procs {
		rows[i] = table.Row{
			fmt.Sprintf("%d", proc.PID),
			proc.ProgramName,
			proc.User,
			proc.CommandLine,
			fmt.Sprintf("%.1f", proc.Metrics.CPU),
			fmt.Sprintf("%.1f", proc.Metrics.MEM),
		}
	}
	
	p.table.SetRows(rows)
	p.table.Focus()
}

func (p *ProcessList) handleWindowSize(msg tea.WindowSizeMsg) {
	p.width = msg.Width
	p.height = msg.Height
	
	tableHeight := 10
	if msg.Height > 0 {
		tableHeight = msg.Height - 10
		if tableHeight < 5 {
			tableHeight = 5
		}
	}
	
	p.table.SetHeight(tableHeight)
}

// Helper functions for debug display
func getSortKeyString(key viewmodel.SortKey) string {
	switch key {
	case viewmodel.SortByNone:
		return "None"
	case viewmodel.SortByPID:
		return "PID"
	case viewmodel.SortByProgramName:
		return "Program"
	case viewmodel.SortByCPU:
		return "CPU"
	case viewmodel.SortByMEM:
		return "MEM"
	default:
		return "Unknown"
	}
}

func getSortOrderString(order viewmodel.SortOrder) string {
	switch order {
	case viewmodel.OrderNone:
		return "None"
	case viewmodel.OrderASC:
		return "ASC"
	case viewmodel.OrderDESC:
		return "DESC"
	default:
		return "Unknown"
	}
}
