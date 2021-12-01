package domain

import "time"

type Role int8

const (
	UserRole Role = iota + 1
	AgentRole
	AdminRole
)

type User struct {
	ID        int
	Email     string
	Password  string
	FirstName string
	LastName  string
	BirthDate time.Time
	Balance   int64
	PIN       string
	RoleID    Role
}

type UserMeta struct {
	ID        int
	FirstName string
	LastName  string
}
