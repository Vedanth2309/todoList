package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct{ Port, MongoURI, MongoDB, JWTSecret, ESURL, ESKey string }

var C Config

// Load tries backend/.env then ../.env. Each file is loaded independently
// (the old godotenv.Load("../.env", ".env") call stopped at the first missing file).
func Load() error {
	for _, f := range []string{".env", "../.env"} {
		if godotenv.Load(f) == nil {
			log.Println("Loaded", f)
		}
	}
	get := func(k string) string { return strings.Trim(strings.TrimSpace(os.Getenv(k)), `"'`) }
	C = Config{Port: get("PORT"), MongoURI: get("MONGODB_URI"), MongoDB: get("MONGODB_DATABASE"), JWTSecret: get("JWT_SECRET"),
		ESURL: strings.TrimRight(get("ELASTICSEARCH_URL"), "/"), ESKey: get("ELASTICSEARCH_API_KEY")}
	if C.Port == "" {
		C.Port = "8082"
	}
	if C.MongoDB == "" {
		C.MongoDB = "personal_tracker"
	}
	switch {
	case C.MongoURI == "":
		return errors.New("MONGODB_URI is empty or .env was not found (looked in ./.env and ../.env)")
	case strings.ContainsAny(C.MongoURI, "<>"):
		return errors.New("MONGODB_URI still contains <username>/<password> placeholders; remove the angle brackets")
	case C.JWTSecret == "":
		return errors.New("JWT_SECRET is empty")
	}
	return nil
}
