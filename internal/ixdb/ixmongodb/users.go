package ixmongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/rdoorn/ixxi/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func (m *MongoDB) GetByEmail(email string) (*models.User, error) {
	//return m.Get(bson.M{"email": email})
	return m.Get(bson.M{"emails": bson.M{"$elemMatch": bson.M{"email": email}}})
}

func (m *MongoDB) GetByUsername(username string) (*models.User, error) {
	return m.Get(bson.M{"username": username})
}

func (m *MongoDB) GetByID(id string) (*models.User, error) {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.Get(bson.M{"_id": bsonID})
}

func (m *MongoDB) Get(filter bson.M) (*models.User, error) {
	collection := m.database.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := models.User{}

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("not found: %s: %s", filter, err)
	}
	return &user, nil
}

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

func (m *MongoDB) Add(u models.User) (string, error) {
	collection := m.database.Collection("users")
	log.Printf("Collection of users: %v", collection)
	/*	if collection == nil {
		var err error
		collection, err = m.CreateCollection("users")
		if err != nil {
			return "", err
		}
	}*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, &u)
	if err != nil {
		return "", err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.String(), nil
	}
	return "", err
}

func (m *MongoDB) SetStorageLimit(id string, limitByte int64) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": bsonID}
	update := bson.M{"$set": bson.M{"storage_max_byte": limitByte}}

	return m.UpdateOne(filter, update)
}

func (m *MongoDB) UpdateOne(filter, update bson.M) error {
	collection := m.database.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
