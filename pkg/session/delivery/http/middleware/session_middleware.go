package middleware

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	"github.com/wascript3r/httputil"
	httpjson "github.com/wascript3r/httputil/json"
)

type HTTPMiddleware struct {
	cookieName     string
	cookieLifetime time.Duration
	secureCookie   bool

	sessionUcase session.Usecase
}

func NewHTTPMiddleware(cookieName string, cookieLifetime time.Duration, secureCookie bool, su session.Usecase) *HTTPMiddleware {
	return &HTTPMiddleware{
		cookieName:     cookieName,
		cookieLifetime: cookieLifetime,
		secureCookie:   secureCookie,

		sessionUcase: su,
	}
}

func (h *HTTPMiddleware) ExtractSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(h.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (h *HTTPMiddleware) SetSessionCookie(w http.ResponseWriter, ss *domain.Session) {
	cookie := &http.Cookie{
		Name:     h.cookieName,
		Value:    url.QueryEscape(ss.ID),
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		MaxAge:   int(h.cookieLifetime.Seconds()),
	}
	http.SetCookie(w, cookie)
}

func (h *HTTPMiddleware) DeleteSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		MaxAge:   0,
	}
	http.SetCookie(w, cookie)
}

func (h *HTTPMiddleware) Authenticated(next httputil.HandleCtx) httputil.HandleCtx {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sessionID, err := h.ExtractSessionID(r)
		if err != nil {
			httpjson.ForbiddenCustom(w, session.NotAuthenticatedError, nil)
			return
		}

		s, err := h.sessionUcase.Validate(ctx, sessionID)
		if err != nil {
			code := errcode.UnwrapErr(err, session.UnknownError)

			if err == session.NotAuthenticatedError || err == session.SessionExpiredError {
				if err == session.SessionExpiredError {
					h.DeleteSessionCookie(w)
				}

				httpjson.ForbiddenCustom(w, code, nil)
				return
			}

			httpjson.InternalErrorCustom(w, code, nil)
			return
		}
		ctx = h.sessionUcase.StoreCtx(ctx, s)

		next(ctx, w, r, p)
	}
}

func (h *HTTPMiddleware) NotAuthenticated(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sessionID, err := h.ExtractSessionID(r)
		if err != nil {
			next(w, r, p)
			return
		}

		_, err = h.sessionUcase.Validate(r.Context(), sessionID)
		if err != nil {
			if err == session.NotAuthenticatedError || err == session.SessionExpiredError {
				if err == session.SessionExpiredError {
					h.DeleteSessionCookie(w)
				}

				next(w, r, p)
				return
			}

			httpjson.InternalErrorCustom(w, errcode.UnwrapErr(err, session.UnknownError), nil)
			return
		}

		httpjson.ServeErr(w, session.AlreadyAuthenticatedError, nil)
	}
}

func (h *HTTPMiddleware) HasRole(role domain.Role) func(next httputil.HandleCtx) httputil.HandleCtx {
	return func(next httputil.HandleCtx) httputil.HandleCtx {
		return h.Authenticated(
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				s, err := h.sessionUcase.LoadCtx(ctx)
				if err != nil {
					httpjson.InternalError(w, nil)
					return
				}

				if !domain.HasRole(s, role) {
					httpjson.ForbiddenCustom(w, session.InsufficientPermissionsError, nil)
					return
				}

				next(ctx, w, r, p)
			},
		)
	}
}
