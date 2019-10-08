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

// Get executes a get of 1 object, and returns one result
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

// UpdateOne updates a single item in a record on mongo
func (m *MongoDB) UpdateOne(filter, update bson.M) error {
	collection := m.database.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// Add writes a user record to the database
func (m *MongoDB) Add(i interface{}) (string, error) {
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
	res, err := collection.InsertOne(ctx, i)
	if err != nil {
		return "", err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.String(), nil
	}
	return "", err
}
