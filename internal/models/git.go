package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	appconst "github.com/intrigues/zeus-automation/internal/constant"
	"github.com/intrigues/zeus-automation/internal/helpers"
)

type Git struct {
	GitUsername string
	GitPassword string
	GitUrl      string
	Directory   string
	GitAuth     *http.BasicAuth
}

func (g *Git) Initialize() error {
	g.SetDirectory()
	g.MakeDirectory()
	g.SetGitAuth()
	g.NoCheckoutClone()
	return nil
}

func (g *Git) SetDirectory() {
	folderName := helpers.GenerateRandomString(20)
	g.Directory = appconst.GetGitRepoDir(folderName)
}

func (g *Git) MakeDirectory() error {
	err := os.MkdirAll(g.Directory, 0775)
	if err != nil {
		log.Println("Error creating directory")
		return err
	}
	return nil
}

func (g *Git) SetGitAuth() error {
	gitAuth := &http.BasicAuth{
		Username: g.GitUsername,
		Password: g.GitPassword,
	}
	g.GitAuth = gitAuth
	return nil
}

func (g *Git) NoCheckoutClone() error {
	_, err := git.PlainClone(
		g.Directory,
		false,
		&git.CloneOptions{
			URL: g.GitUrl,
			Auth: &http.BasicAuth{
				Username: g.GitUsername,
				Password: g.GitPassword,
			},
		},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (g *Git) GetRepo() (*git.Repository, error) {
	gitRepo, err := git.PlainOpen(g.Directory)
	if err != nil {
		log.Println("error in opening git repo directory:", err)
		return nil, err
	}
	return gitRepo, nil
}

func (g *Git) GetRemoteOrigin() (*git.Remote, error) {
	gitRepo, err := g.GetRepo()
	if err != nil {
		return nil, err
	}
	gitRemote, err := gitRepo.Remote("origin")
	if err != nil {
		log.Println("error setting remote", err)
		return nil, err
	}
	return gitRemote, nil
}

func (g *Git) FetchRemote() error {
	gitRemote, err := g.GetRemoteOrigin()
	if err != nil {
		return err
	}
	err = gitRemote.Fetch(&git.FetchOptions{
		Auth:     g.GitAuth,
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
	})
	if err != nil {
		log.Println("error in fetching:", err)
		return err
	}
	return nil
}

func (g *Git) GetTree() (*git.Worktree, error) {
	gitRepo, err := g.GetRepo()
	if err != nil {
		return nil, err
	}

	w, err := gitRepo.Worktree()
	if err != nil {
		log.Println("error in getting work tree", err)
		return nil, err
	}
	return w, nil
}

func (g *Git) CreateNewBranch(branch string) error {
	gitRepo, err := g.GetRepo()
	err = gitRepo.CreateBranch(&config.Branch{
		Name: branch,
	})
	if err != nil {
		log.Println("error in creating new branch:", err)
		return err
	}
	return nil
}

func (g *Git) CheckoutAndCreateNewBranch(newBranch string, oldBranch string) error {
	g.CheckoutToBranch(oldBranch)
	newBranchName := fmt.Sprintf("refs/heads/%s", newBranch)
	newBranchRef := plumbing.ReferenceName(newBranchName)

	w, err := g.GetTree()
	if err != nil {
		log.Println("error in checking out to old branch:", err)
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: newBranchRef,
		Create: true,
	})
	if err != nil {
		log.Println("error in checking out to new branch:", err)
		return err
	}
	return nil
}

func (g *Git) CheckoutToBranch(branch string) error {
	branchName := fmt.Sprintf("refs/heads/%s", branch)
	branchRef := plumbing.ReferenceName(branchName)

	w, err := g.GetTree()
	if err != nil {
		log.Println("error in checking out to branch:", err)
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
		Create: false,
	})
	if err != nil {
		log.Println("error in checking out to branch:", err)
		return err
	}
	return nil
}

func (g *Git) ListBranches() (map[string][]string, error) {
	refPrefix := "refs/heads/"
	branchList := make(map[string][]string)
	gitRemote, err := g.GetRemoteOrigin()
	if err != nil {
		return nil, err
	}
	refList, err := gitRemote.List(&git.ListOptions{Auth: g.GitAuth})
	if err != nil {
		return make(map[string][]string), err
	}
	for _, ref := range refList {
		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchList["branches"] = append(branchList["branches"], refName[len(refPrefix):])
	}
	return branchList, nil
}

func (g *Git) GetListOffiles() ([]string, error) {

	gitRepo, err := g.GetRepo()
	if err != nil {
		return nil, err
	}
	// Get head ref
	head, _ := gitRepo.Head()
	// get latest commit
	commit, _ := gitRepo.CommitObject(head.Hash())
	// get tree object of the latest commit
	tObj, err := commit.Tree()

	if err != nil {
		log.Println("error in listing files:", err)
		return nil, err
	}

	var files []string
	// create list of file names in the latest commit
	for item := range tObj.Entries {
		files = append(files, tObj.Entries[item].Name)
	}

	return files, nil
}

func (g *Git) AddChangesToWorkTree(fileName string, renderedTemplateFile string) error {

	w, err := g.GetTree()
	if err != nil {
		log.Println("error in getting git tree:", err)
		return err
	}

	helpers.MakeDirectory(filepath.Dir(fmt.Sprintf("%s/%s", g.Directory, fileName)))
	newFile, err := os.Create(fmt.Sprintf("%s/%s", g.Directory, fileName))

	if err != nil {
		log.Println("error creating new file:", err)
		return err
	}
	newFile.Write([]byte(renderedTemplateFile))
	newFile.Close()
	w.Add(fileName)

	return nil
}

func (g *Git) CommitAndPush(msg string) error {

	w, err := g.GetTree()
	if err != nil {
		log.Println("error in getting git tree:", err)
		return err
	}

	w.Commit(msg, &git.CommitOptions{})

	gitRepo, err := g.GetRepo()
	if err != nil {
		return err
	}

	err = gitRepo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       g.GitAuth,
	})
	if err != nil {
		log.Println("error pushing in git:", err)
		return err
	}

	return nil
}
