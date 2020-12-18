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

// Init init of database
func (db *DB) Init() error {
	qUsers := `CREATE TABLE IF NOT EXISTS users (
		id bigserial not null PRIMARY KEY,
		login varchar not null unique,
		password varchar not null,
		role varchar not null
	)`
	_, err := db.cdb.Exec(qUsers)
	if err != nil {
		db.log.Errorf("not create users: %v", err)
		return err
	}
	db.log.Service("Init table users")
	qMessages := `CREATE TABLE IF NOT EXISTS messages (
		id bigserial not null PRIMARY KEY,
		owner_id BIGINT REFERENCES users (id),
		user_id INTEGER,
		description VARCHAR(100),
		text TEXT not null
	)`
	_, err = db.cdb.Exec(qMessages)
	if err != nil {
		db.log.Errorf("not create messages: %v", err)
		return err
	}
	db.log.Service("Init table messages")
	// create superUser
	qSuperUser := `INSERT INTO users (login, password, role) VALUES ($1, $2, $3)`
	_, err = db.cdb.Exec(qSuperUser, "admin", "q1w2e3r4t5y6", "super_user")
	if err != nil {
		db.log.Warningf("not create superuser: %v", err)
	} else {
		db.log.Service("create Super User")
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
