package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/furkanpala/post-app/internal/core"
	"github.com/furkanpala/post-app/internal/database"
)

// var accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
// var refreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
var accessTokenSecret = "accessTokenSecret"
var refreshTokenSecret = "refreshTokenSecret"

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

func HandleRegister(w http.ResponseWriter, r *http.Request) *HTTPError {
	var user core.User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&user); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Invalid JSON",
				Detail: err.Error(),
			},
			Code: 400,
		}
	}

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

	return nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) *HTTPError {
	var user core.User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&user); err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Invalid JSON",
				Detail: err.Error(),
			},
			Code: 400,
		}
	}

	if err := ValidateUser(&user); err != nil {
		return err
	}

	// Generate access token
	accessTokenString, err := GenerateToken(15*time.Minute, user.Username, accessTokenSecret)
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
	refreshTokenString, err := GenerateToken(168*time.Hour, user.Username, refreshTokenSecret)
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	responseBody := SuccessfulLoginResponse{
		AccessToken: accessTokenString,
		TokenType:   "bearer",
		ExpiresIn:   15 * 60 * 60,
	}
	cookie := http.Cookie{
		Name:     "jid",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(168 * time.Hour),
		HttpOnly: true,
		Path:     "/token",
	}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(responseBody)
	return nil
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func GenerateToken(expireTime time.Duration, username, secretKey string) (string, error) {
	claims := &Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireTime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}

func ValidateUser(user *core.User) *HTTPError {
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

	return nil
}

func RefreshToken(w http.ResponseWriter, r *http.Request) *HTTPError {
	// parse cookie
	cookie, err := r.Cookie("jid")
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Forbidden",
				Detail: err.Error(),
			},
			Code: 403,
		}
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenSecret), nil
	})

	if err != nil {
		switch err.(type) {
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

	isInBlacklist, err := database.FindToken(token.Raw)
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

	accessTokenString, err := GenerateToken(15*time.Minute, claims.Username, accessTokenSecret)
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
	refreshTokenString, err := GenerateToken(168*time.Hour, claims.Username, refreshTokenSecret)
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	responseBody := SuccessfulLoginResponse{
		AccessToken: accessTokenString,
		TokenType:   "bearer",
		ExpiresIn:   15 * 60 * 60,
	}
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

func HandleLogout(w http.ResponseWriter, r *http.Request) *HTTPError {
	// parse cookie
	cookie, err := r.Cookie("jid")
	if err != nil {
		return &HTTPError{
			Cause: err,
			Info: ErrorMessage{
				Title:  "Forbidden",
				Detail: err.Error(),
			},
			Code: 403,
		}
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenSecret), nil
	})

	if err != nil {
		switch err.(type) {
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

	isInBlacklist, err := database.FindToken(token.Raw)
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

	if err := database.BlacklistToken(token.Raw); err != nil {
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
}
