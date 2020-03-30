package jwttoken

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// GenerateToken function generates a token with given
// expire time
// username as payload
// secret key to sign the token
// token also includes an jti generated with uuid
func GenerateToken(expireTime time.Duration, username, secret string) (string, error) {
	claims := &Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireTime).Unix(),
			Id:        uuid.NewV4().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}
