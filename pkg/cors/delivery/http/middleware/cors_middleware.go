package middleware

import (
	"net/http"
)

type HTTPMiddleware struct {
	origin string
}

func NewHTTPMiddleware(origin string) *HTTPMiddleware {
	return &HTTPMiddleware{origin}
}

func (h *HTTPMiddleware) EnableCors(hnd http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", h.origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		hnd.ServeHTTP(w, r)
	})
}
