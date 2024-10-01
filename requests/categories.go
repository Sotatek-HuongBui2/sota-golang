package requests

type CreateOrUpdateCategory struct {
	Id           string `json:"id"`
	ParentId     string `json:"parent_id"`
	CategoryName string `json:"category_name"`
	IsActive     bool   `json:"is_active"`
}

type GetCategories struct {
	GetList
}
