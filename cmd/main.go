package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkanpala/post-app/internal/database"
	httphandlers "github.com/furkanpala/post-app/internal/http/handlers"
	"github.com/furkanpala/post-app/internal/http/middleware"

	"github.com/gorilla/handlers"

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

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	credentials := handlers.AllowCredentials()

	// Auth
	router.Handle("/login", httphandlers.RouteHandler(httphandlers.HandleLogin)).Methods("POST")
	router.Handle("/register", httphandlers.RouteHandler(httphandlers.HandleRegister)).Methods("POST")
	router.Handle("/token", httphandlers.RouteHandler(httphandlers.RefreshToken)).Methods("POST")
	router.Handle("/token/logout", httphandlers.RouteHandler(httphandlers.HandleLogout)).Methods("POST")

	// Post API
	router.Handle("/posts", httphandlers.RouteHandler(httphandlers.GetPosts)).Methods("GET")
	router.Handle("/posts/amount", httphandlers.RouteHandler(httphandlers.GetPostsAmount)).Methods("GET")
	router.Handle("/posts/{page}", httphandlers.RouteHandler(httphandlers.GetPostsOnPage)).Methods("GET")
	router.Handle("/posts", middleware.AuthMiddleware(httphandlers.RouteHandler(httphandlers.AddPost)))

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins, credentials)(router)))
}
