package domain

import (
	"errors"
	"time"
)

type Role int8

const (
	ClientRole Role = iota + 1
	AgentRole
	AdminRole
)

var ErrInvalidUserRole = errors.New("invalid user role")

type User struct {
	ID        int
	Email     string
	Password  string
	FirstName string
	LastName  string
	BirthDate time.Time
	Balance   float32
	PIN       string
	RoleID    Role
}

type UserCredentials struct {
	ID       int
	RoleID   Role
	Password string
}

type UserMeta struct {
	ID        int
	FirstName string
	LastName  string
}

type UserSensitiveMeta struct {
	*UserMeta
	PIN string
}
