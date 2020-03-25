package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/furkanpala/post-app/internal/core"
	"github.com/gorilla/mux"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3000"
	}
	router := mux.NewRouter()

	router.HandleFunc("/", getIndex).Methods("GET")
	router.HandleFunc("/login", handleLogin).Methods("POST")
	router.HandleFunc("/register", handleRegister).Methods("POST")
	// router.HandleFunc("/logout", handleLogout).Methods("POST")

	//TODO: Handle error
	db, _ := sql.Open("sqlite3", "./post-app.db")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var user core.User
	var hasher core.Hasher
	hasher = &user

	//TODO: Handle error
	_ = json.NewDecoder(r.Body).Decode(&user)

	hasher.HashPassword()

	//TODO: Handle error

}

func handleLogin(w http.ResponseWriter, r *http.Request) {

}

func getIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
