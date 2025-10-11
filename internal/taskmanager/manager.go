package taskmanager

import (
	"fmt"
	"sync"
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/proc"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager/providers"
)

type TaskManager struct {
	pollingInterval time.Duration
	systemMonitor   proc.SystemMonitor
	mode            Mode
	procProvider    providers.ProcessProvider

	activeProcesses map[int]TaskProcess // PID to task
	mu sync.RWMutex
	
	snapshotChan chan<- Snapshot
}

var procProviders = map[Mode]providers.ProcessProvider{
	Hypr: providers.NewHyprlandProvider(),
}

func NewTaskManager(mode string, pollInterval time.Duration, snapshotChan chan Snapshot) (*TaskManager, error) {
	m, ok := stringToMode[mode]
	if !ok {
		return nil, fmt.Errorf("invalid mode for task manager: %s", mode)
	}
	procProvider := procProviders[m]
	systemMonitor, err := proc.Init(time.Second * 4)
	if err != nil {
		return nil, err
	}
	
	activeProcesses := make(map[int]TaskProcess)
	return &TaskManager{mode: m, pollingInterval: pollInterval, systemMonitor: *systemMonitor, procProvider: procProvider, activeProcesses: activeProcesses , snapshotChan: snapshotChan}, nil
}

func (t *TaskManager) Start() {
	ticker := time.NewTicker(t.pollingInterval)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				go t.updateTaskProcesses()
			case <-done:
				logger.Log.Info("stopping task manager")
				return
			}
		}
	}()

	time.Sleep(30 * time.Second)
	done <- true
	time.Sleep(1 * time.Second)
}


func (t *TaskManager) updateTaskProcesses() {
	procs, err := t.procProvider.GetProcs()
	if err != nil {
		logger.Log.Error("could not get pids from proc provider", "err: ", err)
		return
	}

	t.deleteInactiveProcesses(procs)
	t.updateActiveProcesses(procs)

	t.sendSnapshot()
}

func (t *TaskManager) deleteInactiveProcesses(procs map[int]providers.Proc) {
	t.mu.Lock()
	defer t.mu.Unlock()	
	logger.Log.Info("active processes now: ", "active procs before", len(t.activeProcesses))
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

	logger.Log.Info("proc details", "active procs", len(t.activeProcesses), "procs deleted", deletedCount)
}


func (t *TaskManager) updateActiveProcesses(procs map[int]providers.Proc) {
	var wg sync.WaitGroup
	for pid, p := range procs {
		wg.Add(1)
		go func(pid int, meta providers.Meta) {
			defer wg.Done()
			m, err := t.systemMonitor.GetMetrics(pid)
			if err != nil {
				logger.Log.Error("could not get system metrics", "error", err)
				return
			}

			t.mu.Lock()
			defer t.mu.Unlock()
			t.activeProcesses[pid] = TaskProcess{
				PID:     pid,
				Meta:    meta,
				Metrics: *m,
			}
		}(pid, p.Meta)
	}
	wg.Wait()
}

func (t *TaskManager) makeSnapshot() Snapshot{
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	procs := make([]TaskProcess,0 , len(t.activeProcesses))
	for _, tp := range t.activeProcesses{
		procs = append(procs, tp)
	}
	return Snapshot{Processes: procs, Timestamp: time.Now()}
}

func (t *TaskManager) sendSnapshot(){
	snapshot := t.makeSnapshot()

	select {
	case t.snapshotChan <- snapshot:

	default:
		logger.Log.Warn("skipped snapshot send - channel full")
	}
}