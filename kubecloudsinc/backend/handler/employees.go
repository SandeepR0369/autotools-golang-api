package handler

import (
	"autotools-golang-api/kubecloudsinc/backend/dbs"
	"autotools-golang-api/kubecloudsinc/backend/middleware"
	"autotools-golang-api/kubecloudsinc/backend/schema"
	"autotools-golang-api/kubecloudsinc/backend/utils"
	"errors"
	"regexp"
	"strings"
	"time"

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

	if txn != nil {
		txn.AddAttribute("httpMethod", r.Method)
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

	if txn != nil {
		txn.AddAttribute("httpMethod", r.Method)
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

	if txn != nil {
		txn.AddAttribute("httpMethod", r.Method)
	}

	var emp schema.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		errorMessage := fmt.Sprintf("Failed to decode request body: %v", err)
		log.Println(errorMessage)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidRequestBody", "AddEmployee")
		return
	}

	if err := validateAddEmployeeInput(&emp); err != nil {
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
	successMessage := fmt.Sprintf("Employee with EmployeeId: %d successfully Added", employeeId)

	response := map[string]string{"message": successMessage}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "AddEmployee")
	}
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())
	userRole, ok := r.Context().Value(middleware.RoleContextKey).(string)
	if !ok {
		utils.SendErrorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("user role not found"), "unique_error_id", "UserRoleNotFound", "UpdateEmployee")
		return
	}

	vars := mux.Vars(r)
	employeeIdStr := vars["employeeId"]
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		log.Printf("Error converting employee ID '%s' to integer: %v", employeeIdStr, err)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidEmployeeIDFormat", "UpdateEmployee")
		return
	}

	// Record a custom event before attempting to update the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("UpdateEmployeeAttempt", map[string]interface{}{
			"employeeId": employeeId,
		})
	}

	if txn != nil {
		txn.AddAttribute("httpMethod", r.Method)
	}

	var emp schema.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		log.Printf("Error decoding employee data for update: %v", err)
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "RequestBodyDecodeError", "UpdateEmployee")
		return
	}

	if err := validateUpdateEmployeeInput(&emp); err != nil {
		utils.SendErrorResponse(w, r, http.StatusBadRequest, err, "unique_error_id", "InvalidRequestBody", "AddEmployee")
		return
	}

	proceedWithUpdate := true
	var restrictedFields = make([]string, 0)

	if userRole == "editor" {
		restrictedFields = checkRestrictedFields(emp)
		if len(restrictedFields) > 0 {
			errMsg := fmt.Sprintf("You don't have enough permissions to update these fields: %s", strings.Join(restrictedFields, ", "))
			log.Println(errMsg)
			utils.SendErrorResponse(w, r, http.StatusForbidden, fmt.Errorf(errMsg), "unique_error_id", "InsufficientPermissions", "UpdateEmployee")
			proceedWithUpdate = false
		}
	} // No else if needed here, admin has full access and others are already blocked by middleware

	if proceedWithUpdate {
		if err := dbs.UpdateEmployeeDB(txn, dbs.DB, employeeId, dbs.Employees(emp)); err != nil {
			log.Printf("Error updating employee with ID %d: %v", employeeId, err)
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "EmployeeUpdateError", "UpdateEmployee")
			return
		}
		log.Printf("Employee with ID %d successfully updated", employeeId)
	}

	// Record a custom event after successfully updating the employee
	if txn != nil {
		txn.Application().RecordCustomEvent("UpdateEmployeeCompleted", map[string]interface{}{
			"employeeId": employeeId,
			"success":    true,
		})
	}

	if proceedWithUpdate {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(emp); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err, "unique_error_id", "JSONEncodingError", "UpdateEmployee")
		}
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

	if txn != nil {
		txn.AddAttribute("httpMethod", r.Method)
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

	if txn != nil {
		txn.AddAttribute("httpMethod", r.Method)
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

// validateEmployeeInput validates the fields of Employee
func validateAddEmployeeInput(emp *schema.Employee) error {
	validJobIDs := validJobIDs()
	if emp.FirstName == nil || *emp.FirstName == "" {
		return errors.New("first name is required")
	}
	if emp.LastName == nil || *emp.LastName == "" {
		return errors.New("last name is required")
	}
	if emp.Email == nil || *emp.Email == "" {
		return errors.New("email is required")
	} else if err := validateEmail(*emp.Email); err != nil {
		return err
	}
	if emp.Phone == nil || *emp.Phone == "" {
		return errors.New("phone number is required")
	} else {
		normalizedPhone, err := normalizePhoneNumber(*emp.Phone)
		if err != nil {
			return err
		}
		*emp.Phone = normalizedPhone // Update the employee's phone number to the normalized format
	}
	if emp.HireDate == nil || *emp.HireDate == "" {
		return errors.New("hire date is required")
	} else if _, err := time.Parse("2006-01-02 15:04:05", *emp.HireDate); err != nil {
		_, err := time.Parse("2006-01-02", *emp.HireDate)
		if err != nil {
			return fmt.Errorf("invalid date format for hire date: %v", err)
		}
	}
	if emp.JobId == nil || *emp.JobId == "" {
		return errors.New("job ID is required")
	} else if _, exists := validJobIDs[*emp.JobId]; !exists {
		return fmt.Errorf("invalid job ID: %s", *emp.JobId)
	}
	if emp.Salary == nil {
		return errors.New("salary is required")
	}
	return nil
}

func validateUpdateEmployeeInput(emp *schema.Employee) error {
	validJobIDs := validJobIDs()
	if emp.Email != nil && *emp.Email != "" {
		if err := validateEmail(*emp.Email); err != nil {
			return err
		}
	}

	if emp.Phone != nil && *emp.Phone != "" {
		normalizedPhone, err := normalizePhoneNumber(*emp.Phone)
		if err != nil {
			return err
		}
		*emp.Phone = normalizedPhone
	}

	if emp.HireDate != nil && *emp.HireDate != "" {
		_, err := time.Parse("2006-01-02 15:04:05", *emp.HireDate)
		if err != nil {
			_, err := time.Parse("2006-01-02", *emp.HireDate)
			if err != nil {
				return fmt.Errorf("invalid date format for hire date: %v", err)
			}
		}
	}

	if emp.JobId != nil && *emp.JobId != "" {
		if _, exists := validJobIDs[*emp.JobId]; !exists {
			return fmt.Errorf("invalid job ID: %s", *emp.JobId)
		}
	}

	return nil
}

// validateEmail checks if the email address is valid
func validateEmail(email string) error {
	// Regular expression to validate email according to the specified rules
	regexPattern := `^(?:[a-zA-Z0-9]+[_\.-]?)*[a-zA-Z0-9]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z]{2,})+$`
	match, _ := regexp.MatchString(regexPattern, email)
	if !match {
		return errors.New("invalid email format")
	}
	return nil
}

// validatePhone checks if the phone number is valid
func normalizePhoneNumber(phone string) (string, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return "", err
	}
	normalizedPhone := reg.ReplaceAllString(phone, "")

	if len(normalizedPhone) != 10 {
		return "", fmt.Errorf("phone number after normalization does not have 10 digits: %s", normalizedPhone)
	}

	reformattedPhone := fmt.Sprintf("%s.%s.%s", normalizedPhone[0:3], normalizedPhone[3:6], normalizedPhone[6:])
	return reformattedPhone, nil
}

func validJobIDs() map[string]bool {
	return map[string]bool{
		"AC_MGR":     true,
		"AC_ACCOUNT": true,
		"AD_ASST":    true,
		"AD_PRES":    true,
		"AD_VP":      true,
		"FI_ACCOUNT": true,
		"FI_MGR":     true,
		"HR_REP":     true,
		"IT_PROG":    true,
		"MK_MAN":     true,
		"MK_REP":     true,
		"PR_REP":     true,
		"PU_CLERK":   true,
		"PU_MAN":     true,
		"SA_MAN":     true,
		"SA_REP":     true,
		"SH_CLERK":   true,
		"ST_CLERK":   true,
		"ST_MAN":     true,
	}
}

func checkRestrictedFields(emp schema.Employee) []string {
	var restrictedFields []string

	if emp.Salary != nil {
		restrictedFields = append(restrictedFields, "salary")
	}
	if emp.JobId != nil {
		restrictedFields = append(restrictedFields, "jobId")
	}
	if emp.HireDate != nil {
		restrictedFields = append(restrictedFields, "hireDate")
	}
	if emp.CommissionPct != nil {
		restrictedFields = append(restrictedFields, "commissionPct")
	}
	if emp.ManagerId != nil {
		restrictedFields = append(restrictedFields, "managerId")
	}
	if emp.DepartmentId != nil {
		restrictedFields = append(restrictedFields, "departmentId")
	}

	return restrictedFields
}
