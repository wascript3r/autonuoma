package ws

import (
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/pool"
)

type Middleware interface {
	ExtractRoom(s *gows.Socket) (pool.RoomName, bool)
	SetRoom(s *gows.Socket, name pool.RoomName)
	DeleteRoom(s *gows.Socket)
}
