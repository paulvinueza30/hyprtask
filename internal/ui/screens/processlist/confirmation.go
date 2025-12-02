package processlist

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmationScreen struct {
	show        bool
	pid         int
	processName string
	command     string
	force       bool
	width       int
	height      int
}

func NewConfirmationScreen() *ConfirmationScreen {
	return &ConfirmationScreen{
		show:        false,
		pid:         0,
		processName: "",
		command:     "",
		force:       false,
		width:       0,
		height:      0,
	}
}

func (c *ConfirmationScreen) Init() tea.Cmd {
	return nil
}

func (c *ConfirmationScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ShowConfirmationMsg:
		c.Show(msg.PID, msg.ProcessName, msg.Command, msg.Force)
		return c, nil
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		return c, nil
	}

	if !c.show {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			c.show = false
			return c, func() tea.Msg {
				return ConfirmKillMsg{
					PID:   c.pid,
					Force: c.force,
				}
			}
		case "esc":
			c.show = false
			return c, func() tea.Msg {
				return CancelKillMsg{}
			}
		}
	}

	return c, nil
}

func (c *ConfirmationScreen) View() string {
	if !c.show {
		return ""
	}

	killType := "SIGTERM"
	if c.force {
		killType = "SIGKILL (force)"
	}

	title := "Kill Process?"
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		MarginBottom(1)

	pidText := fmt.Sprintf("PID: %d", c.pid)
	processText := fmt.Sprintf("Process: %s", c.processName)
	commandText := fmt.Sprintf("Command: %s", c.command)
	signalText := fmt.Sprintf("Signal: %s", killType)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginBottom(1)

	confirmText := "Press Enter to confirm, Esc to cancel"
	confirmStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true).
		MarginTop(1)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(title),
		infoStyle.Render(pidText),
		infoStyle.Render(processText),
		infoStyle.Render(commandText),
		infoStyle.Render(signalText),
		confirmStyle.Render(confirmText),
	)

	dialogWidth := 60
	if c.width > 0 && c.width < dialogWidth {
		dialogWidth = c.width - 4
	}

	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(dialogWidth).
		Align(lipgloss.Center)

	dialog := dialogStyle.Render(content)

	centered := lipgloss.Place(
		c.width,
		c.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)

	return centered
}

func (c *ConfirmationScreen) Show(pid int, processName, command string, force bool) {
	c.show = true
	c.pid = pid
	c.processName = processName
	c.command = command
	c.force = force
}

func (c *ConfirmationScreen) Hide() {
	c.show = false
}

func (c *ConfirmationScreen) SetSize(width, height int) {
	c.width = width
	c.height = height
}

type ShowConfirmationMsg struct {
	PID         int
	ProcessName string
	Command     string
	Force       bool
}

type ConfirmKillMsg struct {
	PID   int
	Force bool
}

type CancelKillMsg struct{}
