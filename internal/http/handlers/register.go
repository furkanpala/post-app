package handlers

import (
	"net/http"

	"github.com/furkanpala/post-app/internal/core"
	"github.com/furkanpala/post-app/internal/database"
	httperror "github.com/furkanpala/post-app/internal/http/error"
	"github.com/furkanpala/post-app/internal/http/request"
)

// HandleRegister function handles the request for /register route.
// Adds the user into database if request is correct
func HandleRegister(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
	var user core.User

	// Parse request body
	if err := request.DecodeRequestBody(r, &user); err != nil {
		return err
	}

	usernameLength, passwordLength := len(user.Username), len(user.Password)

	// Check username length
	if usernameLength < 3 {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Too short username",
				Detail: "Minimum 3 characters",
			},
			Code: 400,
		}
	}

	// Check password length
	if passwordLength < 6 {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Too short password",
				Detail: "Minimum 6 characters",
			},
			Code: 400,
		}
	}

	// Check if users already exists
	userExists, err := database.FindUser(user.Username)

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

	if userExists != nil {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "User already exists",
				Detail: "",
			},
			Code: 200,
		}
	}

	// Hash user's password
	if err := user.HashPassword(); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Add into database
	if err := database.AddUser(&user); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Register successful

	return nil
}
