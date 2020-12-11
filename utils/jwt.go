package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/nergilz/tmpserver/store"
)

// Claims : castom, not use!
type Claims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

// CreateJWTtoken ...
func CreateJWTtoken(u *store.UserModel, secret []byte) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = u.ID
	claims["login"] = u.Login
	claims["role"] = u.Role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// VerifyJWTtoken parse token
func VerifyJWTtoken(accessToken string, signinKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return signinKey, nil
	})
	if err != nil {
		return nil, err
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok && token.Valid {
		return nil, errors.New("bad auth token")
	}
	return mapClaims, nil
}

// CheckJWTtoken extract token metadata
func CheckJWTtoken(claims jwt.MapClaims) (*store.UserModel, error) {
	newErr := errors.New("Bad extract token metadata")
	userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["id"]), 10, 64)
	if err != nil {
		return nil, err
	}
	userLogin, ok := claims["login"].(string)
	if !ok {
		return nil, newErr
	}
	userRole, ok := claims["role"].(string)
	if !ok {
		return nil, newErr
	}
	u := &store.UserModel{
		ID:    userID,
		Login: userLogin,
		Role:  userRole,
	}
	return u, nil
}
