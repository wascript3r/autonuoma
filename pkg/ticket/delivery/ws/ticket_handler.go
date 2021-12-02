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

	r.HandleMethod("client/ticket/new", client.Wrap(handler.NewTicket))
	r.HandleMethod("agent/ticket/accept", agent.Wrap(handler.AcceptTicket))

	r.HandleMethod("client/ticket/end", client.Wrap(handler.EndTicket))
	r.HandleMethod("agent/ticket/end", agent.Wrap(handler.EndTicket))

	r.HandleMethod("client/ticket/open", client.Wrap(handler.OpenTicket))
	r.HandleMethod("agent/ticket/open", agent.Wrap(handler.OpenTicket))

	r.HandleMethod("client/ticket/close", client.Wrap(handler.CloseTicket))
	r.HandleMethod("agent/ticket/close", agent.Wrap(handler.CloseTicket))

	r.HandleMethod("client/tickets", client.Wrap(handler.AllTickets))
	r.HandleMethod("agent/tickets", agent.Wrap(handler.AllTickets))
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

func (w *WSHandler) EndTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &ticket.EndReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	err = w.ticketUcase.End(ctx, ss.UserID, ss.RoleID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	router.WriteRes(s, &r.Method, nil)
}

func (w *WSHandler) OpenTicket(ctx context.Context, s *gows.Socket, r *router.Request) {
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	req := &ticket.GetMessagesReq{}

	err = json.Unmarshal(r.Params, req)
	if err != nil {
		router.WriteBadRequest(s, &r.Method)
		return
	}

	res, err := w.ticketUcase.GetMessages(ctx, ss.UserID, ss.RoleID, req)
	if err != nil {
		serveError(s, r, err)
		return
	}

	if res.Ticket.Status != domain.EndedTicketStatus {
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
	ss, err := w.sessionUcase.LoadCtx(ctx)
	if err != nil {
		serveError(s, r, err)
		return
	}

	res, err := w.ticketUcase.GetTickets(ctx, ss.UserID, ss.RoleID)
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

		res, err := w.ticketUcase.GetTickets(ctx, 0, domain.AgentRole)
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
