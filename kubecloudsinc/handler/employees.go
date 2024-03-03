package handler

import (
	"autotools-golang-api/kubecloudsinc/dbs"
	"autotools-golang-api/kubecloudsinc/schema"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := dbs.QueryEmployees(dbs.DB)
	if err != nil {
		log.Printf("Error querying all employees: %v", err)
		http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
		return
	}

	// Instead of printing, send the employees back as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding employees to JSON: %v", err)
		// Consider logging the error and deciding on the best HTTP status to return
		http.Error(w, "Error processing data", http.StatusInternalServerError)
	}
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	employeeIdStr := queryValues.Get("employeeId")
	lastName := queryValues.Get("lastName")

	var employeeId int
	var err error
	if employeeIdStr != "" {
		employeeId, err = strconv.Atoi(employeeIdStr)
		if err != nil {
			log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
			http.Error(w, "Invalid employee ID format", http.StatusBadRequest)
			return
		}
	}

	log.Printf("Fetching employee with ID: %d and lastName: %s", employeeId, lastName)
	employees, err := dbs.QueryEmployee(dbs.DB, employeeId, lastName)
	if err != nil {
		log.Printf("Error querying employee with ID %d: %v", employeeId, err)
		http.Error(w, "Failed to fetch employee", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding employee(s) to JSON: %v", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
	}
}

func AddEmployee(w http.ResponseWriter, r *http.Request) {
	var emp schema.Employee
	// err := json.NewDecoder(r.Body).Decode(&emp)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		errorMessage := fmt.Sprintf("Failed to decode request body: %v", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	// Assuming `db` is your database connection available globally or passed in some way
	employeeId, err := dbs.InsertEmployee(dbs.DB, dbs.Employees(emp))

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to add employee: %v", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	log.Printf("Employee added with ID: %d", employeeId) // Log success

	response := map[string]int{"employeeId": employeeId}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateEmployee handles the HTTP request for updating an employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeIdStr := vars["employeeId"]
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
		http.Error(w, "Invalid employee ID format", http.StatusBadRequest)
		return
	}

	var emp schema.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		log.Printf("Error decoding employee data for update: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to update employee with ID: %d", employeeId)

	// Call the database operation from the dbop package
	if err := dbs.UpdateEmployeeDB(dbs.DB, employeeId, dbs.Employees(emp)); err != nil {
		log.Printf("Error updating employee with ID %d: %v", employeeId, err)
		http.Error(w, fmt.Sprintf("Failed to update employee: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Employee with ID %d successfully updated", employeeId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp) // Optionally return the updated employee object or a success message
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeIdStr := vars["employeeId"]
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
		http.Error(w, "Invalid employee ID format", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to delete employee with ID: %d", employeeId)
	// Call a function to delete the employee by ID. This function needs to be implemented.
	err = dbs.DeleteEmployeeByID(dbs.DB, employeeId)
	if err != nil {
		log.Printf("Error deleting employee with ID %d: %v", employeeId, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Employee with ID %d successfully deleted", employeeId)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee successfully deleted"})

	// w.WriteHeader(http.StatusNoContent) // 204 No Content is often used for successful DELETE requests
}
