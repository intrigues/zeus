package constant

import (
	"fmt"
	"os"
)

func GetDataDir() string {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/var/data/zeus"
	}
	return dataDir
}

func GetDatabaseDir() string {
	dbDir := fmt.Sprintf("%s/db/app.db", GetDataDir())
	return dbDir
}

func GetTemplateDir(id string) string {
	templateDir := fmt.Sprintf("%s/templates/%s", GetDataDir(), id)
	return templateDir
}

func GetGitRepoDir(id string) string {
	return fmt.Sprintf("%s/%s/%s", GetDataDir(), "gitRepos", id)
}

func GetAutomationBranchPrefix() string {
	prefix := os.Getenv("AUTOMATION_NEW_BRANCH_PREFIX")
	if prefix == "" {
		prefix = "zeus-automation"
	}
	return prefix
}

func CreateAutomationBranchName(name string) string {
	return fmt.Sprintf("%s-%s", GetAutomationBranchPrefix(), name)
}
