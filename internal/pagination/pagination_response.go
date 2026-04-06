package pagination

type PaginationRes struct {
	Data      any   `json:"data"`
	Total     int32 `json:"total"`
	Page      int32 `json:"page"`
	Limit     int32 `json:"limit"`
	TotalPage int32 `json:"total_page"`
}
