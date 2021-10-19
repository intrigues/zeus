package main

import (
	"net/http"

	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/justinas/nosurf"
)

// NoSurf is the csrf protection middleware
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		// TODO - get this from env
		HttpOnly: false,
		Path:     "/",
		// TODO - get this from env
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Please login to continue")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "currentuser")
	if exists {
		currentUser := app.Session.Get(r.Context(), "currentuser").(models.Users)
		return !(currentUser.Status == 0)
	}
	return exists
}
