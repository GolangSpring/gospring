package service

import . "go-spring/domains/repository"

type IInventoryService interface {
	GetInventory(id string) (*Inventory, error)
	CreateInventory(inventory *Inventory) (*Inventory, error)
	GetAllInventory() ([]*Inventory, error)
	UpdateInventory(id string, inventory *Inventory) (*Inventory, error)
	DeleteInventory(id string) (any, error)
}

type InventoryService struct{}

func (service *InventoryService) PostConstruct() {

}

func (service *InventoryService) GetInventory(id string) (*Inventory, error) {
	//TODO implement me
	panic("implement me")
}

func (service *InventoryService) CreateInventory(inventory *Inventory) (*Inventory, error) {
	//TODO implement me
	panic("implement me")
}

func (service *InventoryService) GetAllInventory() ([]*Inventory, error) {
	//TODO implement me
	panic("implement me")
}

func (service *InventoryService) UpdateInventory(id string, inventory *Inventory) (*Inventory, error) {
	//TODO implement me
	panic("implement me")
}

func (service *InventoryService) DeleteInventory(id string) (any, error) {
	//TODO implement me
	panic("implement me")
}

func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

var _ IInventoryService = (*InventoryService)(nil)
