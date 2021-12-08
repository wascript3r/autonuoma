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
	r.HandleMethod("client/ticket/message/new", client.Wrap(handler.NewMessage))
	r.HandleMethod("agent/ticket/message/new", agent.Wrap(handler.NewMessage))
}

func serveError(s *gows.Socket, r *router.Request, err error) {
	code := errcode.UnwrapErr(err, user.UnknownError)
	router.WriteErr(s, code, &r.Method)
}

func (w *WSHandler) NewMessage(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &message.SendReq{}

	err = json.Unmarshal(r.Data, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	err = w.messageUcase.Send(ctx, ss.UserID, ss.RoleID, req)
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
			Error:  nil,
			Method: &method,
			Data:   tm,
		})
	}
}
