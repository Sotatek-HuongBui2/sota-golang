package requests

type CreateOrUpdateWarehouse struct {
	Id            string `json:"id"`
	WarehouseName string `json:"warehouse_name"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
	Region        string `json:"region"`
	RegionCode    string `json:"region_code"`
	IsActive      bool   `json:"is_active"`
}

type GetWarehouses struct {
	GetList
}
