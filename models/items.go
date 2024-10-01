package models

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Items struct {
	Id           string `gorm:"PRIMARY_KEY;" json:"id"`
	OrderId      string `gorm:"NOT_NULL" json:"orders_id"`
	ProductId    string `gorm:"NOT_NULL" json:"product_id"`
	OrderedQty   string `gorm:"NOT_NULL" json:"ordered_qty"`
	FulfilledQty string `gorm:"NOT_NULL" json:"fulfilled_qty"`
	RefundQty    string `gorm:"NOT_NULL" json:"refund_qty"`
	SpecialPrice string `gorm:"NOT_NULL" json:"special_price"`
	Price        string `gorm:"NOT_NULL" json:"price"`
	SKU          string `gorm:"NOT_NULL" json:"sku"`
	Barcode      string `gorm:"NOT_NULL;unique" json:"barcode"`
	Metadata     string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (item *Items) BeforeCreate(tx *gorm.DB) (err error) {
	item.Id = uuid.New().String()
	if item == nil {
		err = errors.New("can't save invalid data")
	}
	return
}
