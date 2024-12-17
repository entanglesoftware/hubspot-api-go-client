package services

import (
	"context"
	"log"

	"hubspot-api-go-client/config"
	"hubspot-api-go-client/models"

	"github.com/entanglesoftware/hubspot-api-go/codegen/crm/objects/contacts"
	"github.com/entanglesoftware/hubspot-api-go/hubspot"
)

func GetContact(client *hubspot.Client) ([]models.Contact, error) {

	apiKey := config.GetHubSpotAPIKey()
	var contactsData []models.Contact

	// Initialize the client
	client.SetAccessToken(apiKey)

	limit := 10

	// Make the API call
	contactsParams := contacts.GetContactsParams{
		Limit: &limit,
	}

	ct := client.Crm().Contacts().Contacts

	response, err := ct.GetContactsWithResponse(context.Background(), &contactsParams)
	if err != nil {
		log.Fatalf("API call failed: %v", err)
	}

	if response.StatusCode() == 200 {
		if response.JSON200 == nil || response.JSON200.Results == nil {
			log.Fatalf("Response contains no results")
		}

		for _, result := range *response.JSON200.Results {

			// return contactsData, nil
			if *result.Properties != nil {
				properties := *result.Properties
				contact := models.Contact{
					ID:        *result.Id,
					FirstName: properties["firstname"],
					LastName:  properties["lastname"],
					Email:     properties["email"],
					CreatedAt: *result.CreatedAt,
					UpdatedAt: *result.UpdatedAt,
				}
				contactsData = append(contactsData, contact)
			} else {
				return contactsData, nil
			}
		}

		return contactsData, nil

	} else {
		log.Fatalf("API call Failed with status code %d: %v", response.StatusCode(), response)
	}

	return contactsData, nil
}
