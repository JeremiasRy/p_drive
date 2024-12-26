package controllers

import (
	"context"
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
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "9fg9d8fb9d8fb9dfb89d8fb" {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")

	token, err := OAUTH_CONFIG.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get token: %s", err), http.StatusBadRequest)
		return
	}

	client := OAUTH_CONFIG.Client(context.Background(), token)
	userResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create request: %s", err), http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()
	log.Printf("%v\n", *userResp)
}

func (ac *AuthController) HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	state := "9fg9d8fb9d8fb9dfb89d8fb"
	url := OAUTH_CONFIG.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
