package ixmemory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	m := New()

	err := m.CreateUser("username", "firstname", "lastname", "email", "passwordSalt", "passwordHash", "activationToken")
	assert.Nil(t, err)
}
