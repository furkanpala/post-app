package httphandlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/furkanpala/post-app/internal/database"
	"github.com/furkanpala/post-app/internal/env"
	httperror "github.com/furkanpala/post-app/internal/http/error"
	jwttoken "github.com/furkanpala/post-app/internal/http/token"
)

// HandleLogout function handles the requests to /token/logout route.
// If given refresh token is valid, then user is logged out.
// Refresh token's jti is added to the blacklist in database.
// Responses with empty refresh token.
func HandleLogout(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
	// Parse cookie
	cookie, httpErr := jwttoken.ParseCookie(r, "jid")
	if httpErr != nil {
		return httpErr
	}

	claims := &jwttoken.Claims{}

	// Verify token
	_, httpErr = jwttoken.VerifyToken(cookie.Value, env.RefreshTokenSecret, claims)
	if httpErr != nil {
		return httpErr
	}

	// Check if given token's jti is already in blacklist
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
	fmt.Printf("logout: %v\n", isInBlacklist)

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

	// Blacklist the token through adding it's "jti" into database
	if err := database.BlacklistToken(claims.Id, claims.ExpiresAt); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}
	newCookie := http.Cookie{
		Name:     "jid",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &newCookie)
	return nil
}
