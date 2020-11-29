package apiserver

import (
	"github.com/nergilz/tmpserver/internal/app/store"
)

// Config server struct
type Config struct {
	BindAddr string
	Store    *store.DBConfig
}

// NewConfig server config
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		Store:    store.NewDBConfig(),
	}
}
