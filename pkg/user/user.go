package user

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type UserInfo struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserProfile struct {
	ID            int       `json:"id"`
	Balance       float32   `json:"balance"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	Birthdate     BirthDate `json:"birthdate"`
	LicenseStatus string    `json:"license,omitempty"`
}

type UserSensitiveInfo struct {
	*UserInfo
	PIN string `json:"pin"`
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

type AuthenticateRes struct {
	UserID int         `json:"userID"`
	RoleID domain.Role `json:"roleID"`
}

// TempToken

type TempToken struct {
	Token string `json:"token"`
}

// UpdateUser

type UpdateReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateRes struct {
	Email string `json:"email"`
}

// GetTrips

const TripDateTimeFormat = "2006-01-02 15:04:05"

type TripsRes struct {
	ID    int     `json:"id"`
	Begin string  `json:"begin_time"`
	End   string  `json:"end_time"`
	From  string  `json:"from"`
	To    string  `json:"to"`
	Price float32 `json:"price"`
}

// Payment

type PaymentRes struct {
	Status  string  `json:"status"`
	Balance float32 `json:"balance"`
}
