package models

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Products struct {
	Id              string `gorm:"PRIMARY_KEY;" json:"id"`
	ParentId        string `gorm:"NOT_NULL" json:"parent_id"`
	ProductName     string `gorm:"NOT_NULL" json:"product_name"`
	SKU             string `gorm:"NOT_NULL" json:"sku"`
	Barcode         string `gorm:"NOT_NULL" json:"barcode"`
	Type            string `gorm:"NOT_NULL" json:"type"`
	Price           string `gorm:"NOT_NULL" json:"price"`
	SpecialPrice    string `gorm:"NOT_NULL" json:"special_price"`
	ManageStock     string `gorm:"NOT_NULL" json:"manage_stock"`
	SafetyThreshold string `gorm:"NOT_NULL" json:"safety_threshold"`
	Taxable         bool   `gorm:"NOT_NULL" json:"taxable"`
	ImageURL        string `gorm:"NOT_NULL" json:"image_url"`
	IsVariant       bool   `gorm:"NOT_NULL" json:"is_variant"`
	IsActive        bool   `gorm:"NOT_NULL" json:"is_active"`
	Metadata        string `json:"meta_data"`

	Options       []Options       `gorm:"foreignKey:product_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"options"`
	CustomOptions []CustomOptions `gorm:"foreignKey:product_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"custom_options"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (p *Products) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id = uuid.New().String()
	if p == nil {
		err = errors.New("can't save invalid data")
	}
	return
}
