package store

import (
	"database/sql"

	"github.com/0LuigiCode0/Library/logger"

	_ "github.com/lib/pq" // ..
)

// Store ...
type Store struct {
	log    *logger.Logger
	config *DBConfig
	db     *sql.DB
}

// NewStore ...
func NewStore(config *DBConfig) *Store {
	return &Store{
		log:    logger.InitLogger(""),
		config: config,
	}
}

// Open ...
func (store *Store) Open() error {
	db, err := sql.Open("postgres", store.config.DatabaseURL)
	store.log.Service("open posgres")
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	store.log.Info("ping posgresql")
	store.db = db
	return nil
}

// Close ...
func (store *Store) Close() {
	store.db.Close()
	store.log.Service("close posgres")
}
