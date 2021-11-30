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

func (w *WSMiddleware) ExtractRoom(s *gows.Socket) (pool.RoomName, bool) {
	data, ok := s.GetData(w.roomKey)
	if !ok {
		return "", false
	}
	r, ok := data.(pool.RoomName)
	return r, ok
}

func (w *WSMiddleware) SetRoom(s *gows.Socket, name pool.RoomName) {
	s.SetData(w.roomKey, name)
}

func (w *WSMiddleware) DeleteRoom(s *gows.Socket) {
	s.DeleteData(w.roomKey)
}
