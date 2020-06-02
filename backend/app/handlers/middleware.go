package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	jwtTokenLifetimeInDays = 7
)

type contextKey string

var usernameCtxKey = contextKey("username")

// LogRequest logs all the incoming requests.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// RecoverPanic recovers the server from panic (if there is any) and responds with internal server error.
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("connection", "close")
				log.Printf("%s\n%s\n", fmt.Errorf("%s", err), debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Authorize decodes the jwt token and authorizes the request by providing the username in the
// Context() of *http.Request.
func (h *Handler) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// The string is in "bearer tokenStringHere", so the token needs to be extracted.
		tokenStr := strings.Split(authHeader[0], " ")[1]
		username, err := decodeJWT(tokenStr, h.JWTSecretKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err)
			return
		}

		// Save the username in the context so that the next handler can retreive it.
		ctx := context.WithValue(r.Context(), usernameCtxKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// decodeJWT validates the token. If successful, it returns the decoded username.
func decodeJWT(tokenStr, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse or validate the token")
	}

	// Get the claims.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if timestamp, ok := claims["iat"].(float64); ok {
			expTime := time.Unix(int64(timestamp), 0).Add(time.Hour * 24 * jwtTokenLifetimeInDays)
			if expTime.Sub(time.Now()) < 0 {
				return "", fmt.Errorf("the token is expired, please login again")
			}

			// Return the username.
			if username, ok := claims["username"].(string); ok {
				return username, nil
			}
		}

		return "", fmt.Errorf("error retreiving the claims of the token")
	}

	return "", fmt.Errorf("the authorization token is invalid")
}