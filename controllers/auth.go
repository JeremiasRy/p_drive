package controllers

import (
	"backend/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	GITHUB_CLIENT_ID     = os.Getenv("GITHUB_CLIENT_ID")
	GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
	BACKEND_BASE_URL     = os.Getenv("BACKEND_BASE_URL")
	REDIRECT_URL         = fmt.Sprintf("%s/login/github/callback", BACKEND_BASE_URL)
	OAUTH_CONFIG         = &oauth2.Config{ClientID: GITHUB_CLIENT_ID, ClientSecret: GITHUB_CLIENT_SECRET, RedirectURL: REDIRECT_URL, Scopes: []string{"user:email"}, Endpoint: github.Endpoint}
)

type AuthController struct {
	us *services.UserService
}

func NewAuthController(us *services.UserService) *AuthController {
	return &AuthController{us: us}
}

func (ac *AuthController) HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "supersecret" {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")

	token, err := OAUTH_CONFIG.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Could not get token: %s\n%s\n", err, code)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	client := OAUTH_CONFIG.Client(r.Context(), token)
	userResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		log.Printf("Could not create request: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var userEmail []struct{ email string }
	err = json.NewDecoder(userResp.Body).Decode(&userEmail)

	if err != nil {
		log.Printf("Failed to parse response body %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ac.us.NewUser(userEmail[0].email)

	if err != nil {
		log.Printf("Failed to create new user %v, %v\n", userEmail, err)
		http.Error(w, fmt.Sprintf("Failed to authenticate: %s", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Success?")
}

func (ac *AuthController) HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	state := "supersecret"
	url := OAUTH_CONFIG.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
