package service

import (
	"context"
	"go-spring/application"
	. "go-spring/pkg/security/repository"
)

type IUserService interface {
	AddUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUserName(ctx context.Context, name string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	UpdateUserRoles(ctx context.Context, userID uint, roles []string) (*User, error)
}

var _ application.IService = (*UserService)(nil)
var _ IUserService = (*UserService)(nil)

type UserService struct {
	Repository IUserRepository
}

func (service *UserService) FindByID(ctx context.Context, id uint) (*User, error) {
	return service.Repository.FindByID(ctx, id)
}

func (service *UserService) FindByUserName(ctx context.Context, name string) (*User, error) {
	return service.Repository.FindByUserName(ctx, name)
}

func (service *UserService) FindByEmail(ctx context.Context, email string) (*User, error) {
	return service.Repository.FindByEmail(ctx, email)
}

func (service *UserService) AddUser(ctx context.Context, user *User) error {
	userFound, err := service.Repository.FindByID(ctx, user.ID)
	if err == nil && userFound != nil {
		return UserExists
	}

	if err = service.Repository.Save(ctx, user); err != nil {
		return err
	}
	return nil
}

func (service *UserService) UpdateUserRoles(ctx context.Context, userID uint, roles []string) (*User, error) {
	user, err := service.Repository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Roles = roles
	return user, service.Repository.UpdateUserRoles(ctx, user, roles)
}

func (service *UserService) PostConstruct() {}

func NewUserService(repository IUserRepository) *UserService {
	return &UserService{
		Repository: repository,
	}
}
