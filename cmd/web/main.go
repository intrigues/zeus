package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/intrigues/zeus-automation/internal/config"
	"github.com/intrigues/zeus-automation/internal/handlers"
	"github.com/intrigues/zeus-automation/internal/helpers"
	"github.com/intrigues/zeus-automation/internal/models"
	"github.com/intrigues/zeus-automation/internal/render"
	"golang.org/x/crypto/bcrypt"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var addr string

// main is the main function
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. Falling back to default configurations.")
	}

	// get environment variables
	host := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8080"
	}

	addr = fmt.Sprintf("%s:%s", host, port)

	err = run()
	if err != nil {
		log.Fatal("Error starting theapplication", err)
	}

	log.Printf("Staring application on %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	// err = srv.ListenAndServeTLS("localhost.pem", "localhost-key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// get environment variables
	runEnv := os.Getenv("RUN_ENV")
	createAdminUser := os.Getenv("CREATE_ADMIN_USER")

	// database initialization
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error initializing database")
		return err
	}
	db.AutoMigrate(&models.Users{}, &models.AutomationTemplates{})
	app.DB = db
	if createAdminUser == "TRUE" {
		createDefaultUser()
	}

	// session store
	gob.Register(models.Users{})
	gob.Register(models.Git{})

	// change this to true when in production
	if runEnv == "PRODUCTION" {
		app.InProduction = true
		log.Println("Application is running in production mode")
	} else {
		app.InProduction = false
		log.Println("Application is running in development mode")
	}

	// logger setup
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return nil
}

func createDefaultUser() {
	defaultAdminUsername := os.Getenv("DEFAULT_ADMIN_USERNAME")
	defaultAdminEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
	defaultAdminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")

	if defaultAdminUsername == "" {
		defaultAdminUsername = "admin"
	}
	if defaultAdminEmail == "" {
		defaultAdminEmail = "admin@example.com"
	}
	if defaultAdminPassword == "" {
		defaultAdminPassword, err := helpers.GenerateRandomString(24)
		if err != nil {
			log.Println("Error creating random string")
		}
		log.Println(fmt.Sprintf("Admin Password: %s", string(defaultAdminPassword)))
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), 0)
	if err != nil {
		log.Println("Error generating default admin password hash.")
	}
	app.DB.Create(&models.Users{
		Username:          defaultAdminUsername,
		Email:             defaultAdminEmail,
		Password:          string(password_hash),
		IncorrectPassword: 0,
		Status:            1,
		Role:              "Admin",
	})
}
