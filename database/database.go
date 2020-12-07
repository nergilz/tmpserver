package database

import (
	"database/sql"

	"github.com/0LuigiCode0/Library/logger"

	_ "github.com/lib/pq" // sql driver
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
func Connect(c *Config) (*DB, error) {
	conn, err := sql.Open("postgres", c.DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	return &DB{
		Config: c,
		cdb:    conn,
	}, nil
}

// Init ..
func (db *DB) Init() error {
	q := `CREATE TABLE IF NOT EXISTS users (
		id bigserial not null primary key,
		login varchar not null unique,
		password varchar not null,
		role varchar not null
	)`
	_, err := db.cdb.Exec(q)
	if err != nil {
		return err
	}
	// create superUser
	q = `INSERT INTO users (login, password, role) VALUES ($1, $2, $3)`
	_, err = db.cdb.Exec(q, "admin@mail.com", "qwerty", "super_admin")
	if err != nil {
		db.log.Warningf("not create superuser: %v", err)
	}
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
