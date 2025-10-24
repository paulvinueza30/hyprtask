package viewmodel

import (
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
)

type SortKey int

const (
	SortByNone SortKey = iota
	SortByPID
	SortByProgramName
	SortByUser
	SortByCPU
	SortByMEM
	// SortByWorkspace
)

var validSortKeys = map[SortKey]bool{
	SortByNone: true,
	SortByPID:  true,
	SortByProgramName: true,
	SortByUser: true,
	SortByCPU:  true,
	SortByMEM:  true,
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

type ViewAction struct {
	NewSortKey SortKey
	NewSortOrder SortOrder
}
type ViewOptions struct {
	SortKey SortKey
	SortOrder SortOrder
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
