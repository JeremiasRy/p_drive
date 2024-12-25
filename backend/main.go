package main

import (
	"backend/controllers"
	"backend/services"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Handle GET request
		fmt.Fprintf(w, "Handling GET request")
	case "POST":
		// Handle POST request
		fmt.Fprintf(w, "Handling POST request")
	default:
		// Handle unsupported methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.Println("Starting to initialize media gateway server")
	db, err := sql.Open("postgres", POSTGRES_URL)
	if err != nil {
		log.Fatalf("Could not connect to the DB: %v\n", err)
	}
	defer db.Close()

	log.Println("Initialized database connection")

	m, err := migrate.New("file:///migrations", POSTGRES_URL)

	if err != nil {
		log.Fatalf("Could not create migrate instance: %v\n", err)
	}

	log.Println("Running migrations")

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run up migrations: %v\n", err)
	}

	log.Println("Migrations succesful")

	fs, err := services.NewFileservice()
	fc := controllers.NewFileController(fs)

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
				log.Println("Handling file upload")
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

	log.Println("Starting a server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
