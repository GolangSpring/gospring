package repository

import (
	"context"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	Save(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id uint) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUserName(ctx context.Context, name string) (*User, error)
	UpdateUserPassword(ctx context.Context, user *User, password string) error
	ActivateUser(ctx context.Context, user *User) error
	UpdateUserRoles(ctx context.Context, user *User, roles []string) error
}

var _ IUserRepository = (*UserRepository)(nil)

type UserRepository struct {
	Engine *gorm.DB
}

func (repo *UserRepository) ActivateUser(ctx context.Context, user *User) error {
	return repo.Engine.WithContext(ctx).Model(user).Update("is_verified", true).Error
}

func (repo *UserRepository) createPreloadTx(ctx context.Context) *gorm.DB {
	return repo.Engine.WithContext(ctx)
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := repo.createPreloadTx(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindByUserName(ctx context.Context, name string) (*User, error) {
	var user User
	err := repo.createPreloadTx(ctx).First(&user, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindAll(ctx context.Context) ([]*User, error) {
	var users []*User
	err := repo.createPreloadTx(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := repo.createPreloadTx(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (repo *UserRepository) Save(ctx context.Context, user *User) error {
	return repo.Engine.WithContext(ctx).Save(user).Error
}

func (repo *UserRepository) DeleteByID(ctx context.Context, id uint) error {
	return repo.Engine.WithContext(ctx).Delete(&User{}, id).Error
}

func (repo *UserRepository) UpdateUserPassword(ctx context.Context, user *User, password string) error {
	return repo.Engine.WithContext(ctx).Model(user).Update("password", password).Error
}

func (repo *UserRepository) UpdateUserRoles(ctx context.Context, user *User, roles []string) error {
	return repo.Engine.WithContext(ctx).Model(user).Update("roles", pq.StringArray(roles)).Error
}

func NewUserRepository(engine *gorm.DB) *UserRepository {
	return &UserRepository{
		Engine: engine,
	}
}
