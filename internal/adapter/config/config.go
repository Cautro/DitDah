package config

import (
	"os"
)

type Config struct {
	Port string
	DBPath string
	JwtSecret string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		panic("db path don't exist")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT secret is not exist in .env")
	}

	return Config{
		Port: port,
		DBPath: dbPath,
		JwtSecret: jwtSecret,
	}, nil
}