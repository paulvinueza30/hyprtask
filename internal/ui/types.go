package ui

type SortKey int 

const (
	SortByNone SortKey = iota
	SortByCPU 
	SortByMEM 
	SortByPID
	// SortByName
	// SortByWorkspace
)

type SortOrder int 

const (
	OrderNone SortOrder = iota
	OrderASC
	OrderDESC
)
type ViewOptions struct{
	SortBy SortKey
	SortOrder SortOrder 
}