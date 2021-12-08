package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/faq"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
)

type HTTPHandler struct {
	faqUcase faq.Usecase
}

func NewHTTPHandler(r *httprouter.Router, fu faq.Usecase) {
	handler := &HTTPHandler{
		faqUcase: fu,
	}

	r.GET("/api/faq", handler.AllFAQ)
}

func serveError(w http.ResponseWriter, err error) {
	if err == faq.InvalidInputError {
		httpjson.BadRequestCustom(w, faq.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, faq.UnknownError)
	if code == faq.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) AllFAQ(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.faqUcase.GetAll(r.Context())
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
