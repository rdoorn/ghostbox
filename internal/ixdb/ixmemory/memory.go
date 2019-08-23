package ixmemory

import (
	"sync"

	"github.com/rdoorn/ixxi/internal/models"
)

type Memory struct {
	users []models.User
	//db
	m sync.RWMutex
}

func New() *Memory {
	return &Memory{}
}
