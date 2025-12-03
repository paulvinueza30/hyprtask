package keymap

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens"
)

// KeyMap defines all the key bindings for the application
type KeyMap struct {
	Quit                            key.Binding
	NavigateLeft                    key.Binding
	NavigateRight                   key.Binding
	NavigateUp                      key.Binding
	NavigateDown                    key.Binding
	ScrollUp                        key.Binding
	ScrollDown                      key.Binding
	ChangeToAllProcsScreen          key.Binding
	ChangeToWorkspaceSelectorScreen key.Binding
	SelectWorkspace                 key.Binding
	SortKeyLeft                     key.Binding
	SortKeyRight                    key.Binding
	ToggleSortOrder                 key.Binding
	KillProcess                     key.Binding
	KillProcessForce                key.Binding
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
	km.setNavigateLeftKeys("left")
	km.setNavigateRightKeys("right")
	km.setNavigateUpKeys("up")
	km.setNavigateDownKeys("down")
	km.setScrollUpKeys("pgup")
	km.setScrollDownKeys("pgdown")
	km.setChangeToAllProcsScreenKeys("p", "ctrl+p")
	km.setChangeToWorkspaceSelectorScreenKeys("w", "ctrl+w")
	km.setSelectWorkspaceKeys("enter", "return")
	km.setSortKeyLeftKeys("[", "<")
	km.setSortKeyRightKeys("]", ">")
	km.setToggleSortOrderKeys("ctrl+o")
	km.setKillProcessKeys("x")
	km.setKillProcessForceKeys("X")
	return km
}

// GetHelpText returns formatted help text for the specified screen type
func (km KeyMap) GetHelpText(screenType screens.ScreenType) string {
	switch screenType {
	case screens.WorkspaceSelector:
		return km.getWorkspaceSelectorHelpText()
	case screens.ProcessList:
		return km.getProcessListHelpText()
	default:
		return "unknown screen type"
	}
}

func (km KeyMap) getWorkspaceSelectorHelpText() string {
	navigateKeys := fmt.Sprintf("%s/%s/%s/%s",
		km.NavigateLeft.Help().Key, km.NavigateRight.Help().Key, km.NavigateUp.Help().Key, km.NavigateDown.Help().Key)
	scrollKeys := fmt.Sprintf("%s/%s", km.ScrollUp.Help().Key, km.ScrollDown.Help().Key)

	return fmt.Sprintf("%s: navigate, %s: scroll, %s: view all processes, %s: select workspace, %s: quit",
		navigateKeys, scrollKeys, km.ChangeToAllProcsScreen.Help().Key, km.SelectWorkspace.Help().Key, km.Quit.Help().Key)
}

func (km KeyMap) getProcessListHelpText() string {
	return fmt.Sprintf("%s: change to workspace view, %s: sort key left, %s: sort key right, %s: toggle sort order, %s: kill process, %s: kill process force, %s: quit",
		km.ChangeToWorkspaceSelectorScreen.Help().Key, km.SortKeyLeft.Help().Key, km.SortKeyRight.Help().Key, km.ToggleSortOrder.Help().Key, km.KillProcess.Help().Key, km.KillProcessForce.Help().Key, km.Quit.Help().Key)
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
	case key.Matches(msg, km.ScrollUp):
		return "scroll_up", true
	case key.Matches(msg, km.ScrollDown):
		return "scroll_down", true
	case key.Matches(msg, km.ChangeToAllProcsScreen):
		return "change_to_all_procs_view", true
	case key.Matches(msg, km.ChangeToWorkspaceSelectorScreen):
		return "change_to_workspace_view", true
	case key.Matches(msg, km.SelectWorkspace):
		return "select_workspace", true
	case key.Matches(msg, km.SortKeyLeft):
		return "sort_key_left", true
	case key.Matches(msg, km.SortKeyRight):
		return "sort_key_right", true
	case key.Matches(msg, km.ToggleSortOrder):
		return "toggle_sort_order", true
	case key.Matches(msg, km.KillProcess):
		return "kill_process", true
	case key.Matches(msg, km.KillProcessForce):
		return "kill_process_force", true
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

func (km *KeyMap) setScrollUpKeys(keys ...string) {
	km.ScrollUp = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "scroll up"),
	)
}

func (km *KeyMap) setScrollDownKeys(keys ...string) {
	km.ScrollDown = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "scroll down"),
	)
}
func (km *KeyMap) setChangeToAllProcsScreenKeys(keys ...string) {
	km.ChangeToAllProcsScreen = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "view all processes"),
	)
}
func (km *KeyMap) setChangeToWorkspaceSelectorScreenKeys(keys ...string) {
	km.ChangeToWorkspaceSelectorScreen = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "change to workspace view"),
	)
}
func (km *KeyMap) setSelectWorkspaceKeys(keys ...string) {
	km.SelectWorkspace = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "select workspace"),
	)
}
func (km *KeyMap) setSortKeyLeftKeys(keys ...string) {
	km.SortKeyLeft = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "sort key left"),
	)
}
func (km *KeyMap) setSortKeyRightKeys(keys ...string) {
	km.SortKeyRight = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "sort key right"),
	)
}
func (km *KeyMap) setToggleSortOrderKeys(keys ...string) {
	km.ToggleSortOrder = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "toggle sort order"),
	)
}
func (km *KeyMap) setKillProcessKeys(keys ...string) {
	km.KillProcess = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "kill process (SIGTERM)"),
	)
}
func (km *KeyMap) setKillProcessForceKeys(keys ...string) {
	km.KillProcessForce = key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], "kill process force (SIGKILL)"),
	)
}