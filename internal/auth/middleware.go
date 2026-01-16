package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type responseWriter struct{
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int){
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{
				ResponseWriter: w,
				status: http.StatusOK,
			}

			next.ServeHTTP(w,r)

			duration := time.Since(start)

			logger.Info("http request",
				zap.String("Method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.status),
				zap.Duration("duration", duration),
		)

		})
	}
}

func JWTMiddleware(secret []byte) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		
			authHeader := r.Header.Get("Authorization")
			if authHeader == ""{
				http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
				return
			}


			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer"{
				http.Error(w, "Invalid Headers", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(
				tokenString,
				&Claims{},
				func(token *jwt.Token) (any, error) {
					return secret, nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(w, "Invalid Header", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				http.Error(w,"Invalid Headers", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), IdContextKey, claims.UserId)
			fmt.Println("this happened")
			next.ServeHTTP(w, r.WithContext(ctx))

		})	
	}
}

