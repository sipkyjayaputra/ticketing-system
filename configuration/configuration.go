package configuration

import (
	"os"
)

const (
	LOCAL       = "local"
	DEVELOPMENT = "development"
	PRODUCTION  = "production"
)

// ENVIRONMENT:
const ENVIRONMENT string = DEVELOPMENT // LOCAL, DEVELOPMENT, PRODUCTION

var env = map[string]map[string]string{
	"local": {
		"PORT":          "8080",
		"MYSQL_HOST":    "localhost",
		"MYSQL_PORT":    "3306",
		"MYSQL_USER":    "root",
		"MYSQL_PASS":    "root",
		"MYSQL_DB_NAME": "ticketing_system",
		"POSTGRES_HOST":  "localhost",
		"POSTGRES_PORT":  "5432",
		"POSTGRES_USER":  "postgres",
		"POSTGRES_PASS":  "admin",
		"POSTGRES_DB_NAME": "ticketing_system",
	},
	"development": {
		"PORT":          "8080",
		"MYSQL_HOST":    "mysql",
		"MYSQL_PORT":    "3306",
		"MYSQL_USER":    "root",
		"MYSQL_PASS":    "root",
		"MYSQL_DB_NAME": "ticketing_system",
		"POSTGRES_HOST":  "postgres", 
		"POSTGRES_PORT":  "5432",
		"POSTGRES_USER":  "postgres",
		"POSTGRES_PASS":  "admin",
		"POSTGRES_DB_NAME": "ticketing_system",
	},
}

// CONFIG : global configuration
var CONFIG = env[ENVIRONMENT]

// Getenv : function for Environment Lookup
func Getenv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitConfig() {
	for key := range CONFIG {
		CONFIG[key] = Getenv(key, CONFIG[key])
		os.Setenv(key, CONFIG[key])
	}
}
