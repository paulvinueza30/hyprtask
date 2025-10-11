package components

import (
	"github.com/charmbracelet/lipgloss"
)

var HeaderStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(0, 1) 