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
	"github.com/wascript3r/gows/pool"
	"github.com/wascript3r/gows/router"
)

var (
	AuthenticatedRoomConfig = pool.NewRoomConfig("auth", false)
)

type WSHandler struct {
	userUcase    user.Usecase
	sessionUcase session.Usecase
	sessionMid   sessionHandler.Middleware
	socketPool   *pool.Pool
}

func NewWSHandler(r *router.Router, notAuth *middleware.Stack, uu user.Usecase, su session.Usecase, sm sessionHandler.Middleware, socketPool *pool.Pool) {
	handler := &WSHandler{
		userUcase:    uu,
		sessionUcase: su,
		sessionMid:   sm,
		socketPool:   socketPool,
	}

	socketPool.CreateRoom(AuthenticatedRoomConfig)

	r.HandleMethod("user/authenticate", notAuth.Wrap(handler.Authenticate))
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
	w.socketPool.JoinRoom(s, AuthenticatedRoomConfig.Name)

	w.socketPool.EmitRoom(AuthenticatedRoomConfig.Name, &router.Response{
		Err:    nil,
		Method: &r.Method,
		Params: nil,
	})
}
