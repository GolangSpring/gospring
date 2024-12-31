package service

import (
	"go-spring/domains/repository"
)

type IPlaybookService interface {
	GetPlaybook(id string) (*repository.Playbook, error)
	CreatePlaybook(book *repository.Playbook) (*repository.Playbook, error)
	GetAllPlaybook() ([]*repository.Playbook, error)
	UpdatePlaybook(id string, book *repository.Playbook) (*repository.Playbook, error)
	DeletePlaybook(id string) error
}

type PlaybookService struct{}

func (service PlaybookService) GetPlaybook(id string) (*repository.Playbook, error) {
	//TODO implement me

	//TODO implement me
	//TODO implement me
	//TODO implement me

	panic("implement me")
}

func (service PlaybookService) CreatePlaybook(book *repository.Playbook) (*repository.Playbook, error) {
	panic("implement me")
}

func (service PlaybookService) GetAllPlaybook() ([]*repository.Playbook, error) {
	//TODO implement me
	panic("implement me")
}

func (service PlaybookService) UpdatePlaybook(id string, book *repository.Playbook) (*repository.Playbook, error) {
	panic("implement me")
}

func (service PlaybookService) DeletePlaybook(id string) error {
	//TODO implement me
	panic("implement me")
}

var _ IPlaybookService = (*PlaybookService)(nil)
