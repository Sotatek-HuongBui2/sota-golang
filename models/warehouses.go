package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Warehouses struct {
	Id            string `gorm:"PRIMARY_KEY;" json:"id"`
	WarehouseName string `gorm:"NOT_NULL;unique" json:"warehouse_name"`
	Country       string `gorm:"NOT_NULL" json:"country"`
	CountryCode   string `gorm:"NOT_NULL" json:"country_code"`
	Region        string `gorm:"NOT_NULL" json:"region"`
	RegionCode    string `gorm:"NOT_NULL" json:"region_code"`
	IsActive      bool   `gorm:"NOT_NULL" json:"is_active"`
	Metadata      string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
