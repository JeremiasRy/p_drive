package controllers

import (
	"backend/services"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type ViewsController struct {
	fs *services.FoldersService
}

func NewViewsController(fs *services.FoldersService) *ViewsController {
	return &ViewsController{fs: fs}
}

func (vc *ViewsController) HandleGetRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/folders/", http.StatusPermanentRedirect)
}

func (vc *ViewsController) HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(filepath.Join("views", "templates", "layout.html"), filepath.Join("views", "templates", "login", "login.html"))
	if err != nil {
		log.Printf("Failed to parse HTML template %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = template.ExecuteTemplate(w, "layout", struct{}{})

	if err != nil {
		log.Printf("Failed to execute template %v\n", err)
	}
}
