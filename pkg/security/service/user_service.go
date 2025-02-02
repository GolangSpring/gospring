package service

import (
	"context"
	"github.com/GolangSpring/gospring/application"
	. "github.com/GolangSpring/gospring/pkg/security/repository"
)

type IUserService interface {
	IUserRepository
	UpdateUserRolesByUserID(ctx context.Context, userID uint, roles []string) (*User, error)
	AddUser(ctx context.Context, user *User) error
	ResetUserPassword(ctx context.Context, user *User, password string) error
}

func NewUserService(repository IUserRepository) *UserService {
	return &UserService{
		IUserRepository: repository,
	}
}

var _ application.IService = (*UserService)(nil)
var _ IUserService = (*UserService)(nil)

type UserService struct {
	IUserRepository
}

func (service *UserService) ResetUserPassword(ctx context.Context, user *User, password string) error {
	return service.IUserRepository.UpdateUserPassword(ctx, user, password)
}

func (service *UserService) UpdateUserPasswordByUserID(ctx context.Context, userID uint, password string) error {
	user, err := service.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Password = password
	return service.IUserRepository.UpdateUserPassword(ctx, user, password)
}

func (service *UserService) AddUser(ctx context.Context, user *User) error {
	userFound, err := service.FindByID(ctx, user.ID)
	if err == nil && userFound != nil {
		return UserExists
	}

	if err = service.Save(ctx, user); err != nil {
		return err
	}
	return nil
}

func (service *UserService) UpdateUserRolesByUserID(ctx context.Context, userID uint, roles []string) (*User, error) {
	user, err := service.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Roles = roles
	return user, service.UpdateUserRoles(ctx, user, roles)
}

func (service *UserService) PostConstruct() {}
