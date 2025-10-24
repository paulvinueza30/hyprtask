package taskmanager

import (
	"time"

	"github.com/paulvinueza30/hyprtask/internal/hypr"
	"github.com/paulvinueza30/hyprtask/internal/metrics"
)

type Meta struct {
	Hyprland *hypr.HyprlandMeta
}

type TaskProcess struct {
	PID         int
	ProgramName string
	User        string
	CommandLine string
	Meta        *Meta
	Metrics     metrics.Metrics
}

type Snapshot struct {
	Processes []TaskProcess
	Timestamp time.Time
}

type TaskAction struct {
	Type    TaskActionType
	Payload KillProcessPayload
}

type TaskActionType int

const (
	TaskActionKill TaskActionType = iota
)

type KillProcessPayload struct {
	PID   int
	Force bool // true for SIGKILL, false for SIGTERM
}
