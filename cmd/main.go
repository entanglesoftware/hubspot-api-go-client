package main

import (
	"hubspot-api-go-client/config"
	"hubspot-api-go-client/hubspot-api-go"
	"log"
)

func main() {
	// Load API key
	apiKey := config.GetHubSpotAPIKey()

	// Initialize HubSpot client
	client := hubspot.NewClient(apiKey)

	// Example: Create a new contact
	err := hubspot.CreateContact(client, "johndoeexamle1@example.com", "John", "Doe")
	if err != nil {
		log.Fatalf("Error creating contact: %v", err)
	}
}
