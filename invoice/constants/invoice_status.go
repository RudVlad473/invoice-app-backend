package invoice

type Status string

const (
	StatusPending Status = "PENDING"
	StatusPaid    Status = "PAID"
	StatusDraft   Status = "DRAFT"
)

var Statuses = []Status{StatusPending, StatusPaid, StatusDraft}
