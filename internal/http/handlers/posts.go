package httphandlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/furkanpala/post-app/internal/core"
	"github.com/furkanpala/post-app/internal/database"
	httperror "github.com/furkanpala/post-app/internal/http/error"
	"github.com/furkanpala/post-app/internal/http/request"
	"github.com/furkanpala/post-app/internal/http/response"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const PostsPerPage = 6

func GetPosts(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
	posts, err := database.GetAllPosts()
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}

	responseBody := response.PostsResponse{
		Posts: posts,
		Count: len(posts),
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}

	return nil
}

// GetPostsOnPage returns a slice of posts which are on a specific page.
func GetPostsOnPage(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
	params := mux.Vars(r)
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Invalid page",
				Detail: err.Error(),
			},
			Code: 400,
		}
	}

	if page <= 0 {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Page not found",
				Detail: "Page must be an integer greater than zero",
			},
			Code: 404,
		}
	}

	postsCount, err := database.CountPosts()
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}

	firstPostIndex := PostsPerPage * (page - 1)
	if firstPostIndex+1 > postsCount {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Page not found",
				Detail: "",
			},
			Code: 404,
		}
	}

	lastPostIndex := PostsPerPage*page - 1

	if lastPostIndex+1 > postsCount {
		lastPostIndex = postsCount - 1
	}

	posts, err := database.GetAllPosts()
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}
	posts = posts[firstPostIndex : lastPostIndex+1]

	responseBody := response.PostsResponse{
		Posts: posts,
		Count: len(posts),
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}

	return nil
}

func AddPost(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
	var post core.Post
	if httpErr := request.DecodeRequestBody(r, &post); httpErr != nil {
		return httpErr
	}

	username := context.Get(r, "username")
	post.User = username.(string)

	errorMessage := ""
	if len(post.Title) == 0 {
		errorMessage += "Title required"
	}
	if len(post.Content) == 0 {
		if errorMessage != "" {
			errorMessage += "|"
		}
		errorMessage += "Content required"
	}

	if errorMessage != "" {
		return &httperror.HTTPError{
			Cause: nil,
			Info: httperror.ErrorMessage{
				Title:  "Invalid post info",
				Detail: errorMessage,
			},
			Code: 400,
		}
	}

	if err := database.AddPost(&post); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Interal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}
	w.WriteHeader(201)

	return nil
}

func GetPostsAmount(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
	count, err := database.CountPosts()
	if err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}

	responseBody := response.PostsResponse{
		Posts: nil,
		Count: count,
	}

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		return &httperror.HTTPError{
			Cause: err,
			Info: httperror.ErrorMessage{
				Title:  "Internal server error",
				Detail: err.Error(),
			},
			Code: 500,
		}
	}

	return nil
}
