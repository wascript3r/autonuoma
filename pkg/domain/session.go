package domain

import "time"

type Role int8

const (
	UserRole Role = iota + 1
	SupportRole
	AdminRole
)

func HasRole(ss *Session, r Role) bool {
	return ss.Role == r
}

type Session struct {
	ID         string
	UserID     int
	Expiration time.Time
	Role       Role
}
