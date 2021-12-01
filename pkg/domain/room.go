package domain

import "errors"

type Room int

const (
	AuthenticatedRoom Room = iota
	AgentRoom
)

var ErrInvalidRoomName = errors.New("invalid room name")

func IsValidRoom(r Room) bool {
	switch r {
	case AuthenticatedRoom, AgentRoom:
		return true
	}
	return false
}
