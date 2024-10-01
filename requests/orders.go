package requests

import "vtcanteen/models"

type CreateOrUpdateOrder struct {
	Id                string                       `json:"id"`
	CustomerId        string                       `json:"customer_id"`
	AcceptedId        string                       `json:"accepted_id"`
	OrderNumber       string                       `json:"order_number"`
	OrderStatus       string                       `json:"order_status"`
	PaymentStatus     string                       `json:"payment_status"`
	FulfillmentStatus string                       `json:"fulfillment_status"`
	OrderBy           string                       `json:"order_by"`
	Address           *models.Addresses            `json:"address"`
	Items             []*models.Items              `json:"items"`
	Transactions      []*CreateOrUpdateTransaction `json:"transactions"`
}

type GetOrders struct {
	GetList
}
