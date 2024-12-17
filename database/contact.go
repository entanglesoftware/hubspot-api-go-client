package database

import (
	"context"
	"hubspot-api-go-client/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
func GetAllContacts(client *mongo.Client, dbName, collectionName string) ([]models.Contact, error) {
	collection := client.Database(dbName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Slice to store all contacts
	var contacts []models.Contact

	// Iterate through the cursor and decode each document
	for cursor.Next(ctx) {
		var contact models.Contact
		if err := cursor.Decode(&contact); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
