package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key"`
	Password    string    `gorm:"type:varchar(200);not null;" json:"-"`
	Fullname    string    `gorm:"type:varchar(100)"`
	Email       string    `gorm:"uniqueIndex;type:varchar(200);not null"`
	PhoneNumber string    `gorm:"type:varchar(255)"`
	Avatar      string    `gorm:"type:varchar(100)"`
	Verified    bool      `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserVerify struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key"`
	UserId      uuid.UUID `gorm:"type:uuid;not null"`
	UniqueToken string    `gorm:"type:varchar(50);not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time ``
	UpdatedAt   time.Time ``
}

func (User) TableName() string {
	return "SYS_User"
}

func (UserVerify) TableName() string {
	return "SYS_UserVerify"
}

type SignUpModel struct {
	Fullname        string `json:"fullname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

type SignInModel struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
