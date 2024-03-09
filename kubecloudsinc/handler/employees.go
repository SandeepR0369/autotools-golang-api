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
	"github.com/newrelic/go-agent/v3/newrelic"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	// Record a custom event before querying all employees
	if txn != nil {
		txn.Application().RecordCustomEvent("GetEmployeesAttempt", map[string]interface{}{})
	}

	employees, err := dbs.QueryEmployees(txn, dbs.DB)
	if err != nil {
		log.Printf("Error querying all employees: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "hagsv123", "NoMatchingRecordFound", "Employee Retrieval")
		return
	}

	// Record a custom event after successfully querying all employees
	if txn != nil {
		txn.Application().RecordCustomEvent("GetEmployeesCompleted", map[string]interface{}{
			"count": len(employees),
		})
	}

	// Instead of printing, send the employees back as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding employees to JSON: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "hagsv123", "JSONEncodingError", "Employee Encoding")
	}
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	queryValues := r.URL.Query()
	employeeIdStr := queryValues.Get("employeeId")
	lastName := queryValues.Get("lastName")
	// Record a custom event before querying the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("GetEmployeeAttempt", map[string]interface{}{
			"employeeId": employeeIdStr,
			"lastName":   lastName,
		})
	}
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
	employees, err := dbs.QueryEmployee(txn, dbs.DB, employeeId, lastName)
	if err != nil {
		if err.Error() == "no record was found with provided identifiers" {
			utils.SendErrorResponse(w, r, http.StatusNotFound, err, "unique_error_id", "NoMatchingRecordFound", "QueryEmployee")
		} else {
			// Handle other errors
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "QueryError", "QueryEmployee")
		}
		return
	}

	// Record a custom event after successfully querying the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("GetEmployeeCompleted", map[string]interface{}{
			"count": len(employees),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding employee(s) to JSON: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "GetEmployee")
	}
}

func AddEmployee(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())
	// Record a custom event before querying the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("AddEmployeeAttempt", map[string]interface{}{})
	}
	var emp schema.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		errorMessage := fmt.Sprintf("Failed to decode request body: %v", err)
		log.Println(errorMessage)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidRequestBody", "AddEmployee")
		return
	}

	// Assuming `db` is your database connection available globally or passed in some way
	employeeId, err := dbs.InsertEmployee(txn, dbs.DB, dbs.Employees(emp))

	// Record a custom event after successfully querying all employees
	if txn != nil {
		txn.Application().RecordCustomEvent("AddEmployeeCompleted", map[string]interface{}{
			"employeeAdded": employeeId,
		})
	}

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
	txn := newrelic.FromContext(r.Context())

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

	// Record a custom event before attempting to update the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("UpdateEmployeeAttempt", map[string]interface{}{
			"employeeId": employeeId,
		})
	}

	log.Printf("Attempting to update employee with ID: %d", employeeId)

	// Call the database operation from the dbop package
	if err := dbs.UpdateEmployeeDB(txn, dbs.DB, employeeId, dbs.Employees(emp)); err != nil {
		log.Printf("Error updating employee with ID %d: %v", employeeId, err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "EmployeeUpdateError", "UpdateEmployee")
		return
	}

	// Record a custom event after successfully updating the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("UpdateEmployeeCompleted", map[string]interface{}{
			"employeeId": employeeId,
			"success":    true,
		})
	}

	log.Printf("Employee with ID %d successfully updated", employeeId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(emp); err != nil {
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "UpdateEmployee")
	}
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())
	vars := mux.Vars(r)
	employeeIdStr := vars["employeeId"]
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidEmployeeIDFormat", "DeleteEmployee")
		return
	}

	// Record a custom event before attempting to delete the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("DeleteEmployeeAttempt", map[string]interface{}{
			"employeeId": employeeId,
		})
	}

	log.Printf("Attempting to delete employee with ID: %d", employeeId)

	// Call a function to delete the employee by ID. This function needs to be implemented.
	err = dbs.DeleteEmployeeByID(txn, dbs.DB, employeeId)
	if err != nil {
		log.Printf("Error deleting employee with ID %d: %v", employeeId, err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "EmployeeDeletionError", "DeleteEmployee")
		return
	}

	// Record a custom event after successfully deleting the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("DeleteEmployeeCompleted", map[string]interface{}{
			"employeeId": employeeId,
			"success":    true,
		})
	}

	log.Printf("Employee with ID %d successfully deleted", employeeId)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Employee successfully deleted"}); err != nil {
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "DeleteEmployee")
	}
}

func GetEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())
	vars := mux.Vars(r)
	employeeIdStr := vars["employeeId"]
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidEmployeeIDFormat", "GetEmployeeProfile")
		return
	}

	// Record a custom event before attempting to get the employee profile
	if txn != nil {
		txn.Application().RecordCustomEvent("GetEmployeeProfileAttempt", map[string]interface{}{
			"employeeId": employeeId,
		})
	}

	log.Printf("Attempting to get employee profile with ID: %d", employeeId)

	// Query database to get employee profile
	employeeProfile, err := dbs.GetEmployeeProfile(txn, dbs.DB, employeeId)
	if err != nil {
		if err.Error() == "no record was found with provided identifiers" {
			utils.SendErrorResponse(w, r, http.StatusNotFound, err, "unique_error_id", "NoMatchingRecordFound", "GetEmployeeProfile")
		} else {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "QueryError", "GetEmployeeProfile")
		}
		return
	}

	// Record a custom event after successfully retrieved the employee profile
	if txn != nil {
		txn.Application().RecordCustomEvent("GetEmployeeProfileCompleted", map[string]interface{}{
			"employeeId": employeeId,
			"success":    true,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employeeProfile); err != nil {
		log.Printf("Error encoding employee profile to JSON: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "GetEmployeeProfile")
		return
	}
}
