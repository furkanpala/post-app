package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkanpala/post-app/internal/database"
	handlers "github.com/furkanpala/post-app/internal/http"

	"github.com/gorilla/mux"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3000"
	}

	//TODO: Handle error
	_ = database.OpenDatabase()

	defer database.CloseDatabase()

	//TODO: Handle error
	_ = database.CreateUsersTable()

	//TODO: Handle error
	_ = database.CreatePostsTable()

	//TODO: Handle error
	_ = database.CreateLikesTable()

	router := mux.NewRouter()
	// router.Handle("/", handlers.RouteHandler(handlers.GetIndex)).Methods("GET")
	// router.Handle("/login", handlers.HandleLogin).Methods("POST")
	router.Handle("/register", handlers.RouteHandler(handlers.HandleRegister)).Methods("POST")
	// router.HandleFunc("/logout", handleLogout).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
