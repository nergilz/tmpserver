package store

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/nergilz/tmpserver/database"
)

// UserStore ..
type UserStore struct {
	db        *database.DB
	jwtSecret []byte
}

// UserModel ...
type UserModel struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// InitUserStore ..
func InitUserStore(db *database.DB) *UserStore {
	us := new(UserStore)
	us.db = db
	return us
}

// Create user
func (us *UserStore) Create(u *UserModel) error {
	var id int
	q := `INSERT INTO users (login, password, role) VALUES ($1,$2,$3) RETURNING id`
	if err := us.db.Conn().QueryRow(q, u.Login, u.Password, u.Role).Scan(&id); err != nil {
		return err
	}
	u.ID = id
	return nil
}

// Delete user
func (us *UserStore) Delete(userID int) error {
	q := `DELETE FROM users WHERE id = $1`
	if err := us.db.Conn().QueryRow(q, userID).Err(); err != nil {
		return err
	}
	return nil
}

// FindByLogin ...
func (us *UserStore) FindByLogin(login string) (*UserModel, error) {
	u := &UserModel{}
	if err := us.db.Conn().QueryRow(
		"SELECT id, login, password FROM users WHERE login = $1",
		login).Scan( // заполняет модель UserModel
		&u.ID,
		&u.Login,
		&u.Password,
	); err != nil {
		return nil, err
	}
	return u, nil
}

// HashPassword ...
func HashPassword(s string) (string, error) {
	data := []byte(s)
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}
	s = hex.EncodeToString(hash.Sum(nil))
	return s, nil
}

// Validate ...
func (u *UserModel) Validate() error {
	if u.Login == "" {
		return errors.New("Login cannnot de empty")
	}
	if u.Password == "" {
		return errors.New("Password cannnot de empty")
	}
	if u.Role == "" || u.Role != "user" {
		return errors.New("Role must be user")
	}
	return nil
}

// GetSecret ...
func (us *UserStore) GetSecret() []byte {
	us.jwtSecret = []byte("captainjacksparrowsayshi")
	return us.jwtSecret
}
