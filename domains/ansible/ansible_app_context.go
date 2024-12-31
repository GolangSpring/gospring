package ansible

import (
	"go-spring/application"
	"go-spring/domains/controller"
	"go-spring/domains/service"
)

const ContextName = "AnsibleAppContext"

func MustNewAnsibleAppContext() *application.ApplicationContext {
	playbookService := service.NewPlaybookService()
	playbookResources := controller.NewPlaybookResources(playbookService)

	inventoryService := service.NewInventoryService()
	inventoryResources := controller.NewInventoryResources(inventoryService)

	return &application.ApplicationContext{
		Name: ContextName,
		Services: []application.IService{
			playbookService, inventoryService,
		},
		Controllers: []application.IController{
			playbookResources, inventoryResources,
		},
	}

}
