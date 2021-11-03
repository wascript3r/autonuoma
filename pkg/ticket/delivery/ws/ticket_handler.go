package ws

import (
	"context"
	"encoding/json"

	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/autonuoma/pkg/ticket"
	"github.com/wascript3r/autonuoma/pkg/user"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/middleware"
	"github.com/wascript3r/gows/pool"
	"github.com/wascript3r/gows/router"
)

type WSHandler struct {
	ticketUcase  ticket.Usecase
	sessionUcase session.Usecase
	ticketMid    Middleware
	socketPool   *pool.Pool
}

func NewWSHandler(r *router.Router, client *middleware.Stack, tu ticket.Usecase, su session.Usecase, tm Middleware, socketPool *pool.Pool) {
	handler := &WSHandler{
		ticketUcase:  tu,
		sessionUcase: su,
		ticketMid:    tm,
		socketPool:   socketPool,
	}

	r.HandleMethod("ticket/new", client.Wrap(handler.NewTicket))
}

func serveError(s *gows.Socket, r *router.Request, err error) {
	code := errcode.UnwrapErr(err, user.UnknownError)
	router.WriteErr(s, code, &r.Method)
}

func (w *WSHandler) NewTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &ticket.CreateReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	tID, err := w.ticketUcase.Create(ctx, ss.UserID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, tID)
}
