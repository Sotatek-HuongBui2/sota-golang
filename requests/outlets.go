package requests

type CreateOrUpdateOutlet struct {
	Id          string `json:"id"`
	WarehouseId string `json:"warehouse_id"`
	OutletName  string `json:"outlet_name"`
	IsActive    bool   `json:"is_active"`
	Metadata    string `json:"meta_data"`
}

type GetOutlets struct {
	GetList
}
