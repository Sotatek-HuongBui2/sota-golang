package models

import (
	"errors"
	"time"

	"vtcanteen/constants"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Orders struct {
	Id                string `gorm:"PRIMARY_KEY;" json:"id"`
	CustomerId        string `gorm:"NOT_NULL" json:"customer_id"`
	AcceptedId        string `gorm:"NOT_NULL" json:"accepted_id"`
	OrderNumber       string `gorm:"NOT_NULL;unique" json:"order_number"`
	OrderStatus       string `gorm:"NOT_NULL" json:"order_status"`
	PaymentStatus     string `gorm:"NOT_NULL" json:"payment_status"`
	FulfillmentStatus string `gorm:"NOT_NULL" json:"fulfillment_status"`
	OrderBy           string `gorm:"NOT_NULL" json:"order_by"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}

func (o *Orders) BeforeCreate(tx *gorm.DB) (err error) {
	o.Id = uuid.New().String()
	if o.Id == constants.EMPTY_STRING {
		err = errors.New("Cannot save invalid data")
	}
	return
}
