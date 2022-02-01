package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/reservation"
	"github.com/wascript3r/autonuoma/pkg/session"
	sessionHandler "github.com/wascript3r/autonuoma/pkg/session/delivery/http"
	"github.com/wascript3r/autonuoma/pkg/user"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/httputil/middleware"
)

type HTTPHandler struct {
	resUcase     reservation.Usecase
	sessionUcase session.Usecase
	sessionMid   sessionHandler.Middleware
}

func NewHTTPHandler(ctx context.Context, r *httprouter.Router, su session.Usecase, sm sessionHandler.Middleware, auth *middleware.StackCtx, ru reservation.Usecase) {
	handler := &HTTPHandler{
		resUcase:     ru,
		sessionUcase: su,
		sessionMid:   sm,
	}

	r.POST("/api/reservation/create", auth.Wrap(ctx, handler.Create))
	r.POST("/api/reservation/cancel", auth.Wrap(ctx, handler.Cancel))
	r.GET("/api/reservation", auth.Wrap(ctx, handler.GetCurrent))
}

func serveError(w http.ResponseWriter, err error) {
	if err == user.InvalidInputError {
		httpjson.BadRequestCustom(w, user.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, user.UnknownError)
	if code == user.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}
	req := &reservation.CreateReq{}

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	id, err := h.resUcase.Create(r.Context(), req, s.UserID)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, id)
}

func (h *HTTPHandler) Cancel(_ context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &reservation.CancelReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	err = h.resUcase.Cancel(r.Context(), req.ReservationID)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, nil)
}

func (h *HTTPHandler) GetCurrent(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	res, err := h.resUcase.GetCurrent(r.Context(), s.UserID)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
