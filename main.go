package main

import (
	"backend/config"
	"backend/controllers"
	"backend/database"
	"backend/middleware"
	"backend/services"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
)

var c = config.GetConfig()
var store = sessions.NewCookieStore([]byte(c.SESSION_KEY))

func waitForDb() (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 1; i < 5; i++ {
		db, err = sql.Open("postgres", c.POSTGRES_URL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
			log.Printf("Database was not available, waiting %s...", time.Second<<i)
			time.Sleep(time.Second << i)
		}
	}
	return nil, fmt.Errorf("could not reach database %v", err)
}

func waitForMigration() (*migrate.Migrate, error) {
	var m *migrate.Migrate
	var err error
	for i := 1; i < 5; i++ {
		m, err = migrate.New("file:///migrations", c.POSTGRES_URL)
		if err != nil {
			log.Printf("Database was not available for migrations waiting %s...", time.Second<<i)
			time.Sleep(time.Second << i)
			continue
		}
		return m, nil
	}
	return nil, fmt.Errorf("failed to create migration %v", err)
}

func main() {
	log.Println("Starting to initialize media gateway server")
	db, err := waitForDb()
	if err != nil {
		log.Fatalf("Could not connect to the DB: %v\n", err)
	}
	defer db.Close()

	log.Println("Initialized database connection")

	m, err := waitForMigration()

	if err != nil {
		log.Fatalf("Could not create migrate instance: %v\n", err)
	}

	log.Println("Running migrations...")

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run up migrations: %v\n", err)
	}

	log.Println("Migrations succesful")

	// session settings
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false,
	}

	md := database.NewMetaDataDb()
	ms := services.NewMetaDataService(db, md)
	ud := database.NewUsersDb()
	fd := database.NewFoldersDatabase()

	fs := services.NewFoldersService(db, fd)
	us := services.NewUserService(db, ud, fd)

	fileService, err := services.NewFileservice(md, db)

	if err != nil {
		log.Fatalf("Failed to initialize file service %v", err)
	}

	fileController := controllers.NewFileController(fileService, ms, db)

	ac := controllers.NewAuthController(us, store)
	vc := controllers.NewViewsController(fs)
	fc := controllers.NewFoldersController(fs)

	if vc == nil {
		log.Fatalf("Failed to initialize view controller")
	}
	fileserver := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", vc.HandleGetRoot)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		os.ReadFile("")
	})

	http.Handle("/folders/", middleware.NewEnsureAuth(us, store, fc.HandleFolders))
	http.Handle("/files/", middleware.NewEnsureAuth(us, store, fileController.HandleFiles))

	http.HandleFunc("/login", vc.HandleGetLogin)
	http.HandleFunc("/login/github", ac.HandleGithubLogin)
	http.HandleFunc("/login/github/callback", ac.HandleGithubCallback)

	log.Println("Starting a server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
