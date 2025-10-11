package ui

import (
	"cmp"
	"slices"
	"sync"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
)

type ViewModel struct {
	snapshotChan <-chan taskmanager.Snapshot 
	
	mode taskmanager.Mode
	// viewLayout any // implement later need to think more
	viewOptions *ViewOptions

	currentSnapshot *taskmanager.Snapshot
	displayData interface{}
	
	mu sync.RWMutex
}

func NewViewModel(mode taskmanager.Mode , snapshotChan chan taskmanager.Snapshot) *ViewModel {
	viewOptions := &ViewOptions{SortBy: SortByNone, SortOrder: OrderNone}
	return &ViewModel{
		snapshotChan: snapshotChan,
		mode: mode,
		viewOptions: viewOptions,
		currentSnapshot: nil,
	}
}

func (v *ViewModel) Start() {
    for snapshot := range v.snapshotChan {
        logger.Log.Info("snapshot received", "timestamp", snapshot.Timestamp)
		v.updateSnapshot(snapshot)
		v.buildDisplayData()
    }
}
func (v *ViewModel) updateSnapshot(s taskmanager.Snapshot){
	v.mu.Lock()
	defer v.mu.Unlock()
	v.currentSnapshot = &s
}

func (v *ViewModel) buildDisplayData() {
	v.mu.Lock()
	defer v.mu.Unlock()
	viewOpts := v.viewOptions
	if viewOpts.SortBy == SortByNone{
		return
	}
	procs := v.currentSnapshot.Processes
	
	slices.SortStableFunc(procs, func(a , b taskmanager.TaskProcess) int {
		var less int
		switch viewOpts.SortBy{
		case SortByCPU:
			less = cmp.Compare(a.Metrics.CPU, b.Metrics.CPU)
		case SortByMEM:
			less = cmp.Compare(a.Metrics.MEM, b.Metrics.MEM)
		case SortByPID:
			less = cmp.Compare(a.PID, b.PID)
		}
		if viewOpts.SortOrder == OrderASC{
			return less
		}
		return -less
	})
}