package user

import (
	"crypto/sha256"
	"encoding/hex"
)

// Roles
const (
	RoleSuperUser = "super_user"
	RoleUser      = "user"
)

// User model
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string
}

func (u *User) hashPass() error {
	data := []byte(u.Password)
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return err
	}
	u.Password = hex.EncodeToString(hash.Sum(nil))

	return nil
}
