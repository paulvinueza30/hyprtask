package taskmanager

import (
	"time"

	"github.com/paulvinueza30/hyprtask/internal/proc"
	"github.com/paulvinueza30/hyprtask/internal/taskmanager/providers"
)

type Mode string

const (
	Hypr Mode = "hypr"
	All  Mode = "all"
)

var stringToMode = map[string]Mode{
	"all":  All,
	"hypr": Hypr,
}

type TaskProcess struct {
	PID     int
	Meta    providers.Meta
	Metrics proc.Metrics
}

type Snapshot struct {
	Processes []TaskProcess
	Timestamp time.Time
}

type TaskAction struct {
	Type    TaskActionType
	Payload interface{}
}

type TaskActionType int
