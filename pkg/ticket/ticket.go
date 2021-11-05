package ticket

// Create

type CreateReq struct {
	Message string `json:"message" validate:"required,m_message"`
}
