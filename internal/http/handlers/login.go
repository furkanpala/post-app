package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/furkanpala/post-app/internal/core"
	"github.com/furkanpala/post-app/internal/database"
	"github.com/furkanpala/post-app/internal/env"
	httperror "github.com/furkanpala/post-app/internal/http/error"
	"github.com/furkanpala/post-app/internal/http/request"
	"github.com/furkanpala/post-app/internal/http/response"
	jwttoken "github.com/furkanpala/post-app/internal/http/token"
)

// HandleLogin function handles the request for /login route.
// Check credentials against database.
// If correct, responses with an access token and a refresh token
func HandleLogin(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
	var user core.User

	// Parse request body
	if err := request.DecodeRequestBody(r, &user); err != nil {
		return err
	}

	// Check credentials
	if err := validateUser(&user); err != nil {
		return err
	}

	// Generate access token with 15 minutes of expire time
	accessTokenString, err := jwttoken.GenerateToken(15*time.Minute, user.Username, env.AccessTokenSecret)
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

	// Generate refresh token with 7 days of expire time
	refreshTokenString, err := jwttoken.GenerateToken(168*time.Hour, user.Username, env.RefreshTokenSecret)
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
	cookie := http.Cookie{
		Name:     "jid",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(168 * time.Hour),
		HttpOnly: true,
		Path:     "/token",
	}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(responseBody)

	// Login successful

	return nil
}

// validateUser function checks if user's credentials valid for log in
func validateUser(user *core.User) *httperror.HTTPError {
	// Check if user exists in database
	dbUser, err := database.FindUser(user.Username)
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

	if dbUser == nil {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Check plain text password against hashed password
	if err := dbUser.Compare(user.Password); err != nil {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Validation successful

	return nil
}
