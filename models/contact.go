package models

import "time"

type Contact struct {
	ID        string    `bson:"hs_contact_id"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"hs_updated_at"`
}
