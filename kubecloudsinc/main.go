// main.go or any initialization part of your REST API package

package main

import (
	"autotools-golang-api/kubecloudsinc/dbs" // Replace with the path to your db package
	"autotools-golang-api/kubecloudsinc/middleware"
	"autotools-golang-api/kubecloudsinc/server"
	"fmt"

	"log"
	"os"

	"github.com/joho/godotenv"
)

var dsn, appName, appKey string

func init() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is not set")
	}
	appName = os.Getenv("NewRelic_AppName")
	if appName == "" {
		log.Fatal("NewRelic_AppName is not set")
	}
	appKey = os.Getenv("NewRelic_Key")
	if appKey == "" {
		log.Fatal("NewRelic_Key is not set")
	}
	fmt.Println("Database DSN:", dsn)
	fmt.Println("New Relic App Name:", appName)
	fmt.Println("New Relic App Key:", appKey)
}
func main() {
	err := dbs.InitDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Initialize New Relic
	app, err := middleware.InitNewRelic(appName, appKey)
	if err != nil {
		log.Fatal("Failed to initialize New Relic:", err)
	}

	// Start the server on port 8080
	err = server.StartServer(":8080", app)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
