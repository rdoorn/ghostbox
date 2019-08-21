package models

type UserInterface interface {
	CreateUser(username, firstname, lastname, email, passwordSalt, passwordHash, activationToken string) error
	ActivateUser(id string) error
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByID(id string) (*User, error)
}

type User struct {
	ID              string   `json:"id"`               // user id generated
	Username        string   `json:"username"`         // username used to display
	Firstname       string   `json:"first_name"`       // first name of client
	Lastname        string   `json:"last_name"`        // last name of client
	Emails          []Email  `json:"emails"`           // emails registered to client
	PasswordSalt    string   `json:"password_salt"`    // salt of password
	PasswordHash    string   `json:"password_hash"`    // hash of password
	Active          bool     `json:"active"`           // is the account activated using the token
	ActivationToken string   `json:"activation_token"` // token to activate the accounts
	Tokens          []string `json:"tokens"`           // access tokens
}

type Email struct {
	Primary  bool   `json:"primary"`  // only primary email is marked
	Email    string `json:"email"`    // email added to account
	Verified bool   `json:"verified"` // is email address verified
}
