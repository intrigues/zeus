package handlers

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	appconst "github.com/intrigues/zeus-automation/internal/constant"
	"github.com/intrigues/zeus-automation/internal/forms"
	"github.com/intrigues/zeus-automation/internal/helpers"
	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
)

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

	// TODO: optimise this in a better way
	// Fetch list of variables for the files
	listOfVariables := make(map[string][]models.AutomationMetadata)
	files := GetTemplateFiles(automationTemplates.ID)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mapping") {
			var t []models.AutomationMetadata
			templateMapping := ReadTemplate(automationTemplates.ID, file.Name())
			json.Unmarshal([]byte(templateMapping), &t)
			filePrefix := strings.Split(file.Name(), ".")[0]
			listOfVariables[filePrefix] = t
		}
	}

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

	// getting current user
	currentUserID := m.App.Session.Get(r.Context(), "currentuser").(models.Users).ID

	// fetching url params
	projectName := chi.URLParam(r, "projectName")
	technology := chi.URLParam(r, "technology")

	// fetching form data
	// gitBranchDropDown := r.Form.Get("gitBranchDropDown")

	// validating and retriving the template
	var automationTemplates models.AutomationTemplates
	m.App.DB.Where("user_id = ? AND project_name = ? AND technology = ?", currentUserID, projectName, technology).First(&automationTemplates)

	// Fetch list of variables for the files
	listOfVariables := make(map[string][]models.AutomationMetadata)
	files := GetTemplateFiles(automationTemplates.ID)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mapping") {
			var t []models.AutomationMetadata
			m.App.InfoLog.Println("Files ---> ", file.Name(), file.IsDir())
			templateMapping := ReadTemplate(automationTemplates.ID, file.Name())
			json.Unmarshal([]byte(templateMapping), &t)

			m.App.InfoLog.Println("t ---> ", t)
			filePrefix := strings.Split(file.Name(), ".")[0]
			listOfVariables[filePrefix] = t
		}
	}

	// Form operation
	form := forms.New(r.PostForm)
	// validating form content
	// TODO: validate dynamic variables from form
	form.Required("gitUrlField", "gitUsernameField", "gitPasswordField")

	// if form is not valid
	if !form.Valid() {
		var automationTemplates models.AutomationTemplates
		m.App.DB.Where("user_id = ? AND project_name = ? AND technology = ?", currentUserID, projectName, technology).First(&automationTemplates)

		listOfVariables := make(map[string][]models.AutomationMetadata)
		files := GetTemplateFiles(automationTemplates.ID)
		for _, file := range files {
			filePrefix := strings.Split(file.Name(), ".")[0]
			if strings.HasSuffix(file.Name(), ".mapping") {
				var t []models.AutomationMetadata
				templateMapping := ReadTemplate(automationTemplates.ID, file.Name())
				json.Unmarshal([]byte(templateMapping), &t)
				listOfVariables[filePrefix] = t
			}
		}

		data := make(map[string]interface{})
		data["formVariables"] = listOfVariables

		render.RenderTemplate(w, r, "automation.new.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	for _, file := range files {
		filePrefix := strings.Split(file.Name(), ".")[0]
		if strings.HasSuffix(file.Name(), ".template") {
			renderedTemplateFile := ReadTemplate(automationTemplates.ID, file.Name())
			for _, variableName := range listOfVariables[filePrefix] {
				m1 := regexp.MustCompile("@@" + variableName.Name + "@@")
				renderedTemplateFile = m1.ReplaceAllString(renderedTemplateFile, form.Get(fmt.Sprintf("%s-%s", filePrefix, variableName.Name)))
			}
		}
	}

	//commiting the file
	// gitRepo := m.App.Session.Get(r.Context(), "gitRepo").(models.Git)

	// _ = gitRepo.PublishChanges(fileName, renderedTemplateFile, gitBranchDropDown)
	// gitRepo.FetchRemote()
	// gitRepo.CheckoutToBranch(gitBranchDropDown)
	// files, err := gitRepo.GetListOffiles()
	// if err != nil {
	// 	m.App.ErrorLog.Println("error publishing the changes to git", err)
	// } else {
	// 	m.App.ErrorLog.Println("files", files)
	// }
	// m.App.InfoLog.Println("changes successfully pushed to git")
	// redirecting on success
	http.Redirect(w, r, "/admin/automation/opt", http.StatusSeeOther)
}

// ReadTemplate reads templatefiles
func ReadTemplate(id string, fileName string) string {
	templateDir := appconst.GetTemplateDir(id)
	data, err := helpers.ReadFile(templateDir, fileName)
	if err != nil {
		fmt.Println("error occured while reading from file", err)
		return "error"
	}
	return data
}

// ReadTemplate reads templatefiles
func GetTemplateFiles(id string) []fs.FileInfo {
	templateDir := appconst.GetTemplateDir(id)
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		log.Fatal(err)
	}
	return files
}
