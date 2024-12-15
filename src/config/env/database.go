package env

import "os"

var PostgresDatabaseVar = map[string]string{
	"DATABASE_URL": os.Getenv("DATABASE_URL"),
}
