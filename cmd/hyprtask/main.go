package main

import (
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

func main() {
	logger.Init()

	snapshotChan := make(chan taskmanager.Snapshot, 1)
	defaultMode := "hypr"
	taskActionChan := make(chan taskmanager.TaskAction, 10)
	tm, err := taskmanager.NewTaskManager(defaultMode, 5*time.Second, snapshotChan, taskActionChan)
	if err != nil {
		return
	}
	viewActionChan := make(chan viewmodel.ViewAction, 10)
	displayDataChan := make(chan viewmodel.DisplayData, 1)
	vm := viewmodel.NewViewModel(taskmanager.Mode(defaultMode), snapshotChan, viewActionChan, displayDataChan)
	go tm.Start()
	vm.Start()

}
