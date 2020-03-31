package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkanpala/post-app/internal/database"
	"github.com/furkanpala/post-app/internal/http/handlers"
	"github.com/furkanpala/post-app/internal/http/middleware"

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

	//TODO: Handle error
	_ = database.CreateBlacklistTable()

	router := mux.NewRouter()

	// Auth
	router.Handle("/login", handlers.RouteHandler(handlers.HandleLogin)).Methods("POST")
	router.Handle("/register", handlers.RouteHandler(handlers.HandleRegister)).Methods("POST")
	router.Handle("/token", handlers.RouteHandler(handlers.RefreshToken)).Methods("POST")
	router.Handle("/token/logout", handlers.RouteHandler(handlers.HandleLogout)).Methods("POST")

	// Post API
	router.Handle("/posts", handlers.RouteHandler(handlers.GetPosts)).Methods("GET")
	router.Handle("/posts/{page}", handlers.RouteHandler(handlers.GetPostsOnPage)).Methods("GET")
	router.Handle("/posts", middleware.AuthMiddleware(handlers.RouteHandler(handlers.AddPost)))
	// protected /posts POST
	// req -> {title,content,user,date*,id*}
	// res -> 201 && all posts on first page
	// on client redirect to first page

	log.Fatal(http.ListenAndServe(":"+port, router))
}
