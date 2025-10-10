package main

import (
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/services/taskmanager"
)

func main() {
	logger.Init()
	
	taskmanager , err := taskmanager.NewTaskManager("hypr" , 5 * time.Second)
	if err != nil{
		return
	}
	taskmanager.Start()	
}