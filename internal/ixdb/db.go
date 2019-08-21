package ixdb

import (
	"fmt"

	"github.com/rdoorn/ixxi/internal/ixdb/ixmemory"
)

type DB struct {
	Provider         string `mapstructure:"provider"`
	ConnectionString string `mapstructure:"connection_string"`
	ApiKey           string `mapstructure:"api_key"`
	ApiSecret        string `mapstructure:"api_secret"`
}

func (d *DB) Verify() error {

	switch d.Provider {
	case "memory":
		return nil
	case "mongodb":
		if d.ConnectionString == "" {
			return fmt.Errorf("missing connection string for provider: %s", d.Provider)
		}
		return nil
	default:
		return fmt.Errorf("unknown provider: %s", d.Provider)
	}
}

func (d *DB) Setup() (interface{}, error) {
	switch d.Provider {
	case "memory":
		provider := ixmemory.New()
		return provider, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", d.Provider)
	}
}
