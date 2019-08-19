package models

type UserStore interface {
	NewUser(firstname, lastname, email, password string) error
}
