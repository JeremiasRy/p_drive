package main

import (
	"backend/controllers"
	"backend/database"
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
)

var (
	POSTGRES_URL      = fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_USER)
	POSTGRES_USER     = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_LOCATION = os.Getenv("POSTGRES_LOCATION")
)

func waitForDb() (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 1; i < 5; i++ {
		db, err = sql.Open("postgres", POSTGRES_URL)
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
		m, err = migrate.New("file:///migrations", POSTGRES_URL)
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

	ud := database.NewUsersDb(db)
	us := services.NewUserService(ud)

	fs, err := services.NewFileservice()
	fc := controllers.NewFileController(fs)

	ac := controllers.NewAuthController(us)

	if err != nil {
		log.Fatalf("Failed to create fileservice client %s", err)
	}

	vc := controllers.NewViewsController()

	if vc == nil {
		log.Fatalf("Failed to initialize view controller")
		os.Exit(1)
	}

	http.HandleFunc("/", vc.HandleGetMain)
	http.HandleFunc("/files", vc.HandleGetFiles)
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				fc.HandlePostUpload(w, r)
			}
		case http.MethodGet:
			{
				vc.HandleGetUpload(w, r)
			}
		default:
			{
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
	http.HandleFunc("/login", vc.HandleGetLogin)
	http.HandleFunc("/login/github", ac.HandleGithubLogin)
	http.HandleFunc("/login/github/callback", ac.HandleGithubCallback)

	log.Println("Starting a server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
