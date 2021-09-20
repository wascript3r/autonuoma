package session

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Create(ctx context.Context, userID int) (*domain.Session, error)
	IsExpired(ss *domain.Session) bool
	Validate(ctx context.Context, id string) (*domain.Session, error)
	Delete(ctx context.Context, id string) error
	GenTempToken(ss *domain.Session) (string, error)
	ValidateTempToken(ctx context.Context, token string) (*domain.Session, error)
	StoreCtx(ctx context.Context, ss *domain.Session) context.Context
	LoadCtx(ctx context.Context) (*domain.Session, error)
}
