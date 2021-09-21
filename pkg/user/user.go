package user

// Create

type CreateReq struct {
	Email           string `json:"email" validate:"required,u_email"`
	Password        string `json:"password" validate:"required,u_password"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,u_password,eqfield=Password"`
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
