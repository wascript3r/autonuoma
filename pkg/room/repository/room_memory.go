package repository

import (
	"errors"
	"sync"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/room"
)

var ErrRoomDoesNotExist = errors.New("room does not exist")

type MemoryRepo struct {
	mx      *sync.RWMutex
	configs map[domain.Room]room.Config
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		mx:      &sync.RWMutex{},
		configs: make(map[domain.Room]room.Config),
	}
}

func (m *MemoryRepo) Set(r domain.Room, c room.Config) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.configs[r] = c
	return nil
}

func (m *MemoryRepo) GetName(r domain.Room) (string, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	rc, ok := m.configs[r]
	if !ok {
		return "", ErrRoomDoesNotExist
	}

	return rc.NameString(), nil
}
