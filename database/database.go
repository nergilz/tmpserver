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
func (db *DB) Init(passwordSuperUser string) error {
	// create table users
	qUsers := `CREATE TABLE IF NOT EXISTS users (
		id bigserial not null PRIMARY KEY,
		login varchar not null unique,
		password varchar not null,
		role varchar not null
	)`
	_, err := db.cdb.Exec(qUsers)
	if err != nil {
		db.log.Errorf("not create table users: %v", err)
		return err
	}
	db.log.Service("init table users")

	// // create table chats
	// qChats := `CREATE TABLE IF NOT EXISTS chats (
	// 	id BIGSERIAL NOT NULL PRIMARY KEY,
	// 	user_id INTEGER[],
	// 	private BOOL
	// )`
	// _, err = db.cdb.Exec(qChats)
	// if err != nil {
	// 	db.log.Errorf("not create table chats", err)
	// 	return err
	// }
	// db.log.Service("init table chats")

	// create table messages
	qMessages := `CREATE TABLE IF NOT EXISTS messages (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		to_id BIGINT REFERENCES users (id),
		from_id INTEGER REFERENCES users (id),
		content TEXT NOT NULL
	)`
	_, err = db.cdb.Exec(qMessages)
	if err != nil {
		db.log.Errorf("not create table messages: %v", err)
		return err
	}
	db.log.Service("init table messages")

	// create superUser
	qSuperUser := `INSERT INTO users (login, password, role) VALUES ($1, $2, $3)`
	_, err = db.cdb.Exec(qSuperUser, "admin", passwordSuperUser, "super_user")
	if err != nil {
		db.log.Warningf("not create superuser: %v", err)
	} else {
		db.log.Service("create super_user")
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
