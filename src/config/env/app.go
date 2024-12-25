package env

var AppVar = map[string]string{
	"SERVER_HOST":          GetEnv("SERVER_HOST"),
	"APP_NAME":             GetEnv("APP_NAME"),
	"PORT":                 GetEnv("PORT"),
	"ENVIRONMENT":          GetEnv("ENVIRONMENT"),
	"ACCESS_TOKEN_SECRET":  GetEnv("ACCESS_TOKEN_SECRET"),
	"ACCESS_TOKEN_AGE":     GetEnv("ACCESS_TOKEN_AGE"),
	"REFRESH_TOKEN_SECRET": GetEnv("REFRESH_TOKEN_SECRET"),
	"REFRESH_TOKEN_AGE":    GetEnv("REFRESH_TOKEN_AGE"),
}
