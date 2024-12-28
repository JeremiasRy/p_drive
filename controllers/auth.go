package controllers

import (
	"backend/config"
	"backend/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type AuthController struct {
	us *services.UserService
	ss *services.SessionService
}

func NewAuthController(us *services.UserService, ss *services.SessionService) *AuthController {
	return &AuthController{us: us, ss: ss}
}

var c = config.GetConfig()

func (ac *AuthController) HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	session := ac.ss.GetSession(w, r)

	if state != session.Values["verifier"] {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")

	token, err := c.OAUTH_CONFIG.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Could not get token: %s\n%s\n", err, code)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	client := c.OAUTH_CONFIG.Client(r.Context(), token)
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

	err = ac.us.NewUser(userEmail[0].email, services.GITHUB)

	if err != nil {
		log.Printf("Failed to create new user %v, %v\n", userEmail, err)
		http.Error(w, fmt.Sprintf("Failed to authenticate: %s", err), http.StatusInternalServerError)
		return
	}

	user, err := ac.us.GetUserByEmail(userEmail[0].email)

	if err != nil {
		log.Printf("Failed to fetch user after creation %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.Values["current_user"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, r.URL.Host, http.StatusPermanentRedirect)
}

func (ac *AuthController) HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	state := oauth2.GenerateVerifier()
	session := ac.ss.GetSession(w, r)

	session.Values["verifier"] = state
	session.Save(r, w)

	url := c.OAUTH_CONFIG.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
