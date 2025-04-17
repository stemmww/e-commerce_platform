package config

import (
	"os"
)

func GetDBConnectionString() string {
	return os.Getenv("POSTGRES_DSN") // e.g., postgres://user:pass@host:port/dbname
}
