// main.go or any initialization part of your REST API package

package main

import (
	"autotools-golang-api/kubecloudsinc/dbs" // Replace with the path to your db package
	"autotools-golang-api/kubecloudsinc/server"
	"log"
	"os"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is not set")
	}

	err := dbs.InitDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Start the server on port 8080
	err = server.StartServer(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
