package domain

type Ticket struct {
	ID      int
	UserID  int
	AgentID *int
	Ended   bool
}
