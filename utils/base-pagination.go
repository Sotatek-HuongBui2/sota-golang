package utils

type IPagination[D any] struct {
	Data       D     `json:"data"`
	TotalCount int   `json:"count"`
	Meta       any   `json:"meta"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
}

func PaginateResult[D any](data D, totalCount int, page int64, limit int64, meta ...any) *IPagination[D] {
	resp := &IPagination[D]{Data: data, TotalCount: totalCount, Meta: meta, Page: page, Limit: limit}
	return resp
}
