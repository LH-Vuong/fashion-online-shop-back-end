package repository

import (
	"context"
	model "online_fashion_shop/api/model/user"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
	CreateUserVerify(context.Context, *model.UserVerify) (*model.UserVerify, error)
	GetVerifyByUniqueToken(context.Context, string) (*model.UserVerify, error)
	GetUserByEmail(context.Context, string) (*model.User, error)
	GetUserById(context.Context, string) (*model.User, error)
	UpdateUserVerify(context.Context, *model.UserVerify) error
	DeleteUserVerify(context.Context, *model.UserVerify) error
	DeleteUser(context.Context, *model.User) (string, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *userRepository) CreateUserVerify(ctx context.Context, userVerify *model.UserVerify) (*model.UserVerify, error) {
	result := r.db.WithContext(ctx).Create(&userVerify)
	if result.Error != nil {
		return nil, result.Error
	}

	return userVerify, nil
}

func (r *userRepository) GetVerifyByUniqueToken(ctx context.Context, uniqueToken string) (*model.UserVerify, error) {
	var userVerify *model.UserVerify
	err := r.db.WithContext(ctx).Where("unique_token = ?", uniqueToken).First(&userVerify).Error

	if err != nil {
		return nil, err
	}

	return userVerify, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {
	var user *model.User

	if err := r.db.WithContext(ctx).Where("id =?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUserVerify(ctx context.Context, verified *model.UserVerify) error {
	if err := r.db.WithContext(ctx).Save(verified).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteUserVerify(ctx context.Context, verified *model.UserVerify) error {
	if err := r.db.WithContext(ctx).Delete(verified).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, user *model.User) (string, error) {
	if err := r.db.WithContext(ctx).Delete(user).Error; err != nil {
		return "", err
	}

	return user.Id.String(), nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
