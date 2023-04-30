package response

type PagingResponse[T any] struct {
	Data   []T    `json:"data"`
	Total  int64  `json:"total"`
	Length int    `json:"length"`
	Status string `json:"status"`
}
