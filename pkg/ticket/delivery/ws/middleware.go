package ws

import (
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/pool"
)

type Middleware interface {
	GetRoomName(ticketID int) pool.RoomName
	CreateOrRejoinRoom(s *gows.Socket, ticketID int) error
	LeaveCurrentRoom(s *gows.Socket) error
	DeleteRoom(ticketID int) error
}
