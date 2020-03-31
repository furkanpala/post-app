package response

import "github.com/furkanpala/post-app/internal/core"

type PostsResponse struct {
	Posts []core.Post `json:"posts"`
	Count int         `json:"count"`
}
