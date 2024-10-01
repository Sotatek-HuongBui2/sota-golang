package requests

import "vtcanteen/models"

type CreateOrUpdateProduct struct {
	Id              string                 `json:"id"`
	ParentId        string                 `json:"parent_id"`
	ProductName     string                 `json:"product_name"`
	SKU             string                 `json:"sku"`
	Barcode         string                 `json:"barcode"`
	Type            string                 `json:"type"`
	Price           string                 `json:"price"`
	SpecialPrice    string                 `json:"special_price"`
	ManageStock     string                 `json:"manage_stock"`
	SafetyThreshold string                 `json:"safety_threshold"`
	Taxable         bool                   `json:"taxale"`
	ImageURL        string                 `json:"image_url"`
	IsVariant       bool                   `json:"is_variant"`
	IsActive        bool                   `json:"is_active"`
	Options         []*models.Options       `json:"options"`
	CustomOptions   *models.CustomOptions `json:"custom_options"`
}

type GetProducts struct {
	GetList
}
