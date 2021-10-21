package constant

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dataDir string

func TestMain(m *testing.M) {
	dataDir = "./data"
	os.Setenv("DATA_DIR", dataDir)
}

func TestGetDataDir(t *testing.T) {
	result := GetDataDir()
	assert.Equal(t, result, dataDir)
}

func TestGetTemplateDir(t *testing.T) {
	id := "9987345989430289"
	result := GetTemplateDir(id)
	assert.Equal(t, result, fmt.Sprintf("./%s/templates/%s", dataDir, id))
}
