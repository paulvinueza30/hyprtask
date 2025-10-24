package messages

import (
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui/screens"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

type WorkspaceDataMsg struct {
	Workspaces []*viewmodel.WorkspaceData
	Count      int
}

func NewWorkspaceDataMsg(hyprData viewmodel.WorkspaceDisplayData) WorkspaceDataMsg {
	return WorkspaceDataMsg{
		Workspaces: hyprData.Workspaces,
		Count:      hyprData.WorkspaceCount,
	}
}

type ChangeScreenMsg[T ScreenMsg] struct {
	ScreenType screens.ScreenType
	ScreenMsg  T
}

type ScreenMsg interface {
	ProcessListMsg | WorkspaceListMsg
}

// Screen-specific message types

type ProcessListMsg struct {
	WorkspaceID   *int                      // nil = all processes, &workspaceID = specific workspace (for logic)
	WorkspaceName *string                   // nil = all processes, &workspaceName = specific workspace (for display)
	Processes     []taskmanager.TaskProcess // actual process data
}

type WorkspaceListMsg struct {
	// Future workspace-specific data
}

func NewChangeScreenMsg[T ScreenMsg](screenType screens.ScreenType, screenMsg T) ChangeScreenMsg[T] {
	return ChangeScreenMsg[T]{
		ScreenType: screenType,
		ScreenMsg:  screenMsg,
	}
}

func NewProcessListMsg(workspaceID *int, workspaceName *string) ProcessListMsg {
	return ProcessListMsg{
		WorkspaceID:   workspaceID,
		WorkspaceName: workspaceName,
		Processes:     []taskmanager.TaskProcess{},
	}
}

func NewAllProcessesMsg() ProcessListMsg {
	return NewProcessListMsg(nil, nil)
}

func NewWorkspaceProcessesMsg(workspaceID int, workspaceName string) ProcessListMsg {
	return NewProcessListMsg(&workspaceID, &workspaceName)
}

type ChangeSortOptionMsg struct{
	Key viewmodel.SortKey
	Order viewmodel.SortOrder
}

func NewChangeSortOptionMsg(key viewmodel.SortKey, order viewmodel.SortOrder) ChangeSortOptionMsg {
	return ChangeSortOptionMsg{
		Key: key,
		Order: order,
	}
}

type KillProcessMsg struct {
	PID   int
	Force bool // true for SIGKILL, false for SIGTERM
}

func NewKillProcessMsg(pid int, force bool) KillProcessMsg {
	return KillProcessMsg{
		PID:   pid,
		Force: force,
	}
}