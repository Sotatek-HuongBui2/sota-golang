package models

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type OptionItems struct {
	Id        string `gorm:"PRIMARY_KEY;" json:"id"`
	OptionId  string `gorm:"PRIMARY_KEY" json:"option_id"`
	ProductId string `gorm:"PRIMARY_KEY" json:"product_id"`
	Metadata  string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (oi *OptionItems) BeforeCreate(tx *gorm.DB) (err error) {
	oi.Id = uuid.New().String()
	if oi == nil {
		err = errors.New("can't save invalid data")
	}
	return
}
