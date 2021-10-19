package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/intrigues/zeus-automation/internal/forms"
	"github.com/intrigues/zeus-automation/internal/helpers"
	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
)

func (m *Repository) GetAutomationAll(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "automation.all.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) GetAutomationNewOpt(w http.ResponseWriter, r *http.Request) {
	currentUserID := m.App.Session.Get(r.Context(), "currentuser").(models.Users).ID
	var automationTemplates []models.AutomationTemplates
	m.App.DB.Where("user_id = ?", currentUserID).Find(&automationTemplates)
	automationTemplateStruct := make(map[string][]string)
	for _, automation := range automationTemplates {
		automationTemplateStruct[automation.ProjectName] = append(automationTemplateStruct[automation.ProjectName], automation.Technology)
	}
	data := make(map[string]interface{})
	data["automationTemplates"] = automationTemplateStruct
	render.RenderTemplate(w, r, "automation.opt.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) CreateAutomationNew(w http.ResponseWriter, r *http.Request) {
	currentUserID := m.App.Session.Get(r.Context(), "currentuser").(models.Users).ID
	projectName := chi.URLParam(r, "projectName")
	technology := chi.URLParam(r, "technology")

	var automationTemplates models.AutomationTemplates
	m.App.DB.Where("user_id = ? AND project_name = ? AND technology = ?", currentUserID, projectName, technology).First(&automationTemplates)

	var listOfVariables []models.AutomationMetadata

	json.Unmarshal([]byte(automationTemplates.TemplateMetaData), &listOfVariables)

	data := make(map[string]interface{})
	data["formVariables"] = listOfVariables
	render.RenderTemplate(w, r, "automation.new.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostCreateAutomationNew(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	currentUserID := m.App.Session.Get(r.Context(), "currentuser").(models.Users).ID
	projectName := chi.URLParam(r, "projectName")
	technology := chi.URLParam(r, "technology")

	fileName := chi.URLParam(r, "filenameField")
	gitBranchDropDown := r.Form.Get("gitBranchDropDown")

	var automationTemplates models.AutomationTemplates
	m.App.DB.Where("user_id = ? AND project_name = ? AND technology = ?", currentUserID, projectName, technology).First(&automationTemplates)
	var listOfVariables []models.AutomationMetadata
	json.Unmarshal([]byte(automationTemplates.TemplateMetaData), &listOfVariables)

	// Form operation
	form := forms.New(r.PostForm)
	// validating form content
	form.Required("gitUrlField", "usernameField", "passwordField", "filenameField")
	for _, variableName := range listOfVariables {
		form.Required(variableName.Name)
	}
	// if form is not valid
	if !form.Valid() {
		var automationTemplates models.AutomationTemplates
		m.App.DB.Where("user_id = ? AND project_name = ? AND technology = ?", currentUserID, projectName, technology).First(&automationTemplates)

		var listOfVariables []models.AutomationMetadata

		json.Unmarshal([]byte(automationTemplates.TemplateMetaData), &listOfVariables)

		data := make(map[string]interface{})
		data["formVariables"] = listOfVariables

		render.RenderTemplate(w, r, "automation.new.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// rendering jenkinfile with user inputs
	renderedTemplateFile := string(automationTemplates.TemplateFile)
	for _, variableName := range listOfVariables {
		m1 := regexp.MustCompile("@@" + variableName.Name + "@@")
		renderedTemplateFile = m1.ReplaceAllString(renderedTemplateFile, form.Get(variableName.Name))
	}

	//commiting the file
	gitRepo := m.App.Session.Get(r.Context(), "gitRepo").(models.Git)
	m.App.ErrorLog.Println("Branch --->", gitRepo.GetHeadBranch())

	// err = gitRepo.PublishChanges(fileName, renderedTemplateFile, gitBranchDropDown)
	if err != nil {
		m.App.ErrorLog.Println("error publishing the changes to git", err)
	}
	m.App.InfoLog.Println("changes successfully pushed to git")
	// redirecting on success
	http.Redirect(w, r, "/admin/automation/opt", http.StatusSeeOther)
}
