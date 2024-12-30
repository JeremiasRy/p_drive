package middleware

import (
	"backend/.gen/personal_drive/public/model"
	"backend/config"
	"backend/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *model.Users)

type EnsureAuth struct {
	handler AuthenticatedHandler
	us      *services.UserService
	store   *sessions.CookieStore
}

var c = config.GetConfig()

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := ea.GetAuthenticatedUser(r)

	if err != nil {
		log.Printf("Failed to authenticate session %v", err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	ea.handler(w, r, user)
}

func (ea *EnsureAuth) GetAuthenticatedUser(r *http.Request) (*model.Users, error) {
	session, err := ea.store.Get(r, c.SESSION_NAME)

	if err != nil {
		return nil, err
	}

	if session.IsNew {
		return nil, fmt.Errorf("session was outdated")
	}

	userId, found := session.Values["current_user"].(string)

	if !found {
		return nil, fmt.Errorf("no current_user in session.Values")
	}

	user, err := ea.us.GetUserByID(userId)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewEnsureAuth(us *services.UserService, store *sessions.CookieStore, handlerToWrap AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{us: us, store: store, handler: handlerToWrap}
}
