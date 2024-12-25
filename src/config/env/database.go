package env

var PostgresDatabaseVar = map[string]string{
	"DATABASE_URL": GetEnv("DATABASE_URL"),
}
