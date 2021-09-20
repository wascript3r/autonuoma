package session

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Repository interface {
	Insert(ctx context.Context, ss *domain.Session) error
	Get(ctx context.Context, id string) (*domain.Session, error)
	Delete(ctx context.Context, id string) error
}
