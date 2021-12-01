package ws

import (
	"context"
	"encoding/json"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/room"
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
	roomUcase    room.Usecase

	socketPool *pool.Pool
}

func NewWSHandler(r *router.Router, client *middleware.Stack, agent *middleware.Stack, tu ticket.Usecase, su session.Usecase, tm Middleware, teb ticket.EventBus, ru room.Usecase, socketPool *pool.Pool) {
	handler := &WSHandler{
		ticketUcase:  tu,
		sessionUcase: su,
		ticketMid:    tm,
		roomUcase:    ru,

		socketPool: socketPool,
	}

	teb.Subscribe(ticket.NewTicketEvent, handler.NewTicketNotification("ticket/notification"))
	r.HandleMethod("ticket/new", client.Wrap(handler.NewTicket))
	r.HandleMethod("ticket/accept", agent.Wrap(handler.AcceptTicket))
	r.HandleMethod("ticket/client/end", client.Wrap(handler.EndClientTicket))
	r.HandleMethod("ticket/agent/end", agent.Wrap(handler.EndAgentTicket))
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

func (w *WSHandler) AcceptTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &ticket.AcceptReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	err = w.ticketUcase.Accept(ctx, ss.UserID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) EndClientTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	err = w.ticketUcase.EndClient(ctx, ss.UserID)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) EndAgentTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &ticket.EndAgentReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	err = w.ticketUcase.EndAgent(ctx, ss.UserID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) NewTicketNotification(method string) func(ctx context.Context) {
	return func(ctx context.Context) {
		rName, err := w.roomUcase.GetName(domain.AgentRoom)
		if err != nil {
			return
		}

		w.socketPool.EmitRoom(pool.RoomName(rName), &router.Response{
			Err:    nil,
			Method: &method,
			Params: nil,
		})
	}
}
