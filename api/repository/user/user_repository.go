package repository

import (
	"context"
	model "online_fashion_shop/api/model/user"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) models.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.Create(user).Error
    if err!= nil {
        return nil, err
    }

    return user, nil
}
