package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/intrigues/zeus-automation/internal/models"
)

func (m *Repository) FetchGitBranch(w http.ResponseWriter, r *http.Request) {
	gitUsername := r.FormValue("git_username")
	gitPassword := r.FormValue("git_password")
	gitUrl := r.FormValue("git_url")
	w.Header().Set("Content-Type", "application/json")

	gitRepo := &models.Git{
		GitUsername: gitUsername,
		GitPassword: gitPassword,
		GitUrl:      gitUrl,
		Directory:   "/tmp/zeus/",
	}
	gitRepo.Initialize()
	listOfBranches, err := gitRepo.ListBranches()
	if err != nil {
		m.App.ErrorLog.Println("error listing branch")
	}
	m.App.Session.Put(r.Context(), "gitRepo", &gitRepo)

	jsonResp, err := json.Marshal(listOfBranches)
	if err != nil {
		m.App.ErrorLog.Printf("error listing branch: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}
