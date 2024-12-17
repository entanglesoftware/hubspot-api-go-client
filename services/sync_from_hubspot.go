package services

import (
	"fmt"
	"hubspot-api-go-client/config"
	"hubspot-api-go-client/hubspot-api-go"
	"hubspot-api-go-client/models"
	"log"
)

func SyncFromHubsport() {

	mongoURI := config.GetMongoURI()
	InitMongoDB(mongoURI, "hubspot")

	var contacts []models.Contact

	// Initialize HubSpot client
	apiKey := config.GetHubSpotAPIKey()
	client := hubspot.NewClient(apiKey)

	contacts, err := GetContact(client)
	if err != nil {
		log.Fatalf("Error fetching contacts: %v", err)
	}

	// Sync contacts to MongoDB
	log.Printf("Found %d contacts. Syncing to MongoDB...\n", len(contacts))
	for _, contact := range contacts {
		err := UpsertContact(contact)
		if err != nil {
			log.Printf("contact %s", contact.ID)
		} else {
			fmt.Printf("Synced contact: %s (%s %s)\n", contact.Email, contact.FirstName, contact.LastName)
		}
	}

	log.Println("Sync completed successfully!")
}
