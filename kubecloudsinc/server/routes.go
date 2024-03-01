package server

import (
	"autotools-golang-api/kubecloudsinc/handler" // Adjust this import path to your project structure
	"autotools-golang-api/kubecloudsinc/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Initialize and return a new HTTP router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Public route
	r.HandleFunc("/login", middleware.Login).Methods("POST")
	// r.HandleFunc("/employees",middleware.IsAuthorized(handler.GetEmployees, "admin")).Methods("GET")
	r.HandleFunc("/employees", middleware.IsAuthorized("admin")(handler.GetEmployees)).Methods("GET")
	return r
}

// StartServer starts the HTTP server on a specified port
func StartServer(port string) error {
	r := NewRouter()
	return http.ListenAndServe(port, r)
}
