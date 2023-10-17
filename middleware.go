package cas

import (
	"log/slog"
	"net/http"
)

// Handler returns a standard http.HandlerFunc, which will check the authenticated status (redirect user go login if needed)
// If the user pass the authenticated check, it will call the h's ServeHTTP method
func (c *Client) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "cas: handling request", "method", r.Method, "url", r.URL)

		setClient(r, c)

		if !IsAuthenticated(r) {
			RedirectToLogin(w, r)
			return
		}

		if r.URL.Path == "/logout" {
			RedirectToLogout(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (c *Client) BranchUnauthenticated(next http.Handler, unauthHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuthenticated(r) {
			slog.InfoContext(r.Context(), "cas: user is authenticated")
			next.ServeHTTP(w, r)
			return
		}
		slog.InfoContext(r.Context(), "cas: user is not authenticated")
		unauthHandler.ServeHTTP(w, r)
	})
}

func (c *Client) BlockUnauthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuthenticated(r) {
			slog.InfoContext(r.Context(), "cas: user is authenticated")
			next.ServeHTTP(w, r)
			return
		}
		slog.InfoContext(r.Context(), "cas: user is not authenticated")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
