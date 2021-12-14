package http

import (
	"context"
	"encoding/json"
	"github.com/wascript3r/autonuoma/pkg/trip"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/user"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/httputil/middleware"
)

type HTTPHandler struct {
	tripUsecase trip.Usecase
}

func NewHTTPHandler(ctx context.Context, r *httprouter.Router, auth *middleware.StackCtx, tu trip.Usecase) {
	handler := &HTTPHandler{
		tripUsecase: tu,
	}

	r.POST("/api/trip/start", auth.Wrap(ctx, handler.StartTrip))
	r.POST("/api/trip/end", auth.Wrap(ctx, handler.EndTrip))
	r.GET("/api/trip/:id", auth.Wrap(ctx, handler.GetById))
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

func (h *HTTPHandler) StartTrip(_ context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &trip.StartReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	id, err := h.tripUsecase.Start(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, id)
}

func (h *HTTPHandler) EndTrip(_ context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &trip.EndReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	err = h.tripUsecase.End(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, nil)
}

func (h *HTTPHandler) GetById(_ context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		serveError(w, err)
	}
	res, err := h.tripUsecase.GetById(r.Context(), id)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
