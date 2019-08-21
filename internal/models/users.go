package models

type UserInterface interface {
	CreateUser(firstname, lastname, email, password string) error
	GetByEmail(email string) (*User, error)
}

type User struct {
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Emails    []Email `json:"emails"`
	Password  string  `json:"password"`
	Active    bool    `json:"active"`
}

type Email struct {
	Primary  bool   `json:"primary"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}
