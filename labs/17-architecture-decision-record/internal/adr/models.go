package adr

type Status string

const (
	StatusProposed   Status = "Proposed"
	StatusAccepted   Status = "Accepted"
	StatusSuperseded Status = "Superseded"
	StatusDeprecated Status = "Deprecated"
	StatusRejected   Status = "Rejected"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusProposed, StatusAccepted, StatusSuperseded, StatusDeprecated, StatusRejected:
		return true
	default:
		return false
	}
}

type Record struct {
	ID           int
	Title        string
	Status       Status
	SupersededBy int // 0 if not superseded
	Supersedes   int // 0 if does not supersede
	Content      string
}
