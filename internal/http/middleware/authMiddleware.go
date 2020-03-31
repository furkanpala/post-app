package middleware

import (
	"net/http"
	"strings"

	"github.com/furkanpala/post-app/internal/env"
	httperror "github.com/furkanpala/post-app/internal/http/error"
	"github.com/furkanpala/post-app/internal/http/handlers"
	jwttoken "github.com/furkanpala/post-app/internal/http/token"
	"github.com/gorilla/context"
)

func AuthMiddleware(next handlers.RouteHandler) handlers.RouteHandler {
	return handlers.RouteHandler(func(w http.ResponseWriter, r *http.Request) *httperror.HTTPError {
		// authorization := strings.Split(strings.Join(r.Header["Authorization"], ""), " ")
		// if authorization

		authorization := strings.Split(r.Header.Get("Authorization"), " ")

		if len(authorization) != 2 {
			return &httperror.HTTPError{
				Cause: nil,
				Info: httperror.ErrorMessage{
					Title:  "Unauthorized",
					Detail: "",
				},
				Code: 401,
			}
		}

		accessTokenString := authorization[1]

		var claims jwttoken.Claims

		_, httpErr := jwttoken.VerifyToken(accessTokenString, env.AccessTokenSecret, &claims)
		if httpErr != nil {
			return httpErr
		}

		context.Set(r, "username", claims.Username)

		next.ServeHTTP(w, r)
		return nil
	})
}
