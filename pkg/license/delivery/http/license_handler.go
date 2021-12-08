package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/license"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/httputil/middleware"
)

type HTTPHandler struct {
	licenseUcase license.Usecase
}

func NewHTTPHandler(ctx context.Context, r *httprouter.Router, agent *middleware.StackCtx, lu license.Usecase) {
	handler := &HTTPHandler{
		licenseUcase: lu,
	}

	r.POST("/api/agent/license/confirm", agent.Wrap(ctx, handler.ConfirmLicense))
	r.POST("/api/agent/license/reject", agent.Wrap(ctx, handler.RejectLicense))
}

func serveError(w http.ResponseWriter, err error) {
	if err == license.InvalidInputError {
		httpjson.BadRequestCustom(w, license.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, license.UnknownError)
	if code == license.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) ConfirmLicense(_ context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &license.ChangeStatusReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	err = h.licenseUcase.Confirm(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, nil)
}

func (h *HTTPHandler) RejectLicense(_ context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &license.ChangeStatusReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	err = h.licenseUcase.Reject(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, nil)
}
