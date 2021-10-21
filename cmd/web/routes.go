package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/intrigues/zeus-automation/internal/config"
	"github.com/intrigues/zeus-automation/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.GetHome)

	mux.Get("/login", handlers.Repo.GetLogin)
	mux.Post("/login", handlers.Repo.PostLogin)

	mux.Get("/signup", handlers.Repo.GetSignup)
	mux.Post("/signup", handlers.Repo.PostSignup)
	mux.Get("/logout", handlers.Repo.GetLogout)

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.GetDashboard)
		mux.Route("/users", func(r chi.Router) {
			r.Get("/all", handlers.Repo.GetUsers)
			r.Route("/{username}", func(r chi.Router) {
				r.Get("/activate", handlers.Repo.ActivateUser)
				r.Get("/deactivate", handlers.Repo.DeactivateUser)
			})
		})
		mux.Route("/automation", func(r chi.Router) {
			r.Get("/opt", handlers.Repo.GetAutomationNewOpt)
			r.Route("/new", func(r chi.Router) {
				r.Get("/{projectName}/{technology}", handlers.Repo.CreateAutomationNew)
				r.Post("/{projectName}/{technology}", handlers.Repo.PostCreateAutomationNew)
			})
		})
		mux.Route("/templates", func(r chi.Router) {
			r.Get("/all", handlers.Repo.GetTemplateAll)
			r.Route("/new", func(r chi.Router) {
				r.Get("/", handlers.Repo.GetTemplateNew)
				r.Post("/", handlers.Repo.PostTemplateNew)

			})
			r.Route("/{templateID}", func(r chi.Router) {
				r.Get("/view", handlers.Repo.GetTemplateAll)
				r.Get("/edit", handlers.Repo.GetTemplateAll)
				r.Get("/delete", handlers.Repo.DeleteTemplate)
			})
		})
		mux.Route("/git", func(r chi.Router) {
			r.Get("/fetch", handlers.Repo.FetchGitBranch)
			r.Post("/fetch", handlers.Repo.FetchGitBranch)
		})
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
