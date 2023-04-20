package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key"`
	Username    string    `gorm:"type:varchar(100);not null"`
	Password    string    `gorm:"type:varchar(200);not null"`
	Fullname    string    `gorm:"type:varchar(100)"`
	Email       string    `gorm:"uniqueIndex;type:varchar(200);not null"`
	Address     string    `gorm:"type:varchar(255)"`
	PhoneNumber string    `gorm:"type:varchar(255)"`
	Avatar      string    `gorm:"type:varchar(100)"`
	Verified    bool      `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time
	CreatedBy   string `gorm:"type:varchar(100)"`
	UpdatedAt   time.Time
	UpdatedBy   string `gorm:"type:varchar(100)"`
}

type UserVerify struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key"`
	Username    string    `gorm:"type:varchar(100);not null"`
	UniqueToken string    `gorm:"type:varchar(50);not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time ``
	CreatedBy   string    `gorm:"type:varchar(100)"`
	UpdatedAt   time.Time ``
	UpdatedBy   string    `gorm:"type:varchar(100)"`
}

func (User) TableName() string {
	return "SYS_User"
}

func (UserVerify) TableName() string {
	return "SYS_UserVerify"
}
