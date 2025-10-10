package providers

import "github.com/paulvinueza30/hyprtask/internal/hypr"

type Meta struct{
	Client *hypr.Client
}
type Proc struct{
	PID int
	Meta Meta
}
type ProcessProvider interface {
	GetProcs() ([]Proc, error)
}
