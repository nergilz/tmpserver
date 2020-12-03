package database

import (
	"database/sql"

	"github.com/0LuigiCode0/Library/logger"

	_ "github.com/lib/pq" // ..
)

// Config ..
type Config struct {
	DatabaseURL string
	log         *logger.Logger
}

// DB ..
type DB struct {
	*Config
	cdb *sql.DB
}

// New DB config
func New(log *logger.Logger) *Config {
	return &Config{
		DatabaseURL: "host=localhost user=postgres dbname=restapi_tmp sslmode=disable",
		log:         log,
	}
}

// Connect to DB
func Connect(c *Config) error {
	conn, err := sql.Open("postgres", c.DatabaseURL)
	if err != nil {
		c.log.Errorf("error open DB : %v", err)
		return err
	}
	if err = conn.Ping(); err != nil {
		c.log.Errorf("error ping DB : %v", err)
		return err
	}
	c.log.Service("open DB & ping DB")
	return nil
}

// Conn return the active connection
func (db *DB) Conn() *sql.DB {
	return db.cdb
}

// Close DB
func (db *DB) Close() error {
	return db.Close()
}
