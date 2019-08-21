package ixmemory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rdoorn/ixxi/internal/models"
)

type Memory struct {
	users []models.User
	m     sync.RWMutex
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) CreateUser(username, firstname, lastname, email, passwordSalt, passwordHash, activationToken string) error {
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
	m.Add(user)
	return nil
}

func (m *Memory) ActivateUser(id string) error {
	m.m.Lock()
	defer m.m.Unlock()
	for i, u := range m.users {
		if u.ID == id {
			m.users[i].ActivationToken = ""
			m.users[i].Active = true
			m.users[i].Emails[0].Verified = true
			return nil
		}
	}
	return fmt.Errorf("ID not found: %s", id)
}

func (m *Memory) GetByEmail(email string) (*models.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	for _, u := range m.users {
		for _, e := range u.Emails {
			if e.Email == email {
				return &u, nil
			}
		}
	}
	return nil, fmt.Errorf("Email not found: %s", email)
}

func (m *Memory) GetByUsername(username string) (*models.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	for _, u := range m.users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("Username not found: %s", username)
}

func (m *Memory) GetByID(id string) (*models.User, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("ID not found: %s", id)
}

func (m *Memory) Add(u models.User) {
	m.m.Lock()
	defer m.m.Unlock()
	u.ID = uuid.New().String()
	m.users = append(m.users, u)
}
