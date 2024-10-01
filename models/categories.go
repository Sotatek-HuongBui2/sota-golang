package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Categories struct {
	Id           string `gorm:"PRIMARY_KEY;" json:"id"`
	ParentId     string `gorm:"NOT_NULL" json:"parent_id"`
	CategoryName string `gorm:"NOT_NULL;unique" json:"category_name"`
	IsActive     bool   `gorm:"NOT_NULL" json:"is_active"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
