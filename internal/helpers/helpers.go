package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"github.com/intrigues/zeus-automation/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError returns error when there is a server error
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError returns error when there is a server error
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// GenerateRandomString generates random string using UUID
// Max length for random string is 32
func GenerateRandomString(length int) string {
	if length > 32 {
		log.Println("Max length can be 32. falling back to length 32.")
		length = 32
	}
	u := uuid.NewString()
	u = strings.Replace(u, "-", "", -1)
	return u[:length]
}

// SaveFile helps to save file on the given path
func SaveFile(f string, path string, fileName string) error {
	location := fmt.Sprintf("%s/%s", path, fileName)
	MakeDirectory(filepath.Dir(location))
	err := os.WriteFile(location, []byte(f), 0644)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating file at: %s", path))
		return err
	}
	return nil
}

func MakeDirectory(path string) error {
	// TODO: fix permission to make it more secure
	isExists, err := exists(path)
	if err != nil {
		log.Println("Error listing directory", err)
		return err
	}
	if !isExists {
		err := os.MkdirAll(path, 0775)
		if err != nil {
			log.Println("Error creating directory", err)
			return err
		}
	}
	return nil
}

func DeleteDirectory(path string) error {
	isExists, err := exists(path)
	if err != nil {
		log.Println("Error listing directory", err)
		return err
	}
	if isExists {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println("Error deleting directory", err)
			return err
		}
	}
	return nil
}

func ReadFile(path string, fileName string) (string, error) {
	var data []byte
	location := fmt.Sprintf("%s/%s", path, fileName)
	isExists, err := exists(location)
	if err != nil {
		log.Println("Error listing directory", err)
		return "", err
	}
	if isExists {
		data, err = os.ReadFile(location)
		if err != nil {
			log.Println("Error deleting directory", err)
			return "", err
		}
	}
	return string(data), nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
