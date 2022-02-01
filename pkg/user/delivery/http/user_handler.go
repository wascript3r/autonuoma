package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/session"
	sessionHandler "github.com/wascript3r/autonuoma/pkg/session/delivery/http"
	"github.com/wascript3r/autonuoma/pkg/user"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/httputil/middleware"
)

type HTTPHandler struct {
	userUcase    user.Usecase
	sessionUcase session.Usecase
	sessionMid   sessionHandler.Middleware
}

func NewHTTPHandler(ctx context.Context, r *httprouter.Router, auth *middleware.StackCtx, notAuth *middleware.Stack, uu user.Usecase, su session.Usecase, sm sessionHandler.Middleware) {
	handler := &HTTPHandler{
		userUcase:    uu,
		sessionUcase: su,
		sessionMid:   sm,
	}

	r.POST("/api/user/register", notAuth.Wrap(handler.RegisterUser))
	r.POST("/api/user/authenticate", notAuth.Wrap(handler.AuthenticateUser))
	r.GET("/api/user/token", auth.Wrap(ctx, handler.GetToken))
	r.GET("/api/user/logout", auth.Wrap(ctx, handler.LogoutUser))
	r.GET("/api/user/info", auth.Wrap(ctx, handler.UserInfo))
	r.GET("/api/user", auth.Wrap(ctx, handler.UserData))
	r.POST("/api/user/update", auth.Wrap(ctx, handler.UpdateUser))
	r.GET("/api/user/trips", auth.Wrap(ctx, handler.GetTrips))
	r.GET("/api/user/payment", auth.Wrap(ctx, handler.Payment))
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

func (h *HTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &user.CreateReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	err = h.userUcase.Create(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, nil)
}

func (h *HTTPHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &user.AuthenticateReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	s, res, err := h.userUcase.Authenticate(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}
	h.sessionMid.SetSessionCookie(w, s)

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) GetToken(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	res, err := h.userUcase.GetTempToken(s)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) LogoutUser(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	err = h.userUcase.Logout(r.Context(), s)
	if err != nil {
		serveError(w, err)
		return
	}
	h.sessionMid.DeleteSessionCookie(w)

	httpjson.ServeJSON(w, nil)
}

func (h *HTTPHandler) UserInfo(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	res := h.userUcase.GetInfo(s.UserID, s.RoleID)
	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) UserData(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	res, err := h.userUcase.GetData(ctx, s.UserID)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) UpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	req := &user.UpdateReq{}

	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.userUcase.UpdateUser(ctx, s.UserID, req)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) GetTrips(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	res, err := h.userUcase.GetTrips(ctx, s.UserID)
	fmt.Println(err)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) Payment(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	res, err := h.userUcase.CheckPayment(ctx, s.UserID)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	httpjson.ServeJSON(w, res)
}
