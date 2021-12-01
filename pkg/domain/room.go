package domain

import "errors"

type Room int

const (
	AuthenticatedRoom Room = iota
	SupportRoom
)

var ErrInvalidRoomName = errors.New("invalid room name")

func IsValidRoom(r Room) bool {
	switch r {
	case AuthenticatedRoom, SupportRoom:
		return true
	}
	return false
}
