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
	ticketMid    Middleware
	sessionUcase session.Usecase
	roomUcase    room.Usecase

	socketPool *pool.Pool
}

func NewWSHandler(r *router.Router, client *middleware.Stack, agent *middleware.Stack, tu ticket.Usecase, teb ticket.EventBus, tm Middleware, su session.Usecase, ru room.Usecase, socketPool *pool.Pool) {
	handler := &WSHandler{
		ticketUcase:  tu,
		ticketMid:    tm,
		sessionUcase: su,
		roomUcase:    ru,

		socketPool: socketPool,
	}

	teb.Subscribe(ticket.NewTicketEvent, handler.TicketNotification("ticket/notification"))
	teb.Subscribe(ticket.AcceptedTicketEvent, handler.TicketNotification("ticket/notification"))
	teb.Subscribe(ticket.EndedTicketEvent, handler.TicketNotification("ticket/notification"))

	teb.Subscribe(ticket.AcceptedTicketEvent, handler.TicketRoomNotification("ticket/notification/accepted"))
	teb.Subscribe(ticket.EndedTicketEvent, handler.TicketRoomNotification("ticket/notification/ended"))

	r.HandleMethod("ticket/new", client.Wrap(handler.NewTicket))
	r.HandleMethod("ticket/accept", agent.Wrap(handler.AcceptTicket))
	r.HandleMethod("ticket/client/end", client.Wrap(handler.ClientEndTicket))
	r.HandleMethod("ticket/agent/end", agent.Wrap(handler.AgentEndTicket))
	r.HandleMethod("ticket/client/open", client.Wrap(handler.ClientOpenTicket))
	r.HandleMethod("ticket/agent/open", agent.Wrap(handler.AgentOpenTicket))
	r.HandleMethod("ticket/client/close", client.Wrap(handler.CloseTicket))
	r.HandleMethod("ticket/agent/close", agent.Wrap(handler.CloseTicket))
	r.HandleMethod("tickets", agent.Wrap(handler.AllTickets))
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

func (w *WSHandler) ClientEndTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	err = w.ticketUcase.ClientEnd(ctx, ss.UserID)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) AgentEndTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &ticket.AgentEndReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	err = w.ticketUcase.AgentEnd(ctx, ss.UserID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) ClientOpenTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	res, err := w.ticketUcase.ClientGetMessages(ctx, ss.UserID)
	if err != nil {
		serveError(s, r, err)
		return
	}

	if !res.Ticket.Ended {
		err = w.ticketMid.CreateOrRejoinRoom(s, res.Ticket.ID)
		if err != nil {
			serveError(s, r, err)
			return
		}
	}

	router.WriteRes(s, &r.Method, res)
}

func (w *WSHandler) AgentOpenTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	req := &ticket.AgentGetMessagesReq{}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	res, err := w.ticketUcase.AgentGetMessages(ctx, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	if !res.Ticket.Ended {
		err = w.ticketMid.CreateOrRejoinRoom(s, res.Ticket.ID)
		if err != nil {
			serveError(s, r, err)
			return
		}
	}

	router.WriteRes(s, &r.Method, res)
}

func (w *WSHandler) CloseTicket(_ context.Context, s *gows.Socket, r *router.Request) {
	err := w.ticketMid.LeaveCurrentRoom(s)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) AllTickets(ctx context.Context, s *gows.Socket, r *router.Request) {
	res, err := w.ticketUcase.GetTickets(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, res)
}

func (w *WSHandler) TicketNotification(method string) func(context.Context, int) {
	return func(ctx context.Context, _ int) {
		rName, err := w.roomUcase.GetName(domain.AgentRoom)
		if err != nil {
			return
		}

		res, err := w.ticketUcase.GetTickets(ctx)
		if err != nil {
			return
		}

		w.socketPool.EmitRoom(pool.RoomName(rName), &router.Response{
			Err:    nil,
			Method: &method,
			Params: res,
		})
	}
}

func (w *WSHandler) TicketRoomNotification(method string) func(context.Context, int) {
	return func(ctx context.Context, ticketID int) {
		rName := w.ticketMid.GetRoomName(ticketID)

		w.socketPool.EmitRoom(rName, &router.Response{
			Err:    nil,
			Method: &method,
			Params: nil,
		})
	}
}
