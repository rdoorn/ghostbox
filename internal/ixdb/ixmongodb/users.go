package ixmongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/rdoorn/ixxi/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser adds a new user to the database with no rights, and not activated yet
func (m *MongoDB) CreateUser(username, firstname, lastname, email, passwordSalt, passwordHash, activationToken string) error {
	emailModel := models.Email{
		Primary:  true,
		Email:    email,
		Verified: false,
	}
	user := models.User{
		Username:        username,
		Firstname:       firstname,
		Lastname:        lastname,
		Emails:          []models.Email{emailModel},
		Active:          false,
		Tokens:          []string{},
		PasswordSalt:    passwordSalt,
		PasswordHash:    passwordHash,
		ActivationToken: activationToken,
	}

	_, err := m.Add(user)
	return err
}

// ActivateUser activates a user, based on a provided activation token
func (m *MongoDB) ActivateUser(id, token string) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": bsonID, "activation_token": token}
	update := bson.M{
		"$set":   bson.M{"active": true, "emails.0.verfied": true},
		"$unset": bson.M{"activation_token": ""},
	}

	collection := m.database.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("No match in activation filter:%s: update:%s", filter, update)
	}
	return nil
}

// GetByEmail returns a user model based on a provided email
func (m *MongoDB) GetByEmail(email string) (*models.User, error) {
	//return m.Get(bson.M{"email": email})
	return m.Get(bson.M{"emails": bson.M{"$elemMatch": bson.M{"email": email}}})
}

// GetByUsername returns a user model based on a provided username
func (m *MongoDB) GetByUsername(username string) (*models.User, error) {
	return m.Get(bson.M{"username": username})
}

// GetByID returns a user model based on a provided user id
func (m *MongoDB) GetByID(id string) (*models.User, error) {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.Get(bson.M{"_id": bsonID})
}

// AddToken adds a token or tokens to an existing user ID
func (m *MongoDB) AddToken(id string, tokens ...string) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": bsonID}

	//update := bson.M{"$addToSet": bson.M{"tokens": tokens}}
	update := bson.M{"$push": bson.M{"tokens": bson.M{"$each": tokens}}}

	collection := m.database.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("No match in update of token, or already added: %s - %s", filter, update)
	}
	return nil
}

// SetStorageLimit sets the storage limit object to a specific user
func (m *MongoDB) SetStorageLimit(id string, limitByte int64) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": bsonID}
	update := bson.M{"$set": bson.M{"storage_max_byte": limitByte}}

	return m.UpdateOne(filter, update)
}
