package middleware

import (
	"errors"
	"fmt"

	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/pool"
)

const (
	DefaultRoomPrefix = "ticket"
	DefaultSocketKey  = "ticket"
)

var ErrTicketIDMismatch = errors.New("ticketID type mismatch")

type WSMiddleware struct {
	socketPool *pool.Pool
}

func NewWSMiddleware(socketPool *pool.Pool) *WSMiddleware {
	return &WSMiddleware{socketPool}
}

func (w *WSMiddleware) GetRoomName(ticketID int) pool.RoomName {
	return pool.RoomName(fmt.Sprintf("%s:%d", DefaultRoomPrefix, ticketID))
}

func (w *WSMiddleware) CreateOrRejoinRoom(s *gows.Socket, ticketID int) error {
	err := w.LeaveCurrentRoom(s)
	if err != nil {
		return err
	}

	name := w.GetRoomName(ticketID)
	if !w.socketPool.RoomExists(name) {
		err := w.socketPool.CreateRoom(pool.NewRoomConfig(name, true))
		if err != nil {
			return err
		}
	}

	err = w.socketPool.JoinRoom(s, name)
	if err != nil {
		return err
	}

	s.SetData(DefaultSocketKey, ticketID)
	return nil
}

func (w *WSMiddleware) LeaveCurrentRoom(s *gows.Socket) error {
	tID, ok := s.GetData(DefaultSocketKey)
	if !ok {
		return nil
	}

	tIDInt, ok := tID.(int)
	if !ok {
		return ErrTicketIDMismatch
	}

	return w.socketPool.LeaveRoom(s, w.GetRoomName(tIDInt))
}

func (w *WSMiddleware) DeleteRoom(ticketID int) error {
	return w.socketPool.DeleteRoom(w.GetRoomName(ticketID))
}
