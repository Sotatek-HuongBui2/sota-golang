package models

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Options struct {
	Id          string `gorm:"PRIMARY_KEY;" json:"id"`
	ProductId   string `gorm:"NOT_NULL" json:"product_id"`
	OptionName  string `gorm:"NOT_NULL;unique" json:"option_name"`
	OptionValue string `gorm:"NOT_NULL" json:"option_value"`
	Position    string `gorm:"NOT_NULL" json:"position"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (o *Options) BeforeCreate(tx *gorm.DB) (err error) {
	o.Id = uuid.New().String()
	if o == nil {
		err = errors.New("can't save invalid data")
	}
	return
}
