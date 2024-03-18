package server

import (
	// Adjust this import path to your project structure
	"autotools-golang-api/kubecloudsinc/backend/handler"
	"autotools-golang-api/kubecloudsinc/backend/middleware"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Initialize and return a new HTTP router
func NewRouter(app *newrelic.Application) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v2/login", middleware.Login).Methods("POST")
	r.HandleFunc("/v2/employees", middleware.IsAuthorized("admin", "editor", "viewer")(handler.GetEmployees)).Methods("GET")
	r.HandleFunc("/v2/employee", middleware.IsAuthorized("admin", "editor", "viewer")(handler.GetEmployee)).Methods("GET")
	r.HandleFunc("/v2/employee/{employeeId}", middleware.IsAuthorized("admin", "editor", "viewer")(handler.GetEmployeeProfile)).Methods("GET")
	r.HandleFunc("/v2/employee", middleware.IsAuthorized("admin", "editor")(handler.AddEmployee)).Methods("POST")
	r.HandleFunc("/v2/employee/{employeeId}", middleware.IsAuthorized("admin", "editor")(handler.UpdateEmployee)).Methods("PUT")
	r.HandleFunc("/v2/employee/{employeeId}", middleware.IsAuthorized("admin")(handler.DeleteEmployee)).Methods("DELETE")

	// Manually register pprof handlers
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register other pprof handlers
	r.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)

	return r
}

// StartServer starts the HTTP server on a specified port
func StartServer(port string, app *newrelic.Application) error {
	r := NewRouter(app)
	//loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	// Setup CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // The frontend origin
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	//http.Handle("/", r)
	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(r))
}
