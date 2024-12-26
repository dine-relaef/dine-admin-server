package env

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var AuthVar = map[string]string{
	"ACCESS_TOKEN_SECRET":  GetEnv("ACCESS_TOKEN_SECRET"),
	"ACCESS_TOKEN_AGE":     GetEnv("ACCESS_TOKEN_AGE"),
	"REFRESH_TOKEN_SECRET": GetEnv("REFRESH_TOKEN_SECRET"),
	"REFRESH_TOKEN_AGE":    GetEnv("REFRESH_TOKEN_AGE"),
	"GOOGLE_CLIENT_ID":     GetEnv("GOOGLE_CLIENT_ID"),
	"GOOGLE_CLIENT_SECRET": GetEnv("GOOGLE_CLIENT_SECRET"),
}

var (
	Config *oauth2.Config
)

func init() {
	var redirectURL string
	if AppVar["ENVIRONMENT"] != "development" {
		redirectURL = "https://" + AppVar["CLIENT_HOST"] + "/api/v1/auth/google/callback"
	} else {
		redirectURL = "http://localhost:8080/api/v1/auth/google/callback"
	}

	Config = &oauth2.Config{
		ClientID:     AuthVar["GOOGLE_CLIENT_ID"],
		ClientSecret: AuthVar["GOOGLE_CLIENT_SECRET"],
		RedirectURL:  redirectURL,
		// Update with your frontend callback URL
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
