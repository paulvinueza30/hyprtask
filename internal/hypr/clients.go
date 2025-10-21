package hypr

import (
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/thiagokokada/hyprland-go"
)

type HyprlandClient struct {
	c *hyprland.RequestClient
}

func NewHyprlandClient() *HyprlandClient {
	c := hyprland.MustClient()
	return &HyprlandClient{c: c}
}

func (c *HyprlandClient) GetHyprlandMeta() (map[int]HyprlandMeta, error) {
	meta := make(map[int]HyprlandMeta)
	clients, err := c.c.Clients()
	if err != nil {
		logger.Log.Error("could not get hyprland clients: " + err.Error())
		return nil, err
	}
	for _, client := range clients {
		meta[client.Pid] = HyprlandMeta{
			Workspace: Workspace{
				ID:   client.Workspace.Id,
				Name: client.Workspace.Name,
			},
			Monitor: client.Monitor,
			Title:   client.Title,
			Class:   client.Class,
			PID:     client.Pid,
		}
	}
	return meta, nil
}