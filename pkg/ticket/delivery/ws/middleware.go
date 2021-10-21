package ws

import (
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/pool"
)

type Middleware interface {
	ExtractRoom(s *gows.Socket) (pool.Room, bool)
	SetRoom(s *gows.Socket, r pool.Room)
	DeleteRoom(s *gows.Socket)
}
