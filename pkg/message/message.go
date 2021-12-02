package message

import (
	"time"

	"github.com/wascript3r/autonuoma/pkg/user"
)

// Send

type SendReq struct {
	TicketID int    `json:"ticketID" validate:"required"`
	Message  string `json:"message" validate:"required,m_message"`
}

type MessageInfo struct {
	User    *user.UserInfo `json:"user"`
	Content string         `json:"content"`
	Time    time.Time      `json:"time"`
}

type TicketMessage struct {
	TicketID int `json:"ticketID"`
	*MessageInfo
}
