package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Customers struct {
	Id          string `gorm:"PRIMARY_KEY;" json:"id"`
	UserId      string `gorm:"NOT_NULL" json:"user_id"`
	UserName    string `gorm:"NOT_NULL;unique" json:"user_name"`
	Password    string `gorm:"NOT_NULL" json:"password"`
	Email       string `gorm:"NOT_NULL;unique" json:"email"`
	FirstName   string `gorm:"NOT_NULL" json:"first_name"`
	MiddleName  string `gorm:"NOT_NULL" json:"middle_name"`
	LastName    string `gorm:"NOT_NULL" json:"last_name"`
	Country     string `gorm:"NOT_NULL" json:"country"`
	CountryCode string `gorm:"NOT_NULL" json:"country_code"`
	Region      string `gorm:"NOT_NULL" json:"region"`
	RegionCode  string `gorm:"NOT_NULL" json:"region_code"`
	Address     string `gorm:"NOT_NULL" json:"address"`
	IsActive    bool   `gorm:"NOT_NULL" json:"is_active"`
	Metadata    string `json:"meta_data"`

	CreatedAt time.Time  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime" json:"deleted_at"`
}
