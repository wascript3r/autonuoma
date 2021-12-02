package user

import (
	"encoding/json"
	"strings"
	"time"
)

type UserInfo struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Create

type BirthDate time.Time

const BirthDateFormat = "2006-01-02"

type CreateReq struct {
	Email     string    `json:"email" validate:"required,u_email"`
	Password  string    `json:"password" validate:"required,u_password"`
	FirstName string    `json:"firstName" validate:"required,u_firstName"`
	LastName  string    `json:"lastName" validate:"required,u_lastName"`
	BirthDate BirthDate `json:"birthDate" validate:"required"`
	PIN       string    `json:"pin" validate:"required,u_pin"`
}

func (bd *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	t, err := time.Parse(BirthDateFormat, s)
	if err != nil {
		return err
	}

	*bd = BirthDate(t)
	return nil
}

func (bd BirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(bd))
}

// Authenticate

type AuthenticateReq struct {
	Email    string `json:"email" validate:"required,u_email"`
	Password string `json:"password" validate:"required,u_password"`
}

// TempToken

type TempToken struct {
	Token string `json:"token"`
}
