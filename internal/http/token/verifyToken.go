package jwttoken

import (
	"github.com/dgrijalva/jwt-go"
	httperror "github.com/furkanpala/post-app/internal/http/error"
)

// VerifyToken function checks if the given token is valid or not
func VerifyToken(tokenString, secret string, claims jwt.Claims) (*jwt.Token, *httperror.HTTPError) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			return nil, &httperror.HTTPError{
				Cause: err,
				Info: httperror.ErrorMessage{
					Title:  "Unauthorized",
					Detail: "Invalid credentials",
				},
				Code: 401,
			}
		default:
			return nil, &httperror.HTTPError{
				Cause: err,
				Info: httperror.ErrorMessage{
					Title:  "Internal server error",
					Detail: "",
				},
				Code: 500,
			}
		}
	}
	if !token.Valid {
		return nil, &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	return token, nil
}
