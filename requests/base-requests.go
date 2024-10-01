package requests

type GetList struct {
	SearchFields string `json:"search_fields"`
	Search       string `json:"search"`
	Filter       string `json:"filter"`
	Sort         string `json:"sort"`
	SortDir      string `json:"sort_dir"`
	Limit        string `json:"limit"`
	Page         string `json:"page"`
}
