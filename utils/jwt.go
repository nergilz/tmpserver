package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nergilz/tmpserver/store"
)

// Claims ..
type Claims struct {
	jwt.StandardClaims
	UserID int    `json:"userid"`
	Login  string `json:"login"`
}

// CreateJWTtoken ...
func CreateJWTtoken(u *store.UserModel, secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["role"] = u.Role
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// // ParseJWTtoken ...
// func ParseJWTtoken(accessToken string, signinKey []byte) (string, err) {
// 	token, err := jwt.Parse(accessToken)
// }

// CheckJWTtoken ...
