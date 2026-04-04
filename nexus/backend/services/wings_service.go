package services

import (
	"nexus/backend/models"
	"nexus/backend/wings"
)

type WingsService struct {
	client *wings.Client
}

func NewWingsService() *WingsService {
	return &WingsService{
		client: wings.NewClient(),
	}
}

func (w *WingsService) GetServerDetails(node models.Node, uuid string) (*wings.ServerDetails, error) {
	return w.client.GetServerDetails(node, uuid)
}

func (w *WingsService) CreateServer(node models.Node, payload wings.CreateServerPayload) error {
	return w.client.CreateServer(node, payload)
}

func (w *WingsService) DeleteServer(node models.Node, uuid string) error {
	return w.client.DeleteServer(node, uuid)
}

func (w *WingsService) SendPowerAction(node models.Node, uuid string, action string) error {
	return w.client.SendPowerAction(node, uuid, action)
}

func (w *WingsService) GetServerResources(node models.Node, uuid string) (*wings.ServerResources, error) {
	return w.client.GetServerResources(node, uuid)
}

func (w *WingsService) GetSystemInfo(node models.Node) (*wings.SystemInfo, error) {
	return w.client.GetSystemInfo(node)
}
