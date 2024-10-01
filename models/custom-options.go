package models

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type CustomOptions struct {
	Id          string         `gorm:"PRIMARY_KEY;" json:"id"`
	ProductId   string         `gorm:"NOT_NULL" json:"product_id"`
	OptionName  string         `gorm:"NOT_NULL;unique" json:"option_name"`
	Position    string         `gorm:"NOT_NULL" json:"position"`
	Price       string         `gorm:"NOT_NULL" json:"price"`
	OptionItems []*OptionItems `gorm:"foreignKey:option_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"option_items"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (co *CustomOptions) BeforeCreate(tx *gorm.DB) (err error) {
	co.Id = uuid.New().String()
	if co == nil {
		err = errors.New("can't save invalid data")
	}
	return
}
