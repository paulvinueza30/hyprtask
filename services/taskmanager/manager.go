package taskmanager

import (
	"fmt"
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/proc"
	"github.com/paulvinueza30/hyprtask/services/taskmanager/managers"
)

type TaskManager struct {
	pollingInterval time.Duration
	systemMonitor   proc.SystemMonitor
	mode            Mode
	procManager     managers.ProcessManager
}

var procManagers = map[Mode]managers.ProcessManager{
	Hypr: managers.NewHyprlandManager(),
}

func NewTaskManager(mode string, pollInterval time.Duration) (*TaskManager, error) {
	m, ok := stringToMode[mode]
	if !ok {
		return nil, fmt.Errorf("invalid mode for task manager: %s", mode)
	}
	procManager := procManagers[m]
	systemMonitor, err := proc.Init()
	if err != nil {
		return nil, err
	}
	return &TaskManager{mode: m, pollingInterval: pollInterval, systemMonitor: *systemMonitor, procManager: procManager}, nil
}

func (t *TaskManager) Start() {
	ticker := time.NewTicker(t.pollingInterval)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				pids, err := t.procManager.GetPIDs()
				if err != nil {
					logger.Log.Error("could not get pids from proc provider", "err: ", err)
					done <- true
					return
				}
				for _, p := range pids {
					_, err := t.systemMonitor.GetUsage(p)
					if err != nil {
						logger.Log.Error("could not get get process usage", "err: ", err)
					}
				}
			case <-done:
				logger.Log.Info("stopping task manager")
				return
			}
		}
	}()

	time.Sleep(10 * time.Second)
	done <- true
	time.Sleep(1 * time.Second)
}

// clients := hyprClient.ListClients()

// for _, c := range clients {
// 	_, err := systemMonitor.GetUsage(c.PID)
// 	if err != nil{
// 		logger.Log.Error("could not get usage for PID", "err: ", err)
// 		return nil, err
// 	}

// }
