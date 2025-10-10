package managers

type ProcessManager interface {
	GetPIDs() ([]int, error)
}
