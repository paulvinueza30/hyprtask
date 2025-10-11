package main

import (
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui"
)

func main() {
	logger.Init()
	
	snapshotChan := make(chan taskmanager.Snapshot, 1)
	defaultMode := "hypr"
	tm, err := taskmanager.NewTaskManager(defaultMode, 5 * time.Second, snapshotChan)
	if err != nil{
		return
	}
	vm := ui.NewViewModel(taskmanager.Mode(defaultMode), snapshotChan)
	go tm.Start()	
	vm.Start()

}