package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/intrigues/zeus-automation/internal/forms"
	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
)

func (m *Repository) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.Users
	m.App.DB.Find(&users)

	data := make(map[string]interface{})
	data["users"] = users
	render.RenderTemplate(w, r, "users.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) ActivateUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	m.App.InfoLog.Println("activating user:", username)
	var user models.Users
	m.App.DB.First(&user, "username = ?", username)
	user.Status = 1
	m.App.DB.Save(&user)
	m.App.Session.Put(r.Context(), "flash", "User activated")
	http.Redirect(w, r, "/admin/users/all", http.StatusSeeOther)
}

func (m *Repository) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	m.App.InfoLog.Println("deactivating user:", username)
	var user models.Users
	m.App.DB.First(&user, "username = ?", username)
	user.Status = 0
	m.App.DB.Save(&user)
	m.App.Session.Put(r.Context(), "flash", "User activated")
	http.Redirect(w, r, "/admin/users/all", http.StatusSeeOther)
}
