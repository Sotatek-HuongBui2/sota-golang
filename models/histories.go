package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Histories struct {
	Id          string `gorm:"PRIMARY_KEY;" json:"id"`
	EntityCode  string `gorm:"NOT_NULL" json:"entity_code"`
	EntityId    string `gorm:"NOT_NULL" json:"entity_id"`
	ActionName  string `gorm:"NOT_NULL" json:"action_name"`
	ProcessedAt string `gorm:"NOT_NULL" json:"processed_at"`
	Metadata    string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
