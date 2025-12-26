package config

import (
	"log"
	"os"
)

type Config struct {
	Port      string
	OBAApiKey string
	OBAURL    string
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var %s", v)
	}

	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func Load() Config {
	port := getEnv("PORT", "8080")
	apiKey := mustEnv("ONEBUSAWAY_API_KEY")
	baseURL := getEnv("ONEBUSAWAY_BASE_URL", "https://api.pugetsound.onebusaway.org/api/where")

	return Config{
		Port:      port,
		OBAApiKey: apiKey,
		OBAURL:    baseURL,
	}

}
