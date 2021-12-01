package room

import "github.com/wascript3r/autonuoma/pkg/domain"

type Usecase interface {
	Register(r domain.Room, c Config) error
	GetName(r domain.Room) (string, error)
}
