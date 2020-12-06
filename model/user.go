package model

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
	Email    string `json:"email"`
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

// Create user ...
func (us *UserStore) Create(u *UserModel) error {
	var id int
	q := `INSERT INTO users (email, password, role) VALUES ($1,$2,$3) RETURNING id`
	if err := us.db.Conn().QueryRow(q, u.Email, u.Password, u.Role).Scan(&id); err != nil {
		return err
	}
	u.ID = id
	return nil
}

// // FindByEmail ...
// func (us *UserStore) FindByEmail(email string) (*UserModel, error) {
// 	u := &UserModel{}
// 	if err := us.db.Conn().QueryRow(
// 		"SELECT id, login, password FROM users WHERE email = $1",
// 		email).Scan( // заполняет модель UserModel
// 		&u.ID,
// 		&u.Email,
// 		&u.Password,
// 	); err != nil {
// 		return nil, err
// 	}
// 	return u, nil
// }
