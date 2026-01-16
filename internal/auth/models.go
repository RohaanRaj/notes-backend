package auth

import (

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const EmailContextKey contextKey = "email"
const IdContextKey contextKey = "user_id"

type LoginReq struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthTokens struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct{
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}
