package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/furkanpala/post-app/internal/core"
	"github.com/furkanpala/post-app/internal/database"
	"github.com/furkanpala/post-app/internal/env"
	uuid "github.com/satori/go.uuid"
)

// generateToken function generates a token with given
// expire time
// username as payload
// secret key to sign the token
// token also includes an jti generated with uuid
func generateToken(expireTime time.Duration, username, secret string) (string, error) {
	claims := &Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireTime).Unix(),
			Id:        uuid.NewV4().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}

// decodeRequestBody function decodes json request r and stores to the value pointed by v
func decodeRequestBody(r *http.Request, v interface{}) *HTTPError {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Invalid JSON",
				Detail: err.Error(),
			},
			Code: 400,
		}
	}

	return nil
}

// validateUser function checks if user's credentials valid for log in
func validateUser(user *core.User) *HTTPError {
	// Check if user exists in database
	dbUser, err := database.FindUser(user.Username)
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	if dbUser == nil {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Check plain text password against hashed password
	if err := dbUser.Compare(user.Password); err != nil {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Validation successful

	return nil
}

// parseCookie function extracts the cookie with given name
func parseCookie(r *http.Request, name string) (*http.Cookie, *HTTPError) {
	cookie, err := r.Cookie(name)

	// If there is no cookie named "jid"
	// returns Forbidden error
	if err != nil {
		return nil, &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Forbidden",
				Detail: err.Error(),
			},
			Code: 403,
		}
	}

	return cookie, nil
}

// verifyToken function checks if the given token is valid or not
func verifyToken(tokenString, secret string, claims jwt.Claims) (*jwt.Token, *HTTPError) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			return nil, &HTTPError{
				Cause: err,
				Info: ErrorMessage{
					Title:  "Unauthorized",
					Detail: "Invalid credentials",
				},
				Code: 401,
			}
		default:
			return nil, &HTTPError{
				Cause: err,
				Info: ErrorMessage{
					Title:  "Internal server error",
					Detail: "",
				},
				Code: 500,
			}
		}
	}
	if !token.Valid {
		return nil, &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	return token, nil
}

// RouteHandler is a custom handler function that returns a custom HTTP Error
type RouteHandler func(http.ResponseWriter, *http.Request) *HTTPError

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

// HandleRegister function handles the request for /register route.
// Adds the user into database if request is correct
func HandleRegister(w http.ResponseWriter, r *http.Request) *HTTPError {
	var user core.User

	// Parse request body
	if err := decodeRequestBody(r, &user); err != nil {
		return err
	}

	usernameLength, passwordLength := len(user.Username), len(user.Password)

	// Check username length
	if usernameLength < 3 {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Too short username",
				Detail: "Minimum 3 characters",
			},
			Code: 400,
		}
	}

	// Check password length
	if passwordLength < 6 {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Too short password",
				Detail: "Minimum 6 characters",
			},
			Code: 400,
		}
	}

	// Check if users already exists
	userExists, err := database.FindUser(user.Username)

	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	if userExists != nil {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "User already exists",
				Detail: "",
			},
			Code: 200,
		}
	}

	// Hash user's password
	if err := user.HashPassword(); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Add into database
	if err := database.AddUser(&user); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Register successful

	return nil
}

// HandleLogin function handles the request for /login route.
// Check credentials against database.
// If correct, responses with an access token and a refresh token
func HandleLogin(w http.ResponseWriter, r *http.Request) *HTTPError {
	var user core.User

	// Parse request body
	if err := decodeRequestBody(r, &user); err != nil {
		return err
	}

	// Check credentials
	if err := validateUser(&user); err != nil {
		return err
	}

	// Generate access token with 15 minutes of expire time
	accessTokenString, err := generateToken(15*time.Minute, user.Username, env.AccessTokenSecret)
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Generate refresh token with 7 days of expire time
	refreshTokenString, err := generateToken(168*time.Hour, user.Username, env.RefreshTokenSecret)
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Set headers according to OAuth 2.0
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	responseBody := SuccessfulLoginResponse{
		AccessToken: accessTokenString,
		TokenType:   "bearer",
		ExpiresIn:   15 * 60 * 60,
	}

	// Set cookie
	cookie := http.Cookie{
		Name:     "jid",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(168 * time.Hour),
		HttpOnly: true,
		Path:     "/token",
	}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(responseBody)

	// Login successful

	return nil
}

// RefreshToken handles the requests for /token route.
// Verifies the incoming refresh token.
// If it is valid, then responses with a new refresh and an access token
func RefreshToken(w http.ResponseWriter, r *http.Request) *HTTPError {

	cookie, httpErr := parseCookie(r, "jid")
	if httpErr != nil {
		return httpErr
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(env.RefreshTokenSecret), nil
	})

	if err != nil {
		switch err.(type) {
		// If token is not valid, return Unauthorized error
		case *jwt.ValidationError:
			return &HTTPError{
				Cause: err,
				Info: ErrorMessage{
					Title:  "Unauthorized",
					Detail: "Invalid credentials",
				},
				Code: 401,
			}
		default:
			return &HTTPError{
				Cause: err,
				Info: ErrorMessage{
					Title:  "Internal server error",
					Detail: "",
				},
				Code: 500,
			}
		}
	}
	// If token is not valid, return Unauthorized error
	if !token.Valid {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Check if the given token's jti is in blacklist
	isInBlacklist, err := database.FindJTI(claims.Id)
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// If token's jti is in blacklist,
	// meaning that if token is already revoked,
	// returns Unauthorized error
	if isInBlacklist {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Generate access token with 15 minutes of expire time
	accessTokenString, err := generateToken(15*time.Minute, claims.Username, env.AccessTokenSecret)
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Generate refresh token with 7 day of expire time
	refreshTokenString, err := generateToken(168*time.Hour, claims.Username, env.RefreshTokenSecret)
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	// Set headers according to OAuth 2.0
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	responseBody := SuccessfulLoginResponse{
		AccessToken: accessTokenString,
		TokenType:   "bearer",
		ExpiresIn:   15 * 60 * 60,
	}

	// Set cookie
	newCookie := http.Cookie{
		Name:     "jid",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(168 * time.Hour),
		HttpOnly: true,
		Path:     "/token",
	}
	http.SetCookie(w, &newCookie)
	json.NewEncoder(w).Encode(responseBody)
	return nil
}

// HandleLogout function handles the requests to /token/logout route.
// If given refresh token is valid, then user is logged out.
// Refresh token's jti is added to the blacklist in database.
// Responses with empty refresh token.
func HandleLogout(w http.ResponseWriter, r *http.Request) *HTTPError {
	// Parse cookie
	cookie, httpErr := parseCookie(r, "jid")
	if httpErr != nil {
		return httpErr
	}

	claims := &Claims{}

	// Verify token
	_, httpErr = verifyToken(cookie.Value, env.RefreshTokenSecret, claims)
	if httpErr != nil {
		return httpErr
	}

	// Check if given token's jti is already in blacklist
	isInBlacklist, err := database.FindJTI(claims.Id)
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}

	if isInBlacklist {
		return &HTTPError{
			Cause: nil,
			Info: ErrorMessage{
				Title:  "Unauthorized",
				Detail: "Invalid credentials",
			},
			Code: 401,
		}
	}

	// Blacklist the token through adding it's "jti" into database
	if err := database.BlacklistToken(claims.Id, claims.ExpiresAt); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Internal server error",
				Detail: "",
			},
			Code: 500,
		}
	}
	newCookie := http.Cookie{
		Name:     "jid",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/token",
	}
	http.SetCookie(w, &newCookie)
	return nil
}
