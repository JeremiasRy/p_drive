package controllers

import (
	"backend/.gen/personal_drive/public/model"
	"backend/services"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ViewsController struct {
	fs *services.FoldersService
}

func NewViewsController(fs *services.FoldersService) *ViewsController {
	return &ViewsController{fs: fs}
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

func (vc *ViewsController) HandleGetHome(w http.ResponseWriter, r *http.Request, u *model.Users) {
	template, err := template.ParseFiles(filepath.Join("views", "templates", "home.html"), filepath.Join("views", "templates", "upload.html"))

	if err != nil {
		log.Printf("Failed to load HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(w, "home", struct{}{})

	if err != nil {
		log.Printf("Failed to execute HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (vc *ViewsController) HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	login, err := vc.readHTMLFile(filepath.Join("views", "login.html"))
	if err != nil {
		log.Printf("Failed to load HTML content %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(login)
}

func (vc *ViewsController) readHTMLFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
