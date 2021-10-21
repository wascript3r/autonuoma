package middleware

import (
	"github.com/wascript3r/autonuoma/pkg/ticket"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/pool"
)

const DefaultRoomKey = "ticket"

type WSMiddleware struct {
	roomKey     string
	ticketUcase ticket.Usecase
}

func NewWSMiddleware(roomKey string, tu ticket.Usecase) *WSMiddleware {
	return &WSMiddleware{roomKey, tu}
}

func (w *WSMiddleware) ExtractRoom(s *gows.Socket) (pool.Room, bool) {
	data, ok := s.GetData(w.roomKey)
	if !ok {
		return "", false
	}
	r, ok := data.(pool.Room)
	return r, ok
}

func (w *WSMiddleware) SetRoom(s *gows.Socket, r pool.Room) {
	s.SetData(w.roomKey, r)
}

func (w *WSMiddleware) DeleteRoom(s *gows.Socket) {
	s.DeleteData(w.roomKey)
}
