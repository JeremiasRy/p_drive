package config

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Config struct {
	GITHUB_CLIENT_ID         string
	GITHUB_CLIENT_SECRET     string
	BACKEND_BASE_URL         string
	REDIRECT_URL             string
	OAUTH_CONFIG             *oauth2.Config
	POSTGRES_URL             string
	POSTGRES_UESR            string
	POSTGRES_PASSWORD        string
	SESSION_NAME             string
	SESSION_KEY              string
	SESSION_CURRENT_USER_KEY string
	SESSION_VERIFIER_KEY     string
}

var instance *Config

func initConfig() {
	GITHUB_CLIENT_ID := os.Getenv("GITHUB_CLIENT_ID")
	GITHUB_CLIENT_SECRET := os.Getenv("GITHUB_CLIENT_SECRET")
	BACKEND_BASE_URL := os.Getenv("BACKEND_BASE_URL")
	REDIRECT_URL := fmt.Sprintf("%s/login/github/callback", BACKEND_BASE_URL)
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	OAUTH_CONFIG := &oauth2.Config{ClientID: GITHUB_CLIENT_ID, ClientSecret: GITHUB_CLIENT_SECRET, RedirectURL: REDIRECT_URL, Scopes: []string{"user:email"}, Endpoint: github.Endpoint}
	POSTGRES_URL := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_USER)
	SESSION_NAME := "user_session"
	SESSION_CURRENT_USER_KEY := "current_user"
	SESSION_VERIFIER_KEY := "verifier_challenge"
	SESSION_KEY := os.Getenv("SESSION_KEY")

	instance = &Config{
		GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, BACKEND_BASE_URL, REDIRECT_URL, OAUTH_CONFIG, POSTGRES_URL, POSTGRES_USER, POSTGRES_PASSWORD, SESSION_NAME, SESSION_KEY, SESSION_CURRENT_USER_KEY, SESSION_VERIFIER_KEY,
	}
}

func GetConfig() *Config {
	if instance == nil {
		initConfig()
	}
	return instance
}
