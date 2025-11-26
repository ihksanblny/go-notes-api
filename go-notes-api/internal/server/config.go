package server

import (
	"log"
	"os"
)

type Config struct {
	Port           string
	AllowedOrigins string
}

func LoadConfig() Config {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	origin := os.Getenv("API_ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}
	log.Printf("Using Config: PORT=%s, ORIGIN=%s", port, origin)

	return Config{
		Port: port,
		AllowedOrigins: origin,
	}
}