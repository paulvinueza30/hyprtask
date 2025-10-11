package main

import (
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
)

func main() {
	logger.Init()
	
	snapshotChan := make(chan taskmanager.Snapshot, 1)
	tm, err := taskmanager.NewTaskManager("hypr" , 5 * time.Second, snapshotChan)
	if err != nil{
		return
	}
	tm.Start()	
}