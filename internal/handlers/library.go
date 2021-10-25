package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	appconst "github.com/intrigues/zeus-automation/internal/constant"
	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
)

func (m *Repository) GetLibraryAll(w http.ResponseWriter, r *http.Request) {
	var listOfTemplates []map[string]string
	for _, template := range retriveTemplateLibrary() {
		if strings.HasSuffix(template, ".mapping") {
			s := strings.SplitN(template, "/", 3)
			t := map[string]string{"project": s[0], "technology": s[1], "file_name": s[2]}
			listOfTemplates = append(listOfTemplates, t)
		}
	}
	data := make(map[string]interface{})
	data["listOfTemplates"] = listOfTemplates
	render.RenderTemplate(w, r, "library.all.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) DownloadLibraryTemplate(w http.ResponseWriter, r *http.Request) {
	// currentUserID := m.App.Session.Get(r.Context(), "currentuser").(models.Users).ID
	// projectName := chi.URLParam(r, "projectName")
	// technology := chi.URLParam(r, "technology")

	// downloadFileContent(projectName, technology)

	http.Redirect(w, r, "/admin/library/all", http.StatusSeeOther)
}

// retriveTemplateLibrary retrives the list of available templates in template library
// result stores the response in required format for trees api
func retriveTemplateLibrary() []string {
	result := []string{}
	url := fmt.Sprintf(
		"%s/repos/%s/git/trees/main?recursive=1",
		appconst.GetGithubHost(),
		appconst.GetTemplateLibraryRepo(),
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := make(map[string][]map[string]string)
	json.Unmarshal([]byte(string(body)), &data)
	for _, v := range data["tree"] {
		if v["type"] == "blob" {
			result = append(result, v["path"])
		}
	}
	return result
}

// func downloadFileContent(projectName string, technology string) {
// 	url := fmt.Sprintf(
// 		"%s/%s/%s/%s",
// 		appconst.GetRawGithubHost(),
// 		appconst.GetTemplateLibraryRepo(),
// 		projectName,
// 		technology,
// 	)
// }
