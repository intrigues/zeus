package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/intrigues/zeus-automation/internal/forms"
	"github.com/intrigues/zeus-automation/internal/helpers"
	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
)

// create new template
func (m *Repository) GetTemplateNew(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "templates.new.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// post method
func (m *Repository) PostTemplateNew(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("projectNameField", "technologyField", "fileTextArea", "templateFileMetadata")
	form.IsJson("templateFileMetadata")
	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Please enter valid details")
		m.App.ErrorLog.Println("Error validating form")
		render.RenderTemplate(w, r, "templates.new.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	m.App.DB.Create(&models.AutomationTemplates{
		ProjectName:      form.Get("projectNameField"),
		Technology:       form.Get("technologyField"),
		TemplateFile:     []byte(form.Get("fileTextArea")),
		TemplateMetaData: []byte(form.Get("templateFileMetadata")),
		User:             m.App.Session.Get(r.Context(), "currentuser").(models.Users),
	})

	http.Redirect(w, r, "/admin/templates/all", http.StatusSeeOther)
}

// get all templates
func (m *Repository) GetTemplateAll(w http.ResponseWriter, r *http.Request) {
	var automationTemplates []models.AutomationTemplates
	m.App.DB.Find(&automationTemplates).Where("user_id = ?", m.App.Session.Get(r.Context(), "currentuser").(models.Users).ID)

	data := make(map[string]interface{})
	data["automationTemplates"] = automationTemplates
	render.RenderTemplate(w, r, "templates.all.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// delete templates
func (m *Repository) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "templateID")
	currentUserID := m.App.Session.Get(r.Context(), "currentuser").(models.Users).ID

	var automationTemplates []models.AutomationTemplates
	m.App.DB.Where("user_id = ? AND id = ?", currentUserID, templateID).Find(&automationTemplates)
	m.App.DB.Delete(&automationTemplates)

	http.Redirect(w, r, "/admin/templates/all", http.StatusSeeOther)
}
