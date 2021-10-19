package helpers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func GenerateRandomString(length int) (string, error) {
	if length > 32 {
		log.Println("Please enter the length smaller than 32.")
		return "", errors.New("UUIDLengthNotSupported")
	}
	u := uuid.NewString()
	u = strings.Replace(u, "-", "", -1)
	return u[:length], nil
}
