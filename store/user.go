package store

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/nergilz/tmpserver/database"
)

// UserStore ..
type UserStore struct {
	db *database.DB
}

// UserModel ...
type UserModel struct {
	ID       int
	Login    string `json:"login"`
	Password string
	Role     string `json:"role"`
}

// RequestUserModel ...
type RequestUserModel struct {
}

// InitUserStore ..
func InitUserStore(db *database.DB) *UserStore {
	us := new(UserStore)
	us.db = db
	return us
}

// Create user ...
func (us *UserStore) Create(u *UserModel) error {
	var id int
	q := `INSERT INTO users (login, password, role) VALUES ($1,$2,$3) RETURNING id`
	if err := us.db.Conn().QueryRow(q, u.Login, u.Password, u.Role).Scan(&id); err != nil {
		return err
	}
	u.ID = id
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

// CheckByLogin ..
func (us *UserStore) CheckByLogin(login string) bool {
	err := us.db.Conn().QueryRow("SELECT id, login, password FROM users WHERE login = $1", login)
	if err != nil {
		return false
	}
	return true
}

func hashPassword(s string) (string, error) {
	data := []byte(s)
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}
	s = hex.EncodeToString(hash.Sum(nil))
	return s, nil
}
