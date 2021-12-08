package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/review"
	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/httputil/middleware"
)

type HTTPHandler struct {
	reviewUcase  review.Usecase
	sessionUcase session.Usecase
}

func NewHTTPHandler(ctx context.Context, r *httprouter.Router, client *middleware.StackCtx, ru review.Usecase, su session.Usecase) {
	handler := &HTTPHandler{
		reviewUcase:  ru,
		sessionUcase: su,
	}

	r.POST("/api/client/ticket/review/submit", client.Wrap(ctx, handler.SubmitReview))
}

func serveError(w http.ResponseWriter, err error) {
	if err == review.InvalidInputError {
		httpjson.BadRequestCustom(w, review.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, review.UnknownError)
	if code == review.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) SubmitReview(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	req := &review.CreateReq{}

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	err = h.reviewUcase.Submit(r.Context(), s.UserID, s.RoleID, req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, nil)
}
