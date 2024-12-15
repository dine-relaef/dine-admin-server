package env

import "os"

var AppVar = map[string]string{
	"APP_NAME":    os.Getenv("APP_NAME"),
	"PORT":        os.Getenv("PORT"),
	"ENVIRONMENT": os.Getenv("ENVIRONMENT"),
}
