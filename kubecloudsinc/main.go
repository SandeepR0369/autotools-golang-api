// main.go or any initialization part of your REST API package

package main

import (
	"kubecloudsinc/dbs" // Replace with the path to your db package
	"kubecloudsinc/server"
	"log"
)

func main() {
	// Example DSN (replace with actual values)

	dsn := "admin/admin123@10.13.18.130:1521/fff"

	err := dbs.InitDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Start the server on port 8080
	err = server.StartServer(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// employees, err := dbs.QueryEmployees(dbs.DB) // Assuming QueryEmployees is a function that returns ([]Employee, error)
	// if err != nil {
	// 	log.Fatalf("Failed to query employees: %v", err)
	// }

	// Print employees
	//	dbs.PrintEmployees(employees)

	// Proceed to set up your REST API routes

}
