package ws

import (
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/router"
)

type Middleware interface {
	Authenticated(next router.Handler) router.Handler
	NotAuthenticated(next router.Handler) router.Handler
	HasRole(role domain.Role) func(next router.Handler) router.Handler
	SetSession(s *gows.Socket, ss *domain.Session)
	DeleteSession(s *gows.Socket)
}
