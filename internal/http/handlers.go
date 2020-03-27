package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/furkanpala/post-app/internal/core"
	"github.com/furkanpala/post-app/internal/database"
)

type RouteHandler func(http.ResponseWriter, *http.Request) *HTTPError

func (fn RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}

	body, parseError := json.Marshal(err)

	if parseError != nil {
		w.WriteHeader(500)
		fmt.Printf("Parse error: %v", parseError)
		return
	}
	fmt.Printf("%v", err.Cause)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	w.Write(body)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) *HTTPError {
	var user core.User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&user); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Invalid JSON",
				Detail: err.Error(),
			},
			Code: 400,
		}
	}

	userExists, err := database.CheckUserExists(user.Username)

	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	if userExists {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "User already exists",
				Detail: "",
			},
			Code: 200,
		}

	}
	if err := user.HashPassword(); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	if err := database.AddUser(&user); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}
	return nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
