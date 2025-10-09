package main

import (
	"github.com/paulvinueza30/hyprtask/internal/hypr"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/proc"
)

func main() {
	logger.Init()
	hypr.Init()

	clients := hypr.ListClients()

	for _, c := range clients {
		proc.GetStats(c.PID)
	}
}
