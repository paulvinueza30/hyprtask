package viewmodel

import (
	"cmp"
	"slices"
	"sync"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
)

type ViewModel struct {
	snapshotChan    <-chan taskmanager.Snapshot
	actionChan      <-chan ViewAction
	displayDataChan chan<- DisplayData

	mode taskmanager.Mode
	// viewLayout any // implement later need to think more
	viewOptions *ViewOptions

	currentSnapshot *taskmanager.Snapshot
	displayData     DisplayData

	mu sync.RWMutex
}

func NewViewModel(mode taskmanager.Mode, ssChan chan taskmanager.Snapshot, acChan chan ViewAction, ddChan chan DisplayData) *ViewModel {
	viewOptions := &ViewOptions{SortBy: SortByNone, SortOrder: OrderNone}
	return &ViewModel{
		snapshotChan:    ssChan,
		actionChan:      acChan,
		displayDataChan: ddChan,
		mode:            mode,
		viewOptions:     viewOptions,
		currentSnapshot: nil,
	}
}

func (v *ViewModel) Start() {
	for {
		select {
		case snapshot := <-v.snapshotChan:
			logger.Log.Info("snapshot received", "snapshot", snapshot)
			v.updateSnapshot(snapshot)
			v.buildDisplayData()
		case action := <-v.actionChan:
			logger.Log.Info("action recevied", "action", action)
			v.handleAction(action)
		}
	}
}

func (v *ViewModel) updateSnapshot(s taskmanager.Snapshot) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.currentSnapshot = &s
}

func (v *ViewModel) buildDisplayData() {
	v.mu.Lock()
	defer v.mu.Unlock()
	procs := v.currentSnapshot.Processes
	var group map[int][]taskmanager.TaskProcess
	if v.mode == taskmanager.Hypr {
		group = v.groupProcsByWorkspace(procs)
	}
	v.applyViewOptions(procs)

	v.displayData = DisplayData{All: procs, Hypr: group}
}
func (v *ViewModel) applyViewOptions(procs []taskmanager.TaskProcess) {
	if v.currentSnapshot == nil {
		return
	}
	viewOpts := v.viewOptions
	if viewOpts.SortBy == SortByNone {
		return
	}

	slices.SortStableFunc(procs, func(a, b taskmanager.TaskProcess) int {
		var less int
		switch viewOpts.SortBy {
		case SortByCPU:
			less = cmp.Compare(a.Metrics.CPU, b.Metrics.CPU)
		case SortByMEM:
			less = cmp.Compare(a.Metrics.MEM, b.Metrics.MEM)
		case SortByPID:
			less = cmp.Compare(a.PID, b.PID)
		}
		if viewOpts.SortOrder == OrderASC {
			return less
		}
		return -less
	})
}

func (v *ViewModel) groupProcsByWorkspace(procs []taskmanager.TaskProcess) map[int][]taskmanager.TaskProcess {
	group := make(map[int][]taskmanager.TaskProcess)
	for _, proc := range procs {
		wID := proc.Meta.Client.Workspace.ID
		group[wID] = append(group[wID], proc)
	}
	for _, procs := range group {
		v.applyViewOptions(procs)
	}
	return group
}
