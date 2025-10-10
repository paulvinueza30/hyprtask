package managers

import "github.com/paulvinueza30/hyprtask/internal/hypr"

type HyprlandManager struct {
	client *hypr.HyprlandClient
}

func NewHyprlandManager() *HyprlandManager {
	client := hypr.NewClient()
	return &HyprlandManager{client: client}
}

func (p *HyprlandManager) GetPIDs() ([]int, error) {
	clients := p.client.ListClients()

	pids := make([]int, len(clients))
	for i, c := range clients {
		pids[i] = c.PID
	}
	return pids, nil
}
