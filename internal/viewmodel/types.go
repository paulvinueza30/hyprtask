package viewmodel

import "github.com/paulvinueza30/hyprtask/internal/taskmanager"

type SortKey int

const (
	SortByNone SortKey = iota
	SortByCPU
	SortByMEM
	SortByPID
	// SortByName
	// SortByWorkspace
)

var validSortKeys = map[SortKey]bool{
	SortByNone: true,
	SortByCPU:  true,
	SortByMEM:  true,
	SortByPID:  true,
}

type SortOrder int

const (
	OrderNone SortOrder = iota
	OrderASC
	OrderDESC
)

var validSortOrders = map[SortOrder]bool{
	OrderNone: true,
	OrderASC:  true,
	OrderDESC: true,
}

type ViewOptions struct {
	SortBy    SortKey
	SortOrder SortOrder
}

type ActionType int

const (
	ActionSetSortKey ActionType = iota
	ActionSetSortOrder
)

type ViewAction struct {
	Type    ActionType
	Payload interface{}
}

type WorkspaceData struct {
	ActiveProcs      []taskmanager.TaskProcess
	ActiveProcsCount int
	TotalCPU         float64
	TotalMEM         float64
	WorkspaceName    string
	WorkspaceID      int
}

type WorkspaceDisplayData struct {
	WorkspaceToProcs map[int]*WorkspaceData // workspace id -> procs in workspace
	Workspaces       []*WorkspaceData
	WorkspaceCount   int
}

type DisplayData struct {
	All  []taskmanager.TaskProcess
	Hypr WorkspaceDisplayData
}
