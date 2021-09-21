package domain

type Role int8

const (
	UserRole Role = iota + 1
	SupportRole
	AdminRole
)

type User struct {
	ID       int
	Email    string
	Password string
	RoleID   Role
}
