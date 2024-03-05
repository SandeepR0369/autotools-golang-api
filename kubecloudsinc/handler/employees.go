package handler

import (
	"autotools-golang-api/kubecloudsinc/dbs"
	"autotools-golang-api/kubecloudsinc/schema"
	"autotools-golang-api/kubecloudsinc/utils"

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
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "hagsv123", "NoMatchingRecordFound", "Employee Retrieval")
		return
	}

	// Instead of printing, send the employees back as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding employees to JSON: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "hagsv123", "JSONEncodingError", "Employee Encoding")
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
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidEmployeeID", "GetEmployee")
			return
		}
	}

	log.Printf("Fetching employee with ID: %d and lastName: %s", employeeId, lastName)
	employees, err := dbs.QueryEmployee(dbs.DB, employeeId, lastName)
	if err != nil {
		if err.Error() == "no record was found with provided identifiers" {
			utils.SendErrorResponse(w, r, http.StatusNotFound, err, "unique_error_id", "NoMatchingRecordFound", "QueryEmployee")
		} else {
			// Handle other errors
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "QueryError", "QueryEmployee")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding employee(s) to JSON: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "GetEmployee")
	}
}

func AddEmployee(w http.ResponseWriter, r *http.Request) {
	var emp schema.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		errorMessage := fmt.Sprintf("Failed to decode request body: %v", err)
		log.Println(errorMessage)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidRequestBody", "AddEmployee")
		return
	}

	// Assuming `db` is your database connection available globally or passed in some way
	employeeId, err := dbs.InsertEmployee(dbs.DB, dbs.Employees(emp))

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to add employee: %v", err)
		log.Println(errorMessage)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "EmployeeInsertionError", "AddEmployee")
		return
	}

	log.Printf("Employee added with ID: %d", employeeId) // Log success

	response := map[string]int{"employeeId": employeeId}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "AddEmployee")
	}
}

// UpdateEmployee handles the HTTP request for updating an employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeIdStr := vars["employeeId"]
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidEmployeeIDFormat", "UpdateEmployee")
		return
	}

	var emp schema.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		log.Printf("Error decoding employee data for update: %v", err)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "RequestBodyDecodeError", "UpdateEmployee")
		return
	}

	log.Printf("Attempting to update employee with ID: %d", employeeId)

	// Call the database operation from the dbop package
	if err := dbs.UpdateEmployeeDB(dbs.DB, employeeId, dbs.Employees(emp)); err != nil {
		log.Printf("Error updating employee with ID %d: %v", employeeId, err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "EmployeeUpdateError", "UpdateEmployee")
		return
	}

	log.Printf("Employee with ID %d successfully updated", employeeId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(emp); err != nil {
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "UpdateEmployee")
	}
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeIdStr := vars["employeeId"]
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidEmployeeIDFormat", "DeleteEmployee")
		return
	}

	log.Printf("Attempting to delete employee with ID: %d", employeeId)
	// Call a function to delete the employee by ID. This function needs to be implemented.
	err = dbs.DeleteEmployeeByID(dbs.DB, employeeId)
	if err != nil {
		log.Printf("Error deleting employee with ID %d: %v", employeeId, err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "EmployeeDeletionError", "DeleteEmployee")
		return
	}

	log.Printf("Employee with ID %d successfully deleted", employeeId)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Employee successfully deleted"}); err != nil {
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "DeleteEmployee")
	}
}
