package jwttoken

import "github.com/dgrijalva/jwt-go"

// Claims is a custom struct for JWT tokens.
// Includes standard claims.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
