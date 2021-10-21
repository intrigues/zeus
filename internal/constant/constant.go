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

func GetTemplateDir(id string) string {
	templateDir := fmt.Sprintf("%s/templates/%s", GetDataDir(), id)
	return templateDir
}
