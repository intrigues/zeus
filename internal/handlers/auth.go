package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/intrigues/zeus-automation/internal/config"
	"github.com/intrigues/zeus-automation/internal/forms"
	"github.com/intrigues/zeus-automation/internal/helpers"
	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
	"golang.org/x/crypto/bcrypt"
)

// Repo the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) GetHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (m *Repository) GetLogin(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	email := r.Form.Get("emailField")
	password := r.Form.Get("passwordField")

	form := forms.New(r.PostForm)
	form.Required("emailField", "passwordField")
	form.IsEmail("emailField")

	if !form.Valid() {
		render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	var user models.Users
	m.App.DB.First(&user, "email = ?", email)

	// authenticate method
	// TODO: isolate this method in models/controllers
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "currentuser", user)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	m.App.InfoLog.Println("Logged in successfully. username:", user.Username)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (m *Repository) GetSignup(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signup.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.InfoLog.Println(r.Form)

	form := forms.New(r.PostForm)
	form.Required("usernameField", "emailField", "passwordField", "roleSelector")
	form.IsEmail("emailField")
	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Please enter valid details")
		m.App.ErrorLog.Println("error validating form")
		render.RenderTemplate(w, r, "signup.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(r.Form.Get("passwordField")), 0)
	if err != nil {
		m.App.ErrorLog.Println("error generating hash from given password")
	}
	m.App.DB.Create(&models.Users{
		Username:          r.Form.Get("usernameField"),
		Email:             r.Form.Get("emailField"),
		Password:          string(password_hash),
		IncorrectPassword: 0,
		Status:            0,
		Role:              r.Form.Get("roleSelector"),
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (m *Repository) GetLogout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (m *Repository) GetDashboard(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "dashboard.page.tmpl", &models.TemplateData{})
}

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
