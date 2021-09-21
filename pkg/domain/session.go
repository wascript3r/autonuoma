package domain

import "time"

func HasRole(ss *Session, r Role) bool {
	return ss.RoleID == r
}

type Session struct {
	ID         string
	UserID     int
	Expiration time.Time
	RoleID     Role
}
