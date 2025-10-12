package keymap

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyMap defines all the key bindings for the application
type KeyMap struct {
	Quit          key.Binding
	NavigateLeft  key.Binding
	NavigateRight key.Binding
	NavigateUp    key.Binding
	NavigateDown  key.Binding
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
	km.setNavigateLeftKeys("left", "h")
	km.setNavigateRightKeys("right", "l")
	km.setNavigateUpKeys("up", "k")
	km.setNavigateDownKeys("down", "j")
	return km
}

// GetHelpText returns formatted help text for navigation
func (km KeyMap) GetHelpText() string {
	leftKeys := km.NavigateLeft.Help().Key
	rightKeys := km.NavigateRight.Help().Key
	upKeys := km.NavigateUp.Help().Key
	downKeys := km.NavigateDown.Help().Key
	return leftKeys + "/" + rightKeys + "/" + upKeys + "/" + downKeys + ": navigate"
}

// HandleKeyMsg processes key messages for navigation
func (km KeyMap) HandleKeyMsg(msg tea.KeyMsg) (string, bool) {
	switch {
	case key.Matches(msg, km.Quit):
		return "quit", true
	case key.Matches(msg, km.NavigateLeft):
		return "navigate_left", true
	case key.Matches(msg, km.NavigateRight):
		return "navigate_right", true
	case key.Matches(msg, km.NavigateUp):
		return "navigate_up", true
	case key.Matches(msg, km.NavigateDown):
		return "navigate_down", true
	default:
		return "", false
	}
}

// Helper methods for initialization

func (km *KeyMap) setNavigateLeftKeys(keys ...string) {
	km.NavigateLeft = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "left"),
	)
}

func (km *KeyMap) setNavigateRightKeys(keys ...string) {
	km.NavigateRight = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "right"),
	)
}

func (km *KeyMap) setNavigateUpKeys(keys ...string) {
	km.NavigateUp = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "up"),
	)
}

func (km *KeyMap) setNavigateDownKeys(keys ...string) {
	km.NavigateDown = key.NewBinding(
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
