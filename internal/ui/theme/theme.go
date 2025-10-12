package theme

import "github.com/charmbracelet/lipgloss"

// Theme struct and sub-structs remain the same
type Theme struct {
    Header        lipgloss.Style
    ViewModel     ViewModelTheme
    WorkspaceView WorkspaceTheme
    ProcessView   ProcessListTheme
    UsageBars     UsageBarTheme
    Help          HelpTheme
    Footer        lipgloss.Style
}

type WorkspaceTheme struct {
    Box          lipgloss.Style
    SelectedBox  lipgloss.Style
    Title        lipgloss.Style
    Details      lipgloss.Style
}

type ProcessListTheme struct {
    HeaderText  lipgloss.Style
    Row         lipgloss.Style
    SelectedRow lipgloss.Style
}
type ViewModelTheme struct {
    Title lipgloss.Style
}

type UsageBarTheme struct {
    Text      lipgloss.Style
    Bar       lipgloss.Style
    BarFill   string
}

type HelpTheme struct {
    Key   lipgloss.Style
    Desc  lipgloss.Style
}

var theme Theme

func Init(){
    theme = NewDefaultTheme()
}

func Get() Theme{
    return theme 
}

func NewDefaultTheme() Theme {
    return Theme{
        Header:        buildHeaderTheme(DefaultColorAccent, DefaultColorForeground),
        ViewModel:     buildViewModelTheme(DefaultColorAccent, DefaultColorForeground),
        Footer:        buildFooterTheme(),
        WorkspaceView: buildWorkspaceTheme(DefaultColorBorder, DefaultColorAccent, DefaultColorForeground, DefaultColorMutedText),
        ProcessView:   buildProcessListTheme(DefaultColorAccent, DefaultColorForeground),
        UsageBars:     buildUsageBarTheme(DefaultColorMutedText, DefaultColorSuccess),
        Help:          buildHelpTheme(DefaultColorAccent, DefaultColorMutedText),
    }
}


func buildHeaderTheme(bg, fg string) lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(bg)).
        Foreground(lipgloss.Color(fg)).
        Bold(true).
        Padding(0, 1)
}

func buildFooterTheme() lipgloss.Style {
    return lipgloss.NewStyle().MarginTop(1)
}

func buildWorkspaceTheme(border, accent, fg, muted string) WorkspaceTheme {
    baseBox := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color(border)).
        Padding(0, 1).
        Width(24).
        Height(5)

    return WorkspaceTheme{
        Box:         baseBox,
        SelectedBox: baseBox.BorderForeground(lipgloss.Color(accent)),
        Title:       lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(fg)),
        Details:     lipgloss.NewStyle().Foreground(lipgloss.Color(muted)),
    }
}

func buildProcessListTheme(accent, fg string) ProcessListTheme {
    return ProcessListTheme{
        HeaderText: lipgloss.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color(accent)),
        Row: lipgloss.NewStyle().
            Padding(0, 1),
        SelectedRow: lipgloss.NewStyle().
            Padding(0, 1).
            Background(lipgloss.Color(accent)).
            Foreground(lipgloss.Color(fg)),
    }
}

func buildUsageBarTheme(muted, success string) UsageBarTheme {
    return UsageBarTheme{
        Text: lipgloss.NewStyle().
            Bold(true),
        Bar: lipgloss.NewStyle().
            Background(lipgloss.Color(muted)).
            Foreground(lipgloss.Color(success)),
        BarFill: "█",
    }
}

func buildHelpTheme(accent, muted string) HelpTheme {
    return HelpTheme{
        Key: lipgloss.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color(accent)),
        Desc: lipgloss.NewStyle().
            Foreground(lipgloss.Color(muted)),
    }
}

func buildViewModelTheme(accent, fg string) ViewModelTheme {
    return ViewModelTheme{
        Title: lipgloss.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color(accent)),
    }
}