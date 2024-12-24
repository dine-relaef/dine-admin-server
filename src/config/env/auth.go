package env

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var AuthVar = map[string]string{
	"ACCESS_TOKEN_SECRET":  os.Getenv("ACCESS_TOKEN_SECRET"),
	"ACCESS_TOKEN_AGE":     os.Getenv("ACCESS_TOKEN_AGE"),
	"REFRESH_TOKEN_SECRET": os.Getenv("REFRESH_TOKEN_SECRET"),
	"REFRESH_TOKEN_AGE":    os.Getenv("REFRESH_TOKEN_AGE"),
	"GOOGLE_CLIENT_ID":     os.Getenv("GOOGLE_CLIENT_ID"),
	"GOOGLE_CLIENT_SECRET": os.Getenv("GOOGLE_CLIENT_SECRET"),
}

var (
	Config *oauth2.Config
)

func init() {
	Config = &oauth2.Config{
		ClientID:     AuthVar["GOOGLE_CLIENT_ID"],
		ClientSecret: AuthVar["GOOGLE_CLIENT_SECRET"],
		RedirectURL:  "http://localhost:8080/api/v1/auth/google/callback",
		// Update with your frontend callback URL
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
