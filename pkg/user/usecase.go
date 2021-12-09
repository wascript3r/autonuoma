package user

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Usecase interface {
	Create(ctx context.Context, req *CreateReq) error
	Authenticate(ctx context.Context, req *AuthenticateReq) (*domain.Session, *AuthenticateRes, error)
	GetTempToken(ss *domain.Session) (*TempToken, error)
	AuthenticateToken(ctx context.Context, req *TempToken) (*domain.Session, error)
	Logout(ctx context.Context, ss *domain.Session) error
	GetInfo(userID int, role domain.Role) *AuthenticateRes
}
