package ws

import (
	"context"
	"encoding/json"

	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/autonuoma/pkg/user"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	"github.com/wascript3r/gows"
	"github.com/wascript3r/gows/middleware"
	"github.com/wascript3r/gows/pool"
	"github.com/wascript3r/gows/router"
)

type WSHandler struct {
	messageUcase message.Usecase
	sessionUcase session.Usecase
	socketPool   *pool.Pool
}

func NewWSHandler(r *router.Router, client *middleware.Stack, agent *middleware.Stack, mu message.Usecase, su session.Usecase, socketPool *pool.Pool) {
	handler := &WSHandler{
		messageUcase: mu,
		sessionUcase: su,
		socketPool:   socketPool,
	}

	r.HandleMethod("ticket/client/message/new", client.Wrap(handler.ClientNewMessage))
	r.HandleMethod("ticket/agent/message/new", agent.Wrap(handler.AgentNewMessage))
}

func serveError(s *gows.Socket, r *router.Request, err error) {
	code := errcode.UnwrapErr(err, user.UnknownError)
	router.WriteErr(s, code, &r.Method)
}

func (w *WSHandler) ClientNewMessage(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &message.ClientSendReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	_, err = w.messageUcase.ClientSend(ctx, ss.UserID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) AgentNewMessage(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &message.AgentSendReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	_, err = w.messageUcase.AgentSend(ctx, ss.UserID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}
