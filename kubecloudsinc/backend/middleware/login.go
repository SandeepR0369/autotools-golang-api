package middleware

// User represents a user with a username, password, and role.
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type contextKey string

// RoleContextKey is the key for role values in the context
const RoleContextKey contextKey = "userRole"

type User struct {
	Username string
	Password string
	Role     string
}

// Users is a mock database of users.
var users = []User{
	{Username: "mazda", Password: "Test1ng!", Role: "admin"},
	{Username: "honda", Password: "Test1ng!", Role: "editor"},
	{Username: "kia", Password: "Test1ng!", Role: "editor"},
	{Username: "benz", Password: "Test1ng!", Role: "viewer"},
	{Username: "toyota", Password: "Test1ng!", Role: "viewer"},
}

var jwtKey = []byte("JAIJAFFA")

// Credentials are used for parsing login requests.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims are used for creating JWT tokens.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authenticate the user
	var userRole string
	for _, user := range users {
		if user.Username == creds.Username && user.Password == creds.Password {
			userRole = user.Role
			break
		}
	}

	if userRole == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Log successful authentication
	log.Printf("User authenticated: %s at %s", creds.Username, time.Now().Format(time.RFC3339))

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		Role:     userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Token generated for user: %s at %s", creds.Username, time.Now().Format(time.RFC3339))

	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })
	// Instead of setting the token as a cookie, return it in the response body
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func IsAuthorized(requiredRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, "Bearer ")
			if len(bearerToken) != 2 {
				http.Error(w, "Invalid Authorization token format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "Invalid token signature", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Check if the user's role matches any of the required roles
			roleIsAllowed := false
			for _, requiredRole := range requiredRoles {
				if claims.Role == requiredRole {
					roleIsAllowed = true
					break
				}
			}

			if !roleIsAllowed {
				msg := fmt.Sprintf("Insufficient permissions: user role %s is not allowed", claims.Role)
				log.Printf("Insufficient permissions: user role %s is not allowed", claims.Role)
				http.Error(w, msg, http.StatusForbidden)
				return
			}

			// User is authorized; add the user's role to the context
			ctxWithRole := context.WithValue(r.Context(), RoleContextKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctxWithRole))
		}
	}
}
