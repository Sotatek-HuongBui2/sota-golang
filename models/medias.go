package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Medias struct {
	Id          string `gorm:"PRIMARY_KEY;" json:"id"`
	EntityCode  string `gorm:"NOT_NULL" json:"entity_code"`
	EntityId    string `gorm:"NOT_NULL" json:"entity_id"`
	MediaURL    string `gorm:"NOT_NULL" json:"media_url"`
	Description string `gorm:"NOT_NULL" json:"description"`
	Name        string `gorm:"NOT_NULL" json:"name"`
	Metadata    string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
