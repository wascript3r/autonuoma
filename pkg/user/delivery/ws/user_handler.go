package ws

import (
	"context"
	"encoding/json"

	"github.com/wascript3r/autonuoma/pkg/session"
	sessionHandler "github.com/wascript3r/autonuoma/pkg/session/delivery/ws"
	"github.com/wascript3r/autonuoma/pkg/user"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/middleware"
	"github.com/wascript3r/gows/router"
)

type WSHandler struct {
	userUcase    user.Usecase
	sessionUcase session.Usecase
	sessionMid   sessionHandler.Middleware
}

func NewWSHandler(r *router.Router, admin *middleware.Stack, notAuth *middleware.Stack, uu user.Usecase, su session.Usecase, sm sessionHandler.Middleware) {
	handler := &WSHandler{
		userUcase:    uu,
		sessionUcase: su,
		sessionMid:   sm,
	}

	r.HandleMethod("authenticate", notAuth.Wrap(handler.Authenticate))
	r.HandleMethod("test", admin.Wrap(handler.Test))
}

func serveError(s *gows.Socket, r *router.Request, err error) {
	code := errcode.UnwrapErr(err, user.UnknownError)
	router.WriteErr(s, code, &r.Method)
}

func (w *WSHandler) Authenticate(ctx context.Context, s *gows.Socket, r *router.Request) {
	req := &user.TempToken{}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	ss, err := w.userUcase.AuthenticateToken(ctx, req)
	if err != nil {
		serveError(s, r, err)
		return
	}
	w.sessionMid.SetSession(s, ss)

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) Test(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, router.Params{
		"sessID": ss.ID,
		"uID":    ss.UserID,
	})
}
