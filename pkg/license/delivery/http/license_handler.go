package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/license"
	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
	"github.com/wascript3r/httputil/middleware"
)

type HTTPHandler struct {
	licenseUcase license.Usecase
	sessionUcase session.Usecase
}

func NewHTTPHandler(ctx context.Context, r *httprouter.Router, agent *middleware.StackCtx, client *middleware.StackCtx, lu license.Usecase, su session.Usecase) {
	handler := &HTTPHandler{
		licenseUcase: lu,
		sessionUcase: su,
	}

	r.POST("/api/agent/license/confirm", agent.Wrap(ctx, handler.ConfirmLicense))
	r.POST("/api/agent/license/reject", agent.Wrap(ctx, handler.RejectLicense))
	r.GET("/api/agent/licenses", agent.Wrap(ctx, handler.AllLicenses))
	r.POST("/api/agent/license/photos", agent.Wrap(ctx, handler.AllPhotos))
	r.POST("/api/license", client.Wrap(ctx, handler.Upload))
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

func (h *HTTPHandler) AllLicenses(_ context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.licenseUcase.GetAllUnconfirmed(r.Context())
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) AllPhotos(_ context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &license.GetPhotosReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.licenseUcase.GetPhotos(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) Upload(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.sessionUcase.LoadCtx(ctx)
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	file, _, err := r.FormFile(license.LicenseKey)
	if err != nil {
		httpjson.BadRequest(w, r.PostForm)
		return
	}
	defer file.Close()

	licenseExpirationDate, err := time.Parse("2006-01-02", r.PostFormValue(license.LicenseExpirationDateKey))
	if err != nil {
		httpjson.BadRequest(w, "data")
		return
	}

	licenseNumber := r.PostFormValue(license.LicenseNumberKey)
	if len(licenseNumber) <= 0 {
		httpjson.BadRequest(w, "number")
		return
	}

	res, err := h.licenseUcase.Upload(ctx, &license.UploadReq{File: file, LicenseExpirationDate: licenseExpirationDate, LicenseNumber: licenseNumber, Uid: s.UserID})
	if err != nil {
		httpjson.InternalError(w, nil)
		return
	}

	httpjson.ServeJSON(w, res)
}
