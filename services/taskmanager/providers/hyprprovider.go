package providers

import "github.com/paulvinueza30/hyprtask/internal/hypr"

type HyprlandProvider struct {
	client *hypr.HyprlandClient
}

func NewHyprlandProvider() *HyprlandProvider {
	client := hypr.NewClient()
	return &HyprlandProvider{client: client}
}

func (p *HyprlandProvider) GetProcs() ([]Proc, error) {
	clients := p.client.ListClients()
	procs := make([]Proc, len(clients))
	for i, c := range clients {
		procs[i].PID = c.PID
		procs[i].Meta.Client = &c
	}
	return procs, nil
}
