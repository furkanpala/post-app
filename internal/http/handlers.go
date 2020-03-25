package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/furkanpala/post-app/internal/core"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user core.User
	var hasher core.Hasher
	hasher = &user

	//TODO: Handle error
	_ = json.NewDecoder(r.Body).Decode(&user)

	hasher.HashPassword()

	//TODO: Handle error
	statement, _ := 

}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
