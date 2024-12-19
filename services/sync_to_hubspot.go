package services

import (
	"bytes"
	"context"
	"encoding/json"
	"hubspot-api-go-client/config"
	"hubspot-api-go-client/database"
	"hubspot-api-go-client/hubspot-api-go"
	"log"

	"github.com/entanglesoftware/hubspot-api-go/codegen/crm/objects/contacts"
)

func SyncToHubsport() {

	mongoURI := config.GetMongoURI()
	// Connect to MongoDB
	mongoClient := database.ConnectMongo(mongoURI)
	defer func() {
		if err := mongoClient.Disconnect(nil); err != nil {
			log.Fatalf("Error disconnecting MongoDB: %v", err)
		}
	}()
	dbName := "hubspot"
	collectionName := "contacts"

	mongoDbcontacts, mongoDbErr := database.GetAllContacts(mongoClient, dbName, collectionName)
	if mongoDbErr != nil {
		log.Fatalf("Enable to fetch contact from mongoDB: %v", mongoDbErr)
	}

	// Initialize HubSpot client
	apiKey := config.GetHubSpotAPIKey()
	client := hubspot.NewClient(apiKey)
	client.SetAccessToken(apiKey)
	contentType := "application/json"

	ct := client.Crm().Contacts().Contacts

	for _, contact := range mongoDbcontacts {
		// Initialize a variable of type Contact
		firstName := contact["first_name"].(string)
		lastName := contact["last_name"].(string)
		email := contact["email"].(string)
		hscontact := contacts.CreateContactJSONBody{
			Properties: map[string]string{
				"firstname": firstName,
				"lastname":  lastName,
				"email":     email,
			},
		}
		log.Printf("Contact First Name: %v\n", firstName)
		body, _ := json.Marshal(hscontact)

		contactId, ok := contact["hs_contact_id"].(string)
		if !ok || contactId == "" {
			// Create a new contact in HubSpot
			response, err := ct.CreateContactWithBodyWithResponse(context.Background(), contentType, bytes.NewReader(body))
			if err != nil {
				log.Printf("Error creating contact in HubSpot: %v", err)
				continue
			}
			if response.StatusCode() == 201 && response.JSON201 != nil {
				newContactID := *response.JSON201.Id
				log.Printf("Created new contact with ID: %s", newContactID)

				// Update MongoDB with the new hs_contact_id
				updateErr := database.UpdateContactHsID(mongoClient, dbName, collectionName, contact["_id"], newContactID)
				if updateErr != nil {
					log.Printf("Error updating MongoDB with hs_contact_id: %v", updateErr)
				}
			} else {
				log.Printf("Failed to create contact: %v", response)
			}
		} else {
			response, err := ct.UpdateContactWithBodyWithResponse(context.Background(), contactId, contentType, bytes.NewReader(body))
			if err != nil {
				log.Printf("Error to update contact %s", contactId)
			}

			if response.StatusCode() == 200 {
				if response.JSON200 == nil || response.JSON200.Id == nil {
					log.Printf("Response contains no results")
				}

				if response.JSON200.Properties != nil {
					log.Printf("Properties: %s", response.JSON200.Properties)
				} else {
					log.Printf("No properties found.")
				}
			} else {
				log.Printf("Test Failed with status code %d: %v", response.StatusCode(), response)
			}
		}
	}
}
