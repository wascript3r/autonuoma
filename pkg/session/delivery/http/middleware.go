package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/httputil"
)

type Middleware interface {
	Authenticated(next httputil.HandleCtx) httputil.HandleCtx
	NotAuthenticated(next httprouter.Handle) httprouter.Handle
	HasRole(role domain.Role) func(next httputil.HandleCtx) httputil.HandleCtx
	SetSessionCookie(w http.ResponseWriter, ss *domain.Session)
	DeleteSessionCookie(w http.ResponseWriter)
}
