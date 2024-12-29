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

type LayoutData struct {
}

func NewViewsController() *ViewsController {
	return &ViewsController{}
}

func (vc *ViewsController) HandleGetUpload(w http.ResponseWriter, r *http.Request) {
	upload, err := vc.readHTMLFile(filepath.Join("views", "upload.html"))
	if err != nil {
		log.Fatalf("Failed to load HTML content %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(upload)
}

func (vc *ViewsController) HandleGetRoot(w http.ResponseWriter, r *http.Request) {
	data := LayoutData{}
	vc.renderTemplate(w, data)
}

func (vc *ViewsController) HandleGetMain(w http.ResponseWriter, r *http.Request, u *model.Users) {
	home, err := vc.readHTMLFile(filepath.Join("views", "files.html"))
	if err != nil {
		log.Fatalf("Failed to load HTML content %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(home)
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

func (vc *ViewsController) renderTemplate(w http.ResponseWriter, data LayoutData) {
	layoutPath := filepath.Join("views", "layout.html")

	template, err := template.ParseFiles(layoutPath)

	if err != nil {
		log.Fatalf("Failed to load HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	template.Execute(w, data)
}

func (vc *ViewsController) readHTMLFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
