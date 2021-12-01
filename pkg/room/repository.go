package room

import "github.com/wascript3r/autonuoma/pkg/domain"

type Repository interface {
	Set(r domain.Room, c Config) error
	GetName(r domain.Room) (string, error)
}
