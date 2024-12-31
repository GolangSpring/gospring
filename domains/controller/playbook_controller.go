package controller

import (
	"github.com/go-fuego/fuego"
	"go-spring/domains/repository"
	"go-spring/domains/service"
)

type PlaybookResources struct {
	PlaybookService service.PlaybookService
}

func (resource PlaybookResources) Routes(s *fuego.Server) {
	playbookGroup := fuego.Group(s, "/playbook")

	fuego.Get(playbookGroup, "/", resource.getAllPlaybook)
	fuego.Post(playbookGroup, "/", resource.postPlaybook)
	fuego.Get(playbookGroup, "/{id}", resource.getPlaybook)
	fuego.Put(playbookGroup, "/{id}", resource.putPlaybook)
	fuego.Delete(playbookGroup, "/{id}", resource.deletePlaybook)
}

func (resource PlaybookResources) getAllPlaybook(c fuego.ContextNoBody) ([]*repository.Playbook, error) {
	return resource.PlaybookService.GetAllPlaybook()
}

func (resource PlaybookResources) postPlaybook(c fuego.ContextWithBody[*repository.Playbook]) (*repository.Playbook, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	return resource.PlaybookService.CreatePlaybook(body)
}

func (resource PlaybookResources) getPlaybook(c fuego.ContextNoBody) (*repository.Playbook, error) {
	id := c.PathParam("id")

	return resource.PlaybookService.GetPlaybook(id)
}

func (resource PlaybookResources) putPlaybook(c fuego.ContextWithBody[*repository.Playbook]) (*repository.Playbook, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	return resource.PlaybookService.UpdatePlaybook(id, body)
}

func (resource PlaybookResources) deletePlaybook(c fuego.ContextNoBody) (any, error) {
	return nil, resource.PlaybookService.DeletePlaybook(c.PathParam("id"))
}
