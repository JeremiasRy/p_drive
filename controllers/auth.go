package controllers

import (
	"backend/config"
	"backend/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type AuthController struct {
	us    *services.UserService
	store *sessions.CookieStore
}

func NewAuthController(us *services.UserService, store *sessions.CookieStore) *AuthController {
	return &AuthController{us: us, store: store}
}

var c = config.GetConfig()

func (ac *AuthController) HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	session, err := ac.store.Get(r, c.SESSION_NAME)

	if err != nil {
		log.Printf("Failed to fetch sesssion %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if state != session.Values[c.SESSION_VERIFIER_KEY] {
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

	err = ac.us.NewUser(r.Context(), userEmail[0].email)
	if err != nil {
		log.Printf("Failed to create new user %v, %v\n", userEmail, err)
		http.Error(w, fmt.Sprintf("Failed to authenticate: %s", err), http.StatusInternalServerError)
		return
	}

	user, err := ac.us.GetUserByEmail(userEmail[0].email)

	if err != nil {
		log.Printf("Database failure when checking does user exist %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.Values[c.SESSION_CURRENT_USER_KEY] = user.ID.String()
	err = session.Save(r, w)

	if err != nil {
		log.Printf("Failed to save session %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/folders/my-drive", http.StatusFound)
}

func (ac *AuthController) HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	state := oauth2.GenerateVerifier()
	session, err := ac.store.Get(r, c.SESSION_NAME)

	if err != nil {
		log.Printf("Failed to fetch sesssion %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.Values[c.SESSION_VERIFIER_KEY] = state

	err = session.Save(r, w)

	if err != nil {
		log.Printf("Failed to save session %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	url := c.OAUTH_CONFIG.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
