package httphandlers

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

	// Check username and password length
	usernameLength, passwordLength := len(user.Username), len(user.Password)
	errorMessage := ""
	if usernameLength < 3 {
		errorMessage += "Too short username - Minimum 3 characters"
	}
	if passwordLength < 6 {
		if errorMessage != "" {
			errorMessage += "|"
		}
		errorMessage += "Too short password - Minimum 6 characters"
	}
	if usernameLength > 20 {
		if errorMessage != "" {
			errorMessage += "|"
		}
		errorMessage += "Too long password - Maximum 20 characters"
	}

	if errorMessage != "" {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Invalid register credentials",
				Detail: errorMessage,
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
			Code: 409,
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

	// User created
	w.WriteHeader(201)

	// Register successful

	return nil
}
