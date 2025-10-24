package main

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager"
	"github.com/paulvinueza30/hyprtask/internal/ui"
	"github.com/paulvinueza30/hyprtask/internal/viewmodel"
)

func main() {
	logger.Init()

	snapshotChan := make(chan taskmanager.Snapshot, 3)
	taskActionChan := make(chan taskmanager.TaskAction, 10)
	tm, err := taskmanager.NewTaskManager(5*time.Second, snapshotChan, taskActionChan)
	if err != nil {
		return
	}
	viewActionChan := make(chan viewmodel.ViewAction, 1)
	displayDataChan := make(chan viewmodel.DisplayData, 1)
	vm := viewmodel.NewViewModel(snapshotChan, viewActionChan, displayDataChan)
	go tm.Start()
	go vm.Start()

	m := ui.NewModel(displayDataChan, viewActionChan, taskActionChan)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		logger.Log.Error("could not start program", "error", err)
		os.Exit(1)
	}
}
