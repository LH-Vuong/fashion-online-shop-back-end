package order

type Status string

const (
	PENDING Status = "PENDING"
	CANCEL  Status = "CANCEL"
	SUCCESS Status = "SUCCESS"
)
