package services

import (
	"backend/config"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionService struct {
	store *sessions.CookieStore
}

var c = config.GetConfig()
var store = sessions.NewCookieStore([]byte(c.SESSION_KEY))

func NewSessionService() *SessionService {
	return &SessionService{store: store}
}

func (ss *SessionService) GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := ss.store.Get(r, c.SESSION_NAME)

	if err != nil {
		log.Printf("Failed to fetch session %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil
	}

	return session
}
