package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type WarehouseItems struct {
	Id              string `gorm:"PRIMARY_KEY;" json:"id"`
	WarehouseId     string `gorm:"PRIMARY_KEY" json:"warehouse_id"`
	ProductId       string `gorm:"PRIMARY_KEY" json:"product_id"`
	AvailaibleQty   string `gorm:"NOT_NULL" json:"available_qty"`
	SafetyThreshold string `gorm:"NOT_NULL" json:"safety_threshold"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
