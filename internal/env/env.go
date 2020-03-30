package env

import "os"

// AccessTokenSecret holds the value of access token secret key
var AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

// RefreshTokenSecret holds the value of refresh token secret key
var RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
