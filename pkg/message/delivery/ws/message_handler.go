package ws

import (
	"context"
	"encoding/json"

	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/session"
	ticketHandler "github.com/wascript3r/autonuoma/pkg/ticket/delivery/ws"
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
	ticketMid    ticketHandler.Middleware
	socketPool   *pool.Pool
}

func NewWSHandler(r *router.Router, client *middleware.Stack, agent *middleware.Stack, mu message.Usecase, meb message.EventBus, su session.Usecase, tm ticketHandler.Middleware, socketPool *pool.Pool) {
	handler := &WSHandler{
		messageUcase: mu,
		sessionUcase: su,
		ticketMid:    tm,
		socketPool:   socketPool,
	}

	meb.Subscribe(message.NewMessageEvent, handler.NewMessageNotification("message/notification"))
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

	err = w.messageUcase.ClientSend(ctx, ss.UserID, req)
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

	err = w.messageUcase.AgentSend(ctx, ss.UserID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) NewMessageNotification(method string) func(context.Context, *message.TicketMessage) {
	return func(ctx context.Context, tm *message.TicketMessage) {
		rName := w.ticketMid.GetRoomName(tm.TicketID)

		w.socketPool.EmitRoom(rName, &router.Response{
			Err:    nil,
			Method: &method,
			Params: tm,
		})
	}
}
