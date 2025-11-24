package pagination

type PaginationRes struct {
	Data      any
	Total     int32
	Page      int32
	Limit     int32
	TotalPage int32
}
