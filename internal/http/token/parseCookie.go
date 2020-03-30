package jwttoken

import (
	"net/http"

	httperror "github.com/furkanpala/post-app/internal/http/error"
)

// ParseCookie function extracts the cookie with given name
func ParseCookie(r *http.Request, name string) (*http.Cookie, *httperror.HTTPError) {
	cookie, err := r.Cookie(name)

	// If there is no cookie named "jid"
	// returns Forbidden error
	if err != nil {
		return nil, &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Forbidden",
				Detail: err.Error(),
			},
			Code: 403,
		}
	}

	return cookie, nil
}
