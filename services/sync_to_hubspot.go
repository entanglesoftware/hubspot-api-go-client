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
		body, err := json.Marshal(hscontact)
		contactId := contact["hs_contact_id"].(string)
		if err != nil {
			log.Fatalf("Error serializing contact properties: %v", err)
		}
		response, err := ct.UpdateContactWithBodyWithResponse(context.Background(), contactId, contentType, bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error to update contact %s", contactId)
		}

		if response.StatusCode() == 200 {
			if response.JSON200 == nil || response.JSON200.Id == nil {
				log.Fatalf("Response contains no results")
			}

			if response.JSON200.Properties != nil {
				log.Fatalf("Properties: %s", response.JSON200.Properties)
			} else {
				log.Fatalf("No properties found.")
			}
		} else {
			log.Printf("Test Failed with status code %d: %v", response.StatusCode(), response)
		}
	}
}
