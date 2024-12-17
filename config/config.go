package config

import (
	"log"
	"os"
)

func GetHubSpotAPIKey() string {
	apiKey := os.Getenv("HS_ACCESS_TOKEN")
	if apiKey == "" {
		log.Fatal("HS_ACCESS_TOKEN is not set in environment variables")
	}
	return apiKey
}

func GetHubSpotBaseURL() string {
	baseUrl := os.Getenv("HS_BASE_URL")
	if baseUrl == "" {
		return "https://api.hubapi.com"
	}
	return baseUrl
}

func GetMongoURI() string {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set in environment variables")
	}
	return mongoURI
}
