package handlers

import (
	"net/http"

	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
)

func (m *Repository) GetDashboard(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "dashboard.page.tmpl", &models.TemplateData{})
}
