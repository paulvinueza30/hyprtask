package main

import (
	"github.com/paulvinueza30/hyprtask/internal/hypr"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/proc"
)

func main() {
	logger.Init()
	hypr.Init()
	systemMonitor , err := proc.Init()
	if err != nil{
		return
	}
	clients := hypr.ListClients()

	for _, c := range clients {
		_, err := systemMonitor.GetUsage(c.PID)
		if err != nil{
			logger.Log.Error("could not get usage for PID %d", c.PID, err)
		}

	}
}