package database

import (
	"context"
	"hubspot-api-go-client/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Contact struct {
	Email     string `bson:"email"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
}

// SaveContact saves a contact to MongoDB.
func SaveContact(client *mongo.Client, dbName, collectionName string, contact models.Contact) (any, error) {
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, contact)
	if err != nil {
		return nil, err
	}
	return result, err
}

// GetContact retrieves a contact by email from MongoDB.
func GetContact(client *mongo.Client, dbName, collectionName, email string) (*models.Contact, error) {
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var contact models.Contact
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&contact)
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetAllContacts retrieves all contacts from MongoDB.
func GetAllContacts(client *mongo.Client, dbName, collectionName string) ([]bson.M, error) {
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contacts []bson.M
	if err = cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}

	return contacts, nil
}

// UpdateContactHsID updates the hs_contact_id of a contact in MongoDB
func UpdateContactHsID(client *mongo.Client, dbName, collectionName string, contactID interface{}, newHsContactID string) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Create a filter to find the contact by its MongoDB `_id`
	filter := bson.M{"_id": contactID}

	// Define the update operation
	update := bson.M{
		"$set": bson.M{
			"hs_contact_id": newHsContactID,
		},
	}

	// Set update options
	opts := options.Update().SetUpsert(false)

	// Perform the update
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error updating hs_contact_id for contact %v: %v", contactID, err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Printf("No contact found with ID %v to update hs_contact_id.", contactID)
	} else {
		log.Printf("Updated hs_contact_id for contact ID %v.", contactID)
	}

	return nil
}
