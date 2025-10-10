package taskmanager

import (
	"fmt"
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/proc"
	"github.com/paulvinueza30/hyprtask/services/taskmanager/providers"
)

type TaskManager struct {
	pollingInterval time.Duration
	systemMonitor   proc.SystemMonitor
	mode            Mode
	procProvider    providers.ProcessProvider
	
	activeProcesses map[int]*TaskProcess // PID to task
}

var procProviders = map[Mode]providers.ProcessProvider{
	Hypr: providers.NewHyprlandProvider(),
}

func NewTaskManager(mode string, pollInterval time.Duration) (*TaskManager, error) {
	m, ok := stringToMode[mode]
	if !ok {
		return nil, fmt.Errorf("invalid mode for task manager: %s", mode)
	}
	procProvider := procProviders[m]
	systemMonitor, err := proc.Init(time.Second * 4)
	if err != nil {
		return nil, err
	}
	
	activeProcesses := make(map[int]*TaskProcess)
	return &TaskManager{mode: m, pollingInterval: pollInterval, systemMonitor: *systemMonitor, procProvider: procProvider, activeProcesses: activeProcesses}, nil
}

func (t *TaskManager) Start() {
	ticker := time.NewTicker(t.pollingInterval)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := t.poll(); err != nil {
					return
				}
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

func (t *TaskManager) poll() error {
	procs, err := t.procProvider.GetProcs()
	if err != nil {
		logger.Log.Error("could not get pids from proc provider", "err: ", err)
		return err
	}
	for _, p := range procs{
		t.activeProcesses[p.PID] = &TaskProcess{PID: p.PID , Meta: p.Meta,}
		go func(pid int) {
			m, err := t.systemMonitor.GetMetrics(pid)
			if err != nil {
				logger.Log.Error("could not get get process usage", "err: ", err)
				return
			}
			proc := t.activeProcesses[pid]
			proc.Metrics = *m
		}(p.PID)
	}
	return nil
}
