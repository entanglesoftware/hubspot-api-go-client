package hubspot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"hubspot-api-go-client/config"

	"github.com/entanglesoftware/hubspot-api-go/codegen/crm/objects/contacts"
	"github.com/entanglesoftware/hubspot-api-go/hubspot"
)

func CreateContact(client *hubspot.Client, email, firstName, lastName string) error {

	apiKey := config.GetHubSpotAPIKey()

	// Initialize the client
	client.SetAccessToken(apiKey)

	// Initialize a variable of type Contact
	contact := contacts.CreateContactJSONBody{
		Properties: map[string]string{
			"firstname": firstName,
			"lastname":  lastName,
			"email":     email,
		},
	}

	// Serialize the contact properties to JSON
	body, err := json.Marshal(contact)
	if err != nil {
		log.Fatalf("Error serializing contact properties: %v", err)
	}

	contentType := "application/json"

	ct := client.Crm().Contacts().Contacts

	response, err := ct.CreateContactWithBodyWithResponse(context.Background(), contentType, bytes.NewReader(body))
	if err != nil {
		log.Fatalf("API call failed: %v", err)
		return err
	}

	fmt.Printf("Contact created successfully: %v\n", response.JSON201.Id)
	return nil
}
