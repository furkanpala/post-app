package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/furkanpala/post-app/internal/database"
	"github.com/furkanpala/post-app/internal/env"
	httperror "github.com/furkanpala/post-app/internal/http/error"
	"github.com/furkanpala/post-app/internal/http/response"
	jwttoken "github.com/furkanpala/post-app/internal/http/token"
)

// RefreshToken handles the requests for /token route.
// Verifies the incoming refresh token.
// If it is valid, then responses with a new refresh and an access token
func RefreshToken(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {

	cookie, httpErr := jwttoken.ParseCookie(r, "jid")
	if httpErr != nil {
		return httpErr
	}

	claims := &jwttoken.Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(env.RefreshTokenSecret), nil
	})

	if err != nil {
		switch err.(type) {
		// If token is not valid, return Unauthorized error
		case *jwt.ValidationError:
			return &httperror.HTTPError{
				Cause: err,
				Info: httperror.ErrorMessage{
					Title:  "Unauthorized",
					Detail: "Invalid credentials",
				},
				Code: 401,
			}
		default:
			return &httperror.HTTPError{
				Cause: err,
				Info: httperror.ErrorMessage{
					Title:  "Internal server error",
					Detail: "",
				},
				Code: 500,
			}
		}
	}
	// If token is not valid, return Unauthorized error
	if !token.Valid {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Check if the given token's jti is in blacklist
	isInBlacklist, err := database.FindJTI(claims.Id)
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	fmt.Printf("refresh: %v\n", isInBlacklist)

	// If token's jti is in blacklist,
	// meaning that if token is already revoked,
	// returns Unauthorized error
	if isInBlacklist {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Generate access token with 15 minutes of expire time
	accessTokenString, err := jwttoken.GenerateToken(15*time.Minute, claims.Username, env.AccessTokenSecret)
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Generate refresh token with 7 day of expire time
	refreshTokenString, err := jwttoken.GenerateToken(168*time.Hour, claims.Username, env.RefreshTokenSecret)
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Set headers according to OAuth 2.0
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	responseBody := response.SuccessfulLoginResponse{
		AccessToken: accessTokenString,
		TokenType:   "bearer",
		ExpiresIn:   15 * 60 * 60,
	}

	// Set cookie
	newCookie := http.Cookie{
		Name:     "jid",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(168 * time.Hour),
		HttpOnly: true,
		Path:     "/token",
	}
	http.SetCookie(w, &newCookie)
	json.NewEncoder(w).Encode(responseBody)
	return nil
}
