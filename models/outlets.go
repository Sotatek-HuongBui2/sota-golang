package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Outlets struct {
	Id          string `gorm:"PRIMARY_KEY;" json:"id"`
	WarehouseId string `gorm:"NOT_NULL" json:"warehouse_id"`
	OutletName  string `gorm:"NOT_NULL;unique" json:"outlet_name"`
	IsActive    bool   `gorm:"NOT_NULL" json:"is_active"`
	Metadata    string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
