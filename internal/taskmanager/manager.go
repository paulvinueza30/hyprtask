package taskmanager

import (
	"sync"
	"syscall"
	"time"

	"github.com/paulvinueza30/hyprtask/internal/hypr"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/metrics"
	"github.com/paulvinueza30/hyprtask/internal/procprovider"
)

type TaskManager struct {
	pollingInterval time.Duration
	systemMonitor   metrics.SystemMonitor
	procProvider    procprovider.ProcProvider
	
	hyprlandClient  *hypr.HyprlandClient

	activeProcesses map[int]TaskProcess // PID to task
	mu              sync.RWMutex

	snapshotChan   chan<- Snapshot
	taskActionChan <-chan TaskAction
}

const (
	DEBUG_MODE = false
)

func NewTaskManager(pollInterval time.Duration, snapshotChan chan Snapshot, taskActionChan chan TaskAction) (*TaskManager, error) {
	procProvider := procprovider.NewProcProvider()
	systemMonitor, err := metrics.NewSystemMonitor(time.Second * 4)
	hyprlandClient := hypr.NewHyprlandClient()
	if err != nil {
		return nil, err
	}

	activeProcesses := make(map[int]TaskProcess)
	return &TaskManager{
		pollingInterval: pollInterval, 
		systemMonitor: *systemMonitor, 
		procProvider: procProvider, 
		hyprlandClient: hyprlandClient, 
		activeProcesses: activeProcesses, 
		snapshotChan: snapshotChan, 
		taskActionChan: taskActionChan,
	}, nil
}

func (t *TaskManager) Start() {
	go t.handleTaskActions()
	
	// Trigger immediate update on startup to eliminate 5-second delay
	go t.updateTaskProcesses()
	
	ticker := time.NewTicker(t.pollingInterval)
	defer ticker.Stop()
	devTicker := time.NewTicker(30 * time.Second)
	defer devTicker.Stop()
	
	for {
		select {
		case <-ticker.C:
			go t.updateTaskProcesses()
		case <-devTicker.C:
			if DEBUG_MODE {
				return
			}
		}
	}
}

func (t *TaskManager) updateTaskProcesses() {
	procs, err := t.procProvider.GetProcs()
	if err != nil {
		logger.Log.Error("could not get pids from proc provider", "err: ", err)
		return
	}

	// Convert slice to map for easier lookup
	procMap := make(map[int]procprovider.Proc)
	for _, proc := range procs {
		procMap[proc.PID] = proc
	}

	t.deleteInactiveProcesses(procMap)
	
	// Progressive loading: send quick snapshot first, then enrich
	// Inject Hyprland metadata immediately (it's fast) so workspaces appear
	t.updateActiveProcessesQuick(procMap)
	t.injectHyprlandMeta()
	t.sendSnapshot()
	
	// Enrich with accurate metrics in background (this requires 4-second sleep)
	go t.enrichProcessesAccurate(procMap)
}

func (t *TaskManager) deleteInactiveProcesses(procs map[int]procprovider.Proc) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	// Take a snapshot of keys to avoid mutating map while iterating
	snapshot := make([]int, 0, len(t.activeProcesses))
	for pid := range t.activeProcesses {
		snapshot = append(snapshot, pid)
	}

	deletedCount := 0
	for _, pid := range snapshot {
		if _, ok := procs[pid]; !ok {
			deletedCount++
			delete(t.activeProcesses, pid)
		}
	}

	// Only log if there are significant changes
	if deletedCount > 0 {
		logger.Log.Info("processes removed", "deleted", deletedCount, "remaining", len(t.activeProcesses))
	}
}

// updateActiveProcessesQuick uses quick metrics for immediate display
func (t *TaskManager) updateActiveProcessesQuick(procs map[int]procprovider.Proc) {
	var wg sync.WaitGroup
	for pid, p := range procs {
		wg.Add(1)
		go func(pid int, proc procprovider.Proc) {
			defer wg.Done()
			m, err := t.systemMonitor.GetQuickMetrics(pid)
			if err != nil {
				logger.Log.Warn("could not get quick metrics giving default values", "error", err)
				m = &metrics.DEFAULT_METRICS
			}

			t.mu.Lock()
			defer t.mu.Unlock()
			t.activeProcesses[pid] = TaskProcess{
				PID:         pid,
				ProgramName: proc.ProgramName,
				User:        proc.User,
				CommandLine: proc.CommandLine,
				Metrics:     *m,
				Meta:        &Meta{},
			}
		}(pid, p)
	}
	wg.Wait()
}

// enrichProcessesAccurate updates processes with accurate metrics in background
func (t *TaskManager) enrichProcessesAccurate(procs map[int]procprovider.Proc) {
	// Update with accurate metrics (requires 4-second sleep per process)
	t.updateActiveProcessesAccurate(procs)
	
	// Send updated snapshot with accurate metrics
	t.sendSnapshot()
}

// updateActiveProcessesAccurate uses accurate metrics with sleep delay
func (t *TaskManager) updateActiveProcessesAccurate(procs map[int]procprovider.Proc) {
	var wg sync.WaitGroup
	for pid := range procs {
		wg.Add(1)
		go func(pid int) {
			defer wg.Done()
			m, err := t.systemMonitor.GetMetrics(pid)
			if err != nil {
				logger.Log.Warn("could not get system metrics giving default values", "error", err)
				return
			}

			t.mu.Lock()
			if taskProcess, ok := t.activeProcesses[pid]; ok {
				taskProcess.Metrics = *m
				t.activeProcesses[pid] = taskProcess
			}
			t.mu.Unlock()
		}(pid)
	}
	wg.Wait()
}

func (t *TaskManager) makeSnapshot() Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()

	procs := make([]TaskProcess, 0, len(t.activeProcesses))
	for _, tp := range t.activeProcesses {
		procs = append(procs, tp)
	}
	return Snapshot{Processes: procs, Timestamp: time.Now()}
}

func (t *TaskManager) sendSnapshot() {
	snapshot := t.makeSnapshot()

	select {
	case t.snapshotChan <- snapshot:

	default:
		logger.Log.Warn("skipped snapshot send - viewmodel is not ready")
	}
}
func (t *TaskManager) injectHyprlandMeta() {
	hyprlandMeta, err := t.hyprlandClient.GetHyprlandMeta()
	if err != nil {
		logger.Log.Error("could not get hyprland meta: " + err.Error())
		return
	}
	
	t.mu.Lock()
	defer t.mu.Unlock()
	
	for pid, meta := range hyprlandMeta {
		if taskProcess, ok := t.activeProcesses[pid]; ok {
			taskProcess.Meta.Hyprland = &meta
			t.activeProcesses[pid] = taskProcess
		} else {
			logger.Log.Warn("process not found in active processes", "pid", pid)
		}
	}
}

func (t *TaskManager) handleTaskActions() {
	for action := range t.taskActionChan {
		logger.Log.Info("Received task action", "action", action)
		t.handleTaskAction(action)
		t.sendSnapshot()
	}
}

func (t *TaskManager) handleTaskAction(action TaskAction) {
	switch action.Type {
	case TaskActionKill:
		t.handleKillProcess(action.Payload)
	}
}

func (t *TaskManager) handleKillProcess(payload KillProcessPayload) {
	signal := syscall.SIGTERM
	if payload.Force {
		signal = syscall.SIGKILL
	}
	
	err := syscall.Kill(payload.PID, signal)
	if err != nil {
		logger.Log.Error("Failed to kill process", "pid", payload.PID, "signal", signal, "error", err)
	} else {
		logger.Log.Info("Successfully killed process", "pid", payload.PID, "signal", signal)
		// Immediately remove from activeProcesses for instant UI feedback
		t.mu.Lock()
		delete(t.activeProcesses, payload.PID)
		t.mu.Unlock()
	}
}
