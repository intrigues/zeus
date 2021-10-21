package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	appconst "github.com/intrigues/zeus-automation/internal/constant"
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
	form.Required("projectNameField", "technologyField")
	// Fetching number of template files for validating those fields
	for key, _ := range r.Form {
		if strings.HasPrefix(key, "file") {
			form.Required(key)
		}
	}
	// form.IsJson("templateFileMetadata")
	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Please enter valid details")
		m.App.ErrorLog.Println("Error validating form")
		render.RenderTemplate(w, r, "templates.new.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	m.App.InfoLog.Println("Fields are validated")
	templateID := helpers.GenerateRandomString(20)
	template := &models.AutomationTemplates{
		ID:          templateID,
		ProjectName: form.Get("projectNameField"),
		Technology:  form.Get("technologyField"),
		User:        m.App.Session.Get(r.Context(), "currentuser").(models.Users),
	}
	m.App.DB.Create(template)

	// Fetching number of template files
	fileCounter := 0
	// calculating number of files
	for key, _ := range r.Form {
		if strings.HasPrefix(key, "fileNameField") {
			fileCounter++
		}
	}

	m.App.InfoLog.Println("File numbers are fetched")
	// We are assuming that fileNameField, fileTemplateField and fileMappingField are in pair.
	for index := 1; index <= fileCounter; index++ {
		fileNameField := form.Get(fmt.Sprintf("fileNameField%d", index))
		fileTemplateField := form.Get(fmt.Sprintf("fileTemplateField%d", index))
		fileMappingField := form.Get(fmt.Sprintf("fileMappingField%d", index))
		SaveTemplate(fileTemplateField, templateID, fmt.Sprintf("%s.template", fileNameField))
		SaveTemplate(fileMappingField, templateID, fmt.Sprintf("%s.mapping", fileNameField))
	}

	m.App.InfoLog.Println("Files are saved")

	// SaveTemplate(form.Get("fileTextArea"), templateID, "Jenkinsfile.template")
	// SaveTemplate(form.Get("templateFileMetadata"), templateID, "Jenkinsfile.mapping")
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
	DeleteTemplate(templateID)
	m.App.DB.Delete(&automationTemplates)

	http.Redirect(w, r, "/admin/templates/all", http.StatusSeeOther)
}

// util functions
// SaveTemplate saves templatefiles
func SaveTemplate(f string, id string, fileName string) {
	templateDir := appconst.GetTemplateDir(id)
	helpers.MakeDirectory(templateDir)
	helpers.SaveFile(f, templateDir, fileName)
}

// DeleteTemplate deletes templatefiles
func DeleteTemplate(id string) {
	templateDir := appconst.GetTemplateDir(id)
	helpers.DeleteDirectory(templateDir)
}
