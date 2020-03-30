package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	httperror "github.com/furkanpala/post-app/internal/http/error"
)

// RouteHandler is a custom handler function that returns a custom HTTP Error
type RouteHandler func(http.ResponseWriter, *http.Request) *httperror.HTTPError

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
	if err != nil {
		fmt.Printf("%v\n", err.Cause)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	w.Write(body)
}
