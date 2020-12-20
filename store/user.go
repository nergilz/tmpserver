package store

import (
	"errors"

	"github.com/nergilz/tmpserver/database"
)

// UserStore ..
type UserStore struct {
	db        *database.DB
	JWTSecret []byte
}

// UserModel ...
type UserModel struct {
	ID       int64  `json:"id"`
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
	var id int64
	q := `INSERT INTO users (login, password, role) VALUES ($1,$2,$3) RETURNING id`
	if err := us.db.Conn().QueryRow(q, u.Login, u.Password, u.Role).Scan(&id); err != nil {
		return err
	}
	u.ID = id
	return nil
}

// Delete user
func (us *UserStore) Delete(userID int64) error {
	q := `DELETE FROM users WHERE id = $1`
	if err := us.db.Conn().QueryRow(q, userID).Err(); err != nil {
		return err
	}
	return nil
}

// FindByID ...
func (us *UserStore) FindByID(id int64) (*UserModel, error) {
	u := &UserModel{}
	if err := us.db.Conn().QueryRow(
		"SELECT id, login, password, role FROM users WHERE id = $1",
		id).Scan( // заполняет модель UserModel
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Role,
	); err != nil {
		return nil, err
	}
	return u, nil
}

// FindByLogin ...
func (us *UserStore) FindByLogin(login string) (*UserModel, error) {
	u := &UserModel{}
	if err := us.db.Conn().QueryRow(
		"SELECT id, login, password, role FROM users WHERE login = $1",
		login).Scan( // заполняет модель UserModel
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Role,
	); err != nil {
		return nil, err
	}
	return u, nil
}

// Validate ...
func (u *UserModel) Validate() error {
	if u.Login == "" {
		return errors.New("Login cannnot de empty")
	}
	if u.Password == "" {
		return errors.New("Password cannnot de empty")
	}
	return nil
}

// GetSecret ...
func (us *UserStore) GetSecret() []byte {
	us.JWTSecret = []byte("captainjacksparrowsayshi")
	return us.JWTSecret
}
