package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ViewsController struct{}

func NewViewsController() *ViewsController {
	return &ViewsController{}
}

func (vc *ViewsController) HandleGetRoot(w http.ResponseWriter, r *http.Request) {
	layout, err := vc.readHTMLFile(filepath.Join("views", "layout.html"))

	if err != nil {
		log.Printf("Failed to load HTML content %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(layout)
}

func (vc *ViewsController) HandleGetMain(w http.ResponseWriter, r *http.Request, u *model.Users) {
	homePath := filepath.Join("views", "home.html")

	template, err := template.ParseFiles(homePath)

	if err != nil {
		log.Fatalf("Failed to load HTML content %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = template.Execute(w, struct {
		UserID string
	}{UserID: u.ID.String()})

	if err != nil {
		log.Printf("Failed to render template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (vc *ViewsController) HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	login, err := vc.readHTMLFile(filepath.Join("views", "login.html"))
	if err != nil {
		log.Fatalf("Failed to load HTML content %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(login)
}

func (vc *ViewsController) renderTemplate(w http.ResponseWriter) {
	layoutPath := filepath.Join("views", "layout.html")

	template, err := template.ParseFiles(layoutPath)

	if err != nil {
		log.Fatalf("Failed to load HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	template.Execute(w, struct{}{})
}

func (vc *ViewsController) readHTMLFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
