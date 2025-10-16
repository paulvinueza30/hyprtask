package messages

import (
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

type ChangeScreenMsg struct {
	ScreenType screens.ScreenType
}

func NewChangeScreenMsg(screenType screens.ScreenType) ChangeScreenMsg {
	return ChangeScreenMsg{
		ScreenType: screenType,
	}
}
