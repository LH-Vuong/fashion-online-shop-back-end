package response

type PagingResponse[T any] struct {
	Data      []T    `json:"data"`
	TotalPage int    `json:"total_page"`
	Page      int    `json:"page"`
	Status    string `json:"status"`
}
