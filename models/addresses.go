package models

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Addresses struct {
	Id          string `gorm:"PRIMARY_KEY;" json:"id"`
	OrderId     string `gorm:"NOT_NULL" json:"order_id"`
	Country     string `gorm:"NOT_NULL" json:"country"`
	CountryCode string `gorm:"NOT_NULL" json:"country_code"`
	Region      string `gorm:"NOT_NULL" json:"region"`
	RegionCode  string `gorm:"NOT_NULL" json:"region_code"`
	Address     string `gorm:"NOT_NULL" json:"address"`
	Type        string `gorm:"NOT_NULL" json:"type"`
	Email       string `gorm:"NOT_NULL" json:"email"`
	Phone       string `gorm:"NOT_NULL" json:"phone"`
	Name        string `gorm:"NOT_NULL" json:"name"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (address *Addresses) BeforeCreate(tx *gorm.DB) (err error) {
	address.Id = uuid.New().String()
	if address == nil {
		err = errors.New("can't save invalid data")
	}
	return
}
