package database

import (
	"database/sql"

	"github.com/0LuigiCode0/Library/logger"
	// sql driver
	_ "github.com/lib/pq"
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
		c.log.Errorf("error open DB : %v", err)
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		c.log.Errorf("error ping DB : %v", err)
		return nil, err
	}
	c.log.Service("open DB & ping DB")
	return &DB{
		Config: c,
		cdb:    conn,
	}, nil
}

// Init ..
func (db *DB) Init() error {
	q := `CREATE TABLE IF NOT EXISTS users (
		id bigserial not null primary key,
		email varchar not null unique,
		password varchar not null,
		role varchar not null
	)`
	_, err := db.cdb.Exec(q)
	if err != nil {
		db.log.Errorf("not create users tabel : %v", err)
		return err
	}
	db.log.Service("DB init: create table: users")

	q = `INSERT INTO users (email, password, role) VALUES ($1, $2, $3)`
	_, err = db.cdb.Exec(q, "admin@mail.com", "qwerty", "super_admin")
	if err != nil {
		db.log.Warningf("not insert test user on tabel : %v", err)
		return nil
	}
	db.log.Service("DB init: create test user")

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
