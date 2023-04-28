package model

import (
	"time"
)

type User struct {
	Id        string    `bson:"_id,omitempty" json:"id"`
	Password  string    `bson:"password" json:"-"`
	Fullname  string    `bson:"fullname" json:"fullname"`
	Email     string    `bson:"email" json:"email"`
	Photo     string    `bson:"photo" json:"photo"`
	Verified  bool      `bson:"verified" json:"verified"`
	Status    string    `bson:"status" json:"status"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type UserVerify struct {
	Id          string    `bson:"_id,omitempty" json:"id"`
	UserId      string    `bson:"user_id" json:"user_id"`
	UniqueToken string    `bson:"unique_token" json:"unique_token"`
	Status      string    `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type UserWishlist struct {
	Id        string    `bson:"_id,omitempty" json:"id"`
	UserId    string    `bson:"user_id" json:"user_id"`
	ProductId string    `bson:"product_id" json:"product_id"`
	Status    string    `bson:"status" json:"status"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type UserAddress struct {
	Id         string    `bson:"_id,omitempty" json:"id"`
	UserId     string    `bson:"user_id" json:"user_id"`
	Name       string    `bson:"name" json:"name"`
	ProvinceId int       `bson:"province_id" json:"province_id"`
	DistrictId int       `bson:"district_id" json:"district_id"`
	IsDefault  bool      `bson:"is_default" json:"is_default"`
	WardCode   string    `bson:"ward_code" json:"ward_code"`
	Address    string    `bson:"status" json:"status"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type CreateUserAddressModel struct {
	ProvinceId int    `json:"province_id" binding:"required"`
	DistrictId int    `json:"district_id" binding:"required"`
	WardCode   string `json:"ward_code" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type UpdateUserAddressModel struct {
	Id         string `json:"id" binding:"required"`
	ProvinceId int    `json:"province_id" binding:"required"`
	DistrictId int    `json:"district_id" binding:"required"`
	WardCode   string `json:"ward_code" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type GetUserAddressListModel struct {
	Page     int64 `json:"page" binding:"required"`
	PageSize int64 `json:"page_size" binding:"required"`
}

type SignUpModel struct {
	Fullname        string `json:"fullname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,min=8"`
}

type SignInModel struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type AddWishListModel struct {
	ProductId string `json:"product_id" binding:"required"`
}

type DeleteWishListModel struct {
	DeleteIds []string `json:"delete_ids" binding:"required"`
}

type GetUserWishlistModel struct {
	Page     int64 `json:"page" binding:"required"`
	PageSize int64 `json:"page_size" binding:"required"`
}

type ChangePasswordModel struct {
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,min=8"`
}
