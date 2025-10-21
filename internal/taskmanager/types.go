package taskmanager

import (
	"time"

	"github.com/paulvinueza30/hyprtask/internal/hypr"
	"github.com/paulvinueza30/hyprtask/internal/metrics"
)

type Meta struct {
	HyprlandMeta *hypr.HyprlandMeta
}

type TaskProcess struct {
	PID     int
	Meta    *Meta
	Metrics metrics.Metrics
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
