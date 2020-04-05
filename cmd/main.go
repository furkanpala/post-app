package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/furkanpala/post-app/internal/database"
	httphandlers "github.com/furkanpala/post-app/internal/http/handlers"
	"github.com/furkanpala/post-app/internal/http/middleware"

	"github.com/gorilla/mux"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3000"
	}

	err := database.OpenDatabase()
	if err != nil {
		log.Fatal("Database error")
	}

	defer database.CloseDatabase()

	err = database.CreateUsersTable()
	if err != nil {
		log.Fatal("Database error")
	}

	err = database.CreatePostsTable()
	if err != nil {
		log.Fatal("Database error")
	}

	err = database.CreateBlacklistTable()
	if err != nil {
		log.Fatal("Database error")
	}

	router := mux.NewRouter()

	// headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// methods := handlers.AllowedMethods([]string{"GET", "POST"})
	// origins := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	// credentials := handlers.AllowCredentials()

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

	spa := spaHandler{staticPath: "dist", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Server is running on port " + port)
	log.Fatal(srv.ListenAndServe())

}
