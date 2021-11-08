package models

import (
	"github.com/pilagod/gorm-cursor-paginator/v2/cursor"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"column:username;type:varchar(255);not null"`
	Password  string `gorm:"column:password;type:varchar(255);not null"`
	Email     string `gorm:"unique;column:email;type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "user"
}

type UserPaginator struct {
	Users []*User
	cursor.Cursor
}
