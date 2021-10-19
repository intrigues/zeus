package handlers

import "net/http"

func (m *Repository) GetHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
