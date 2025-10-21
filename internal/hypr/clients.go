package hypr

import (
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/thiagokokada/hyprland-go"
)

type HyprlandClient struct {
	c *hyprland.RequestClient
}

func NewClient() *HyprlandClient {
	c := hyprland.MustClient()
	return &HyprlandClient{c: c}
}

func (c *HyprlandClient) ListClients() []Client {
	clients, err := c.c.Clients()
	if err != nil {
		logger.Log.Error("could not get hyprland clients: " + err.Error())
		return nil
	}
	res := []Client{}
	for _, client := range clients {
		res = append(res, Client{
			Class: client.Class,
			Workspace: Workspace{
				ID:   client.Workspace.Id,
				Name: client.Workspace.Name,
			},
			Monitor: client.Monitor,
			Title:   client.Title,
			PID:     client.Pid,
		})
	}
	return res
}
