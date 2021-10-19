package constant

import "os"

func GetDataDirectory() string {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/var/data/zeus"
	}
	return dataDir
}
