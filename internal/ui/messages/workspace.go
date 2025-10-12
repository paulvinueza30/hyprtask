package messages

import "github.com/paulvinueza30/hyprtask/internal/viewmodel"

// WorkspaceDataMsg contains workspace data for the workspace screen
type WorkspaceDataMsg struct {
	Workspaces []*viewmodel.WorkspaceData
	Count      int
}

// NewWorkspaceDataMsg creates a new WorkspaceDataMsg from WorkspaceDisplayData
func NewWorkspaceDataMsg(hyprData viewmodel.WorkspaceDisplayData) WorkspaceDataMsg {
	return WorkspaceDataMsg{
		Workspaces: hyprData.Workspaces,
		Count:      hyprData.WorkspaceCount,
	}
}
