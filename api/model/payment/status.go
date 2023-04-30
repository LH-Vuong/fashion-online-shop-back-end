package payment

type Status string

const (
	StatusInit    Status = "INIT"
	StatusPending Status = "PENDING"
	StatusCancel  Status = "CANCEL"
	StatusError   Status = "ERROR"

	StatusApproved = "APPROVED"
)
