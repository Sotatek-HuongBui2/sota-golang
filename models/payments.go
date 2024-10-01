package models

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Payments struct {
	Id            string `gorm:"PRIMARY_KEY;" json:"id"`
	TransactionId string `gorm:"NOT_NULL" json:"transaction_id"`
	PaidAmount    string `gorm:"NOT_NULL" json:"paid_amount"`
	Type          string `gorm:"NOT_NULL" json:"type"`
	ProcessedAt   string `gorm:"NOT_NULL" json:"processed_at"`
	Metadata      string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (payment *Payments) BeforeCreate(tx *gorm.DB) (err error) {
	payment.Id = uuid.New().String()
	if payment == nil {
		err = errors.New("can't save invalid data")
	}
	return
}
