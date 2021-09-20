package middleware

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/router"
)

const DefaultSessionKey = "sess"

type WSMiddleware struct {
	sessionKey   string
	sessionUcase session.Usecase
}

func NewWSMiddleware(sessionKey string, su session.Usecase) *WSMiddleware {
	return &WSMiddleware{sessionKey, su}
}

func (w *WSMiddleware) ExtractSession(s *gows.Socket) (*domain.Session, bool) {
	data, ok := s.GetData(w.sessionKey)
	if !ok {
		return nil, false
	}
	ss, ok := data.(*domain.Session)
	return ss, ok
}

func (w *WSMiddleware) SetSession(s *gows.Socket, ss *domain.Session) {
	s.SetData(w.sessionKey, ss)
}

func (w *WSMiddleware) DeleteSession(s *gows.Socket) {
	s.DeleteData(w.sessionKey)
}

func (w *WSMiddleware) Authenticated(next router.Handler) router.Handler {
	return func(ctx context.Context, s *gows.Socket, r *router.Request) {
		ss, ok := w.ExtractSession(s)
		if !ok {
			router.WriteErr(s, session.NotAuthenticatedError, &r.Method)
			return
		}

		if w.sessionUcase.IsExpired(ss) {
			w.DeleteSession(s)
			router.WriteErr(s, session.SessionExpiredError, &r.Method)
			return
		}
		ctx = w.sessionUcase.StoreCtx(ctx, ss)

		next(ctx, s, r)
	}
}

func (w *WSMiddleware) NotAuthenticated(next router.Handler) router.Handler {
	return func(ctx context.Context, s *gows.Socket, r *router.Request) {
		ss, ok := w.ExtractSession(s)
		if !ok {
			next(ctx, s, r)
			return
		}

		if w.sessionUcase.IsExpired(ss) {
			w.DeleteSession(s)
			next(ctx, s, r)
			return
		}

		router.WriteErr(s, session.AlreadyAuthenticatedError, &r.Method)
	}
}

func (w *WSMiddleware) HasRole(role domain.Role) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		return w.Authenticated(
			func(ctx context.Context, s *gows.Socket, r *router.Request) {
				ss, err := w.sessionUcase.LoadCtx(ctx)
				if err != nil {
					router.WriteInternalError(s, &r.Method)
					return
				}

				if !domain.HasRole(ss, role) {
					router.WriteErr(s, session.InsufficientPermissionsError, &r.Method)
					return
				}

				next(ctx, s, r)
			},
		)
	}
}
