package model

import (
	"github.com/nergilz/tmpserver/database"
)

// UStore ..
type UStore struct {
	db *database.DB
}

// UModel ...
type UModel struct {
	ID       int
	Email    string
	Password string
	Role     string
}

// InitUserStore ..
func InitUserStore(db *database.DB) *UStore {
	us := new(UStore)
	us.db = db
	return us
}

// Create ..
func (us *UStore) Create(u *UModel) error {
	var id int
	q := `INSERT INTO users (email, password, role) VALUES ($1,$2,$3) RETURNING id`
	if err := us.db.Conn().QueryRow(q, u.Email, u.Password, u.Role).Scan(&id); err != nil {
		return err
	}
	u.ID = id
	return nil
}
