package main

import (
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

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
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

	spa := spaHandler{staticPath: "./web/post-app/dist", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
