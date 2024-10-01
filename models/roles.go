package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Roles struct {
	Id          string `gorm:"PRIMARY_KEY;" json:"id"`
	RoleName    string `gorm:"NOT_NULL;unique" json:"role_name"`
	Permissions string `gorm:"NOT_NULL" json:"permissions"`
	Level       int    `gorm:"NOT_NULL" json:"level"`
	Metadata    string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
