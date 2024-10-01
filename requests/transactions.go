package requests

import "vtcanteen/models"

type CreateOrUpdateTransaction struct {
	Id          string             `gorm:"PRIMARY_KEY;" json:"id"`
	OrderId     string             `gorm:"NOT_NULL" json:"order_id"`
	PaidAmount  string             `gorm:"NOT_NULL" json:"paid_amount"`
	Type        string             `gorm:"NOT_NULL" json:"type"`
	ProcessedAt string             `gorm:"NOT_NULL" json:"processed_at"`
	Metadata    string             `json:"meta_data"`
	Payments    []*models.Payments `json:"payments"`
}
