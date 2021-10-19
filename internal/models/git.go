package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	memory "github.com/go-git/go-git/v5/storage/memory"
)

type Git struct {
	GitUsername   string
	GitPassword   string
	GitUrl        string
	gitRepo       *git.Repository
	gitRemote     *git.Remote
	gitBranch     map[string][]string
	gitAuth       *http.BasicAuth
	gitStorage    *memory.Storage
	gitFileSystem billy.Filesystem
}

func (g *Git) GetHeadBranch() string {
	b, _ := g.gitRepo.Head()
	return string(b.Name())
}

func (g *Git) Initialize() error {
	g.SetGitAuth()
	g.SetGitStorage()
	g.SetGitFileSystem()
	g.NoCheckoutClone()
	g.SetRemote()
	return nil
}

func (g *Git) SetGitAuth() error {
	gitAuth := &http.BasicAuth{
		Username: g.GitUsername,
		Password: g.GitPassword,
	}
	g.gitAuth = gitAuth
	return nil
}

func (g *Git) SetGitStorage() error {
	g.gitStorage = memory.NewStorage()
	return nil
}

func (g *Git) SetGitFileSystem() error {
	g.gitFileSystem = memfs.New()
	return nil
}

func (g *Git) NoCheckoutClone() error {
	gitRepo, err := git.Clone(
		g.gitStorage,
		g.gitFileSystem,
		&git.CloneOptions{
			URL:        g.GitUrl,
			Auth:       g.gitAuth,
			NoCheckout: true,
		},
	)
	if err != nil {
		fmt.Println(err)
		log.Println("error cloning repository", err)
		return err
	}
	g.gitRepo = gitRepo
	return nil
}

func (g *Git) SetRemote() error {
	gitRemote, err := g.gitRepo.Remote("origin")
	if err != nil {
		fmt.Println(err)
		log.Println("error setting remote", err)
		return err
	}
	g.gitRemote = gitRemote
	return nil
}

func (g *Git) ListBranches() (map[string][]string, error) {
	refPrefix := "refs/heads/"
	branchList := make(map[string][]string)
	refList, err := g.gitRemote.List(&git.ListOptions{Auth: g.gitAuth})
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
	g.gitBranch = branchList

	return branchList, nil
}

func (g *Git) PublishChanges(fileName string, renderedTemplateFile string, gitBranchDropDown string) error {
	w, err := g.gitRepo.Worktree()
	if err != nil {
		fmt.Println(err)
		log.Println("error in checking out to branch:", err)
		return err
	}

	err = g.gitRemote.Fetch(&git.FetchOptions{
		Auth:     g.gitAuth,
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
	})
	if err != nil {
		fmt.Println(err)
		log.Println("error in fetching:", err)
		return err
	}

	branchName := fmt.Sprintf("refs/heads/%s", gitBranchDropDown)
	branchRef := plumbing.ReferenceName(branchName)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
		Create: false,
	})
	if err != nil {
		fmt.Println(err)
		log.Println("error in checking out to branch:", err)
		return err
	}

	filePath := fileName
	newFile, err := g.gitFileSystem.Create(filePath)
	if err != nil {
		fmt.Println(err)
		log.Println("error creating new file:", err)
		return err
	}
	newFile.Write([]byte(renderedTemplateFile))
	newFile.Close()
	w.Add(filePath)
	w.Commit("added using zeus", &git.CommitOptions{})
	err = g.gitRepo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       g.gitAuth,
	})
	if err != nil {
		fmt.Println(err)
		log.Println("error pushing in git:", err)
		return err
	}

	return nil
}
