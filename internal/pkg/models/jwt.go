package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

/* JWT DATA */

type JwtData struct {
	jwt.StandardClaims
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (data *JwtData) Marshal(lifetime time.Duration, secret []byte) (string, error) {
	data.StandardClaims.ExpiresAt = time.Now().Add(lifetime).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(secret)
}

func (data *JwtData) UnMarshal(tokenString string, secret []byte) error {
	token, err := jwt.ParseWithClaims(tokenString, data, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if _, ok := token.Claims.(*JwtData); !ok || !token.Valid {
		return err
	}
	return nil
}
