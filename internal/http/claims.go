package handlers

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Username string
	jwt.StandardClaims
}
