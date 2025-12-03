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
	confirmation  *ConfirmationScreen
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
		confirmation: NewConfirmationScreen(),
	}
}

func (p *ProcessList) Init() tea.Cmd {
	return nil
}

func (p *ProcessList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case ConfirmKillMsg:
		return p, func() tea.Msg {
			return messages.NewKillProcessMsg(typedMsg.PID, typedMsg.Force)
		}
	case CancelKillMsg:
		return p, nil
	}
	
	if p.confirmation.show {
		updatedConfirmation, cmd := p.confirmation.Update(msg)
		p.confirmation = updatedConfirmation.(*ConfirmationScreen)
		if cmd != nil {
			return p, cmd
		}
	}

	switch typedMsg := msg.(type) {
	case messages.ProcessListMsg:
		p.stateManager.setState(typedMsg)
		p.updateTableWithProcesses(p.stateManager.getProcs())
		return p, nil
	case ShowConfirmationMsg:
		p.confirmation.SetSize(p.width, p.height)
		updatedConfirmation, cmd := p.confirmation.Update(msg)
		p.confirmation = updatedConfirmation.(*ConfirmationScreen)
		return p, cmd
	case tea.WindowSizeMsg:
		p.handleWindowSize(typedMsg)
		updatedConfirmation, _ := p.confirmation.Update(msg)
		p.confirmation = updatedConfirmation.(*ConfirmationScreen)
		return p, nil
	case tea.KeyMsg:
		updatedTable, cmd := p.table.Update(msg)
		p.table = updatedTable
		p.stateManager.updateTable(&p.table)
		if cmd != nil {
			return p, cmd
		}
		return p, p.stateManager.handleKeyMsg(typedMsg)
	}
	return p, nil
}

func (p *ProcessList) View() string {
	wsName := p.stateManager.getWorkspaceName()
	wsNameStr := "all processes"
	if wsName != nil {
		wsNameStr = "workspace " + *wsName
	}

	title := fmt.Sprintf("Process List for %s", wsNameStr)
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render(title)

	p.updateColumnHeaders()
	
	tableView := p.table.View()
	tableHelp := "Table Help: " + p.table.HelpView()
	
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

	processListView := lipgloss.JoinVertical(lipgloss.Center, headerStyled, tableStyled, helpStyled, instructionsStyled)

	if p.confirmation.show {
		return p.confirmation.View()
	}

	return processListView
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
	p.updateColumnHeaders()
	p.table.Focus()
	p.stateManager.updateTable(&p.table)
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

func (p *ProcessList) updateColumnHeaders() {
	sortKey := p.stateManager.state.sortOptions.key
	sortOrder := p.stateManager.state.sortOptions.order
	
	currentColumns := p.table.Columns()
	baseTitles := []string{"PID", "Program", "User", "Command", "CPU%", "Mem%"}
	
	var arrow string
	switch sortOrder {
	case viewmodel.OrderASC:
		arrow = " ↑"
	case viewmodel.OrderDESC:
		arrow = " ↓"
	default:
		arrow = ""
	}
	
	columnIndex := p.getColumnIndexForSortKey(sortKey)
	
	newColumns := make([]table.Column, len(currentColumns))
	for i := range currentColumns {
		if i < len(baseTitles) {
			title := baseTitles[i]
			if sortKey != viewmodel.SortByNone && i == columnIndex {
				title = baseTitles[i] + arrow
			}
			newColumns[i] = table.Column{
				Title: title,
				Width: currentColumns[i].Width,
			}
		} else {
			newColumns[i] = currentColumns[i]
		}
	}
	
	p.table.SetColumns(newColumns)
}

func (p *ProcessList) getColumnIndexForSortKey(sortKey viewmodel.SortKey) int {
	switch sortKey {
	case viewmodel.SortByPID:
		return 0
	case viewmodel.SortByProgramName:
		return 1
	case viewmodel.SortByUser:
		return 2
	case viewmodel.SortByCPU:
		return 4
	case viewmodel.SortByMEM:
		return 5
	default:
		return -1
	}
}
