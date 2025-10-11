package providers

import "github.com/paulvinueza30/hyprtask/internal/hypr"

type HyprlandProvider struct {
	client *hypr.HyprlandClient
}

func NewHyprlandProvider() *HyprlandProvider {
	client := hypr.NewClient()
	return &HyprlandProvider{client: client}
}

func (p *HyprlandProvider) GetProcs() (map[int]Proc, error) {
	clients := p.client.ListClients()
	procs := make(map[int]Proc, len(clients))
	for _, c := range clients {
		procs[c.PID] = Proc{Meta: Meta{Client: &c} , PID: c.PID}
	}
	return procs, nil
}
