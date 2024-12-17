package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"hubspot-api-go-client/config"
	"hubspot-api-go-client/models"

	"github.com/entanglesoftware/hubspot-api-go/codegen/crm/objects/contacts"
	"github.com/entanglesoftware/hubspot-api-go/hubspot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDBclient *mongo.Client
var contactCollection *mongo.Collection

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

func InitMongoDB(uri, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoDBclient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	contactCollection = mongoDBclient.Database(dbName).Collection("contacts")
	log.Println("Connected to MongoDB!")
}

func UpsertContact(contact models.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"contact_id": contact.ID}
	update := bson.M{"$set": contact}
	options := options.Update().SetUpsert(true)

	_, err := contactCollection.UpdateOne(ctx, filter, update, options)
	return err
}
