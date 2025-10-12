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
			v.updateSnapshot(snapshot)
			v.buildDisplayData()
		case action := <-v.actionChan:
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
	var wsDisplayData WorkspaceDisplayData
	if v.mode == taskmanager.Hypr {
		wsDisplayData = v.buildWorkspaceDisplayData(procs)
	}
	v.applyViewOptions(procs)

	v.displayData = DisplayData{All: procs, Hypr: wsDisplayData}

	// Send DisplayData to UI
	logger.Log.Info("About to send DisplayData", "displayData", v.displayData)
	v.sendDisplayData()
}

func (v *ViewModel) sendDisplayData() {
	select {
	case v.displayDataChan <- v.displayData:
		logger.Log.Info("DisplayData sent to UI successfully")
	default:
		logger.Log.Warn("DisplayData channel full, dropping message")
	}
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

func (v *ViewModel) buildWorkspaceDisplayData(procs []taskmanager.TaskProcess) WorkspaceDisplayData {

	workspaceToWorkspaceData := make(map[int]*WorkspaceData)

	workspaceCount := 0
	for _, proc := range procs {
		wID := proc.Meta.Client.Workspace.ID
		var wsData *WorkspaceData
		wsData, ok := workspaceToWorkspaceData[wID]
		if !ok {
			workspaceToWorkspaceData[wID] = &WorkspaceData{}
			wsData = workspaceToWorkspaceData[wID]
			wsData.WorkspaceName = proc.Meta.Client.Workspace.Name
			wsData.WorkspaceID = wID
		}
		wsData.TotalCPU += proc.Metrics.CPU
		wsData.TotalMEM += proc.Metrics.MEM
		wsData.ActiveProcs = append(wsData.ActiveProcs, proc)
	}
	workspaceCount = len(workspaceToWorkspaceData)
	workspaces := make([]*WorkspaceData, 0, workspaceCount)
	for _, wsData := range workspaceToWorkspaceData {
		wsData.ActiveProcsCount = len(wsData.ActiveProcs)
		v.applyViewOptions(wsData.ActiveProcs)
		workspaces = append(workspaces, wsData)

	}
	// Sort workspaces by name for consistent ordering
	slices.SortStableFunc(workspaces, func(a, b *WorkspaceData) int {
		return cmp.Compare(a.WorkspaceName, b.WorkspaceName)
	})
	return WorkspaceDisplayData{WorkspaceCount: workspaceCount, WorkspaceToProcs: workspaceToWorkspaceData, Workspaces: workspaces}
}
