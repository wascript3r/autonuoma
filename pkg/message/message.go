package message

// Send

type ClientSendReq struct {
	Message string `json:"message" validate:"required,m_message"`
}
