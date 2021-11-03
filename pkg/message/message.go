package message

// Send

type ClientSendReq struct {
	Message string `json:"message" validate:"required,u_message"`
}
