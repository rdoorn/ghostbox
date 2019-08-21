package ixmemory

import (
	"fmt"
	"sync"

	"github.com/rdoorn/ixxi/internal/models"
)

type Memory struct {
	users []models.User
	m     sync.RWMutex
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) CreateUser(firstname, lastname, email, password string) error {
	emailModel := models.Email{
		Primary:  true,
		Email:    email,
		Verified: false,
	}
	user := models.User{
		Firstname: firstname,
		Lastname:  lastname,
		Emails:    []models.Email{emailModel},
		Active:    false,
	}
	m.Add(user)
	return nil
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

func (m *Memory) Add(u models.User) {
	m.m.Lock()
	defer m.m.Unlock()
	m.users = append(m.users, u)
}
