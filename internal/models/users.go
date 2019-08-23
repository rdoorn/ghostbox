package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInterface interface {
	CreateUser(username, firstname, lastname, email, passwordSalt, passwordHash, activationToken string) error
	ActivateUser(id, token string) error
	AddToken(id string, tokens ...string) error
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByID(id string) (*User, error)
	SetStorageLimit(id string, byte int64) error
}

type User struct {
	//BsonID          bson.ObjectId   `json:"_id" bson:"_id,omitempty"`                // user id generated
	ID              interface{} `json:"id" bson:"_id,omitempty"`                    // user id generated
	Username        string      `json:"username" bson:"username"`                   // username used to display
	Firstname       string      `json:"first_name" bson:"first_name"`               // first name of client
	Lastname        string      `json:"last_name" bson:"last_name"`                 // last name of client
	Emails          []Email     `json:"emails" bson:"emails"`                       // emails registered to client
	PasswordSalt    string      `json:"password_salt" bson:"password_salt"`         // salt of password
	PasswordHash    string      `json:"password_hash" bson:"password_hash"`         // hash of password
	Active          bool        `json:"active" bson:"active"`                       // is the account activated using the token
	ActivationToken string      `json:"activation_token" bson:"activation_token"`   // token to activate the accounts
	Tokens          []string    `json:"tokens" bson:"tokens"`                       // access tokens
	CollectionRoot  string      `json:"root" bson:"root"`                           // root directory of the collections
	StorageMaxByte  int64       `json:"storage_max_byte" bson:"storage_max_byte"`   // max storage allowed in Byte
	StorageUsedByte int64       `json:"storage_used_byte" bson:"storage_used_byte"` // actual storage used in Byte
}

type Email struct {
	Primary  bool   `json:"primary" bson:"primary"`   // only primary email is marked
	Email    string `json:"email" bson:"email"`       // email added to account
	Verified bool   `json:"verified" bson:"verified"` // is email address verified
}

func (u *User) Id() string {
	switch u.ID.(type) {
	case string:
		return u.ID.(string)
	case primitive.ObjectID:
		return u.ID.(primitive.ObjectID).Hex()
	default:
		return fmt.Sprintf("unknown ID type: %T", u.ID)
	}
}
