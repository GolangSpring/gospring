package controller

import (
	"github.com/go-fuego/fuego"
	"go-spring/domains/repository"
	"go-spring/domains/service"
)

type InventoryResources struct {
	// TODO add resources
	InventoryService service.IInventoryService
}

func NewInventoryResources(inventoryService service.IInventoryService) InventoryResources {
	return InventoryResources{
		InventoryService: inventoryService,
	}
}

func (resource InventoryResources) Routes(server *fuego.Server) {
	inventoryGroup := fuego.Group(server, "/inventory")

	fuego.Get(inventoryGroup, "/", resource.getAllInventory)
	fuego.Post(inventoryGroup, "/", resource.postInventory)
	fuego.Get(inventoryGroup, "/{id}", resource.getInventory)
	fuego.Put(inventoryGroup, "/{id}", resource.putInventory)
	fuego.Delete(inventoryGroup, "/{id}", resource.deleteInventory)
}

func (resource InventoryResources) getAllInventory(c fuego.ContextNoBody) ([]*repository.Inventory, error) {
	return resource.InventoryService.GetAllInventory()
}

func (resource InventoryResources) postInventory(c fuego.ContextWithBody[*repository.Inventory]) (*repository.Inventory, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	return resource.InventoryService.CreateInventory(body)
}

func (resource InventoryResources) getInventory(c fuego.ContextNoBody) (*repository.Inventory, error) {
	id := c.PathParam("id")

	return resource.InventoryService.GetInventory(id)
}

func (resource InventoryResources) putInventory(c fuego.ContextWithBody[*repository.Inventory]) (*repository.Inventory, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	return resource.InventoryService.UpdateInventory(id, body)
}

func (resource InventoryResources) deleteInventory(c fuego.ContextNoBody) (any, error) {
	return resource.InventoryService.DeleteInventory(c.PathParam("id"))
}
