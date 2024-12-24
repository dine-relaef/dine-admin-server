package env

import "os"

var AppVar = map[string]string{
	"APP_NAME":             os.Getenv("APP_NAME"),
	"PORT":                 os.Getenv("PORT"),
	"ENVIRONMENT":          os.Getenv("ENVIRONMENT"),
	"ACCESS_TOKEN_SECRET":  os.Getenv("ACCESS_TOKEN_SECRET"),
	"ACCESS_TOKEN_AGE":     os.Getenv("ACCESS_TOKEN_AGE"),
	"REFRESH_TOKEN_SECRET": os.Getenv("REFRESH_TOKEN_SECRET"),
	"REFRESH_TOKEN_AGE":    os.Getenv("REFRESH_TOKEN_AGE"),
}
