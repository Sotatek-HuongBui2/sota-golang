package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ProductCategories struct {
	CategoryId string `gorm:"PRIMARY_KEY;" json:"category_id"`
	ProductId  string `gorm:"PRIMARY_KEY;" json:"product_id"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
