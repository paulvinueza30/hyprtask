package keymap

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyMap defines all the key bindings for the application
type KeyMap struct {
	Quit                   key.Binding
	NavigateWorkspaceLeft  key.Binding
	NavigateWorkspaceRight key.Binding
	NavigateProcessUp      key.Binding
	NavigateProcessDown    key.Binding
}

var keyMap KeyMap

func Init() {
	keyMap = NewDefaultKeyMap()
}

func Get() KeyMap {
	return keyMap
}

func NewDefaultKeyMap() KeyMap {
	km := KeyMap{}
	km.setQuitKeys("q", "ctrl+c")
	km.setNavigateWorkspaceLeftKeys("h", "left")
	km.setNavigateWorkspaceRightKeys("l", "right")
	km.setNavigateProcessUpKeys("up", "k")
	km.setNavigateProcessDownKeys("down", "j")
	return km
}

// GetHelpText returns formatted help text for workspace navigation
func (km KeyMap) GetWorkspaceHelpText() string {
	leftKeys := km.NavigateWorkspaceLeft.Help().Key
	rightKeys := km.NavigateWorkspaceRight.Help().Key
	return leftKeys + "/" + rightKeys + ": navigate"
}

// GetHelpText returns formatted help text for process navigation
func (km KeyMap) GetProcessHelpText() string {
	upKeys := km.NavigateProcessUp.Help().Key
	downKeys := km.NavigateProcessDown.Help().Key
	return upKeys + "/" + downKeys + ": navigate"
}

// HandleWorkspaceKeyMsg processes key messages for workspace navigation
func (km KeyMap) HandleWorkspaceKeyMsg(msg tea.KeyMsg) (string, bool) {
	switch {
	case key.Matches(msg, km.Quit):
		return "quit", true
	case key.Matches(msg, km.NavigateWorkspaceLeft):
		return "navigate_left", true
	case key.Matches(msg, km.NavigateWorkspaceRight):
		return "navigate_right", true
	default:
		return "", false
	}
}

// HandleProcessKeyMsg processes key messages for process navigation
func (km KeyMap) HandleProcessKeyMsg(msg tea.KeyMsg) (string, bool) {
	switch {
	case key.Matches(msg, km.Quit):
		return "quit", true
	case key.Matches(msg, km.NavigateProcessUp):
		return "navigate_up", true
	case key.Matches(msg, km.NavigateProcessDown):
		return "navigate_down", true
	default:
		return "", false
	}
}

// Helper methods for initialization

func (km *KeyMap) setNavigateWorkspaceLeftKeys(keys ...string) {
	km.NavigateWorkspaceLeft = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "left"),
	)
}

func (km *KeyMap) setNavigateWorkspaceRightKeys(keys ...string) {
	km.NavigateWorkspaceRight = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "right"),
	)
}

func (km *KeyMap) setNavigateProcessUpKeys(keys ...string) {
	km.NavigateProcessUp = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "up"),
	)
}

func (km *KeyMap) setNavigateProcessDownKeys(keys ...string) {
	km.NavigateProcessDown = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "down"),
	)
}

func (km *KeyMap) setQuitKeys(keys ...string) {
	km.Quit = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "quit"),
	)
}
