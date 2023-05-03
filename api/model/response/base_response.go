package response

type BaseResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
