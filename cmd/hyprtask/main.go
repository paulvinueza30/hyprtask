package main

import (
	"github.com/paulvinueza30/hyprtask/internal/hypr"
	"github.com/paulvinueza30/hyprtask/internal/logger"
)


func main() {
	logger.Init()
	hypr.Init()
	clients := hypr.ListClients()
	for _, c := range clients{
		logger.Log.Debug(c.Title)
	}
}