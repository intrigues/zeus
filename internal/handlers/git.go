package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/intrigues/zeus-automation/internal/models"
)

func (m *Repository) FetchGitBranch(w http.ResponseWriter, r *http.Request) {
	m.App.InfoLog.Println("Fetching Git Branch")
	gitUsername := r.FormValue("git_username")
	gitPassword := r.FormValue("git_password")
	gitUrl := r.FormValue("git_url")
	w.Header().Set("Content-Type", "application/json")

	gitRepo := &models.Git{
		GitUsername: gitUsername,
		GitPassword: gitPassword,
		GitUrl:      gitUrl,
	}

	gitRepo.Initialize()
	m.App.InfoLog.Println("Fetching Git Branch: git initialized")
	listOfBranches, err := gitRepo.ListBranches()
	m.App.InfoLog.Println("Fetching Git Branch: Branch fetched")
	if err != nil {
		m.App.ErrorLog.Println("error in listing branch")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}
	m.App.Session.Put(r.Context(), "gitRepo", &gitRepo)

	jsonResp, err := json.Marshal(listOfBranches)
	if err != nil {
		m.App.ErrorLog.Printf("error in matchaling branch: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}
