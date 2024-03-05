package dbs

import (
	schema "autotools-golang-api/kubecloudsinc/schema"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/godror/godror"
)

var (
	// Database connection pool
	DB *sql.DB
)

type Employees schema.Employee

// InitDB initializes the database connection using the provided DSN
func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("godror", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Check the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

func QueryEmployees(db *sql.DB) ([]Employees, error) {
	log.Println("Making a DB call to get all employees")
	// Define the SQL query
	query := `SELECT employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, commission_pct, manager_id, department_id FROM employees`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows and scan the results into the Employee struct
	var employees []Employees
	for rows.Next() {
		var emp Employees
		// Assuming commission_pct can be null, it's handled as sql.NullFloat64
		var commissionPct sql.NullFloat64
		err := rows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Phone, &emp.HireDate, &emp.JobId, &emp.Salary, &commissionPct, &emp.ManagerId, &emp.DepartmentId)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		employees = append(employees, emp)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return employees, nil
}

func QueryEmployee(db *sql.DB, employeeId int, lastName string) ([]Employees, error) {
	log.Println("Making a DB call to get employee")

	// First, check if the employee exists
	err := checkEmployeeExistence(db, employeeId, lastName)
	if err != nil {
		return nil, err
	}

	// Initialize the base query without the WHERE clause
	baseQuery := `SELECT employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, commission_pct, manager_id, department_id FROM employees`

	var queryParams []interface{}
	var conditions []string

	// Append conditions based on provided arguments
	if employeeId > 0 {
		conditions = append(conditions, "employee_id = :1")
		queryParams = append(queryParams, employeeId)
	}
	if lastName != "" {
		conditions = append(conditions, "last_name LIKE ?")
		queryParams = append(queryParams, "%"+lastName+"%")
	}

	// If there are conditions, append them to the baseQuery
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Debugging: Log the final query and parameters
	log.Println("Final Query:", baseQuery)
	log.Println("Parameters:", queryParams)

	// This is purely for debugging; do not use this to execute the query.
	debugQuery := DebugQuery(baseQuery, queryParams)
	log.Println("Debug Query:", debugQuery)

	// Execute the query using the built query string and parameters
	rows, err := db.Query(baseQuery, queryParams...)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var employees []Employees
	for rows.Next() {
		var emp Employees
		err := rows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Phone, &emp.HireDate, &emp.JobId, &emp.Salary, &emp.CommissionPct, &emp.ManagerId, &emp.DepartmentId)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return employees, nil
}

func InsertEmployee(db *sql.DB, emp Employees) (int, error) {
	log.Printf("Making a DB call to insert employee")
	query := `INSERT INTO employees (employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, commission_pct, manager_id, department_id)
              VALUES (:1, :2, :3, :4, :5, TO_DATE(:6, 'YYYY-MM-DD HH24:MI:SS'), :7, :8, :9, :10, :11) RETURNING employee_id INTO :12`

	// Preparing a variable to hold the returned employee_id
	var returnedEmployeeId int

	args := []interface{}{
		emp.EmployeeId,
		emp.FirstName,
		emp.LastName,
		emp.Email,
		emp.Phone,
		emp.HireDate, // Passed directly into TO_TIMESTAMP
		emp.JobId,
		emp.Salary,
		emp.CommissionPct, // Assuming the driver can handle nil directly; if not, use sql.NullFloat64
		emp.ManagerId,     // Assuming the driver can handle nil directly; if not, use sql.NullInt32
		emp.DepartmentId,
		sql.Out{Dest: &returnedEmployeeId}, // For capturing the RETURNING value
	}

	if _, err := db.Exec(query, args...); err != nil {
		log.Printf("Failed to insert employee: %v", err)
		return 0, fmt.Errorf("failed to insert employee: %v", err)
	}

	if emp.EmployeeId != nil {
		return *emp.EmployeeId, nil
	}

	return 0, errors.New("failed to insert employee")
}

func UpdateEmployeeDB(db *sql.DB, employeeId int, emp Employees) error {
	log.Printf("Making a DB call to update employeeId: %d", employeeId)
	// First, check if the employee exists
	err := checkEmployeeExistence(db, employeeId, "")
	if err != nil {
		return err
	}

	// Initialize the base query and argument counter
	query := "UPDATE employees SET "
	var args []interface{}
	var updates []string
	var argCount int = 1

	// Dynamically build the updates list and args based on non-nil fields in emp
	addUpdate := func(field interface{}, fieldName string) {
		// Always add the field name to the updates list
		updates = append(updates, fmt.Sprintf("%s = :%d", fieldName, argCount))

		// Add the field value to args slice
		args = append(args, field)
		argCount++
	}

	// Check each field and append to updates and args if not nil
	addUpdate(emp.FirstName, "first_name")
	addUpdate(emp.LastName, "last_name")
	addUpdate(emp.Email, "email")
	addUpdate(emp.Phone, "phone_number")
	if emp.HireDate != nil && *emp.HireDate != "" {
		layout := "2006-01-02 15:04:05"
		hireDate, err := time.Parse(layout, *emp.HireDate)
		if err != nil {
			log.Printf("Error parsing hire date '%s': %v", *emp.HireDate, err)
			return fmt.Errorf("error parsing hire date: %v", err)
		}
		addUpdate(hireDate, "hire_date")
	}
	addUpdate(emp.JobId, "job_id")
	addUpdate(emp.Salary, "salary")
	addUpdate(emp.CommissionPct, "commission_pct")
	addUpdate(emp.ManagerId, "manager_id")
	addUpdate(emp.DepartmentId, "department_id")

	// If no fields were updated, return an error
	if len(updates) == 0 {
		return errors.New("no fields provided for update")
	}

	// Finalize the query by appending the WHERE clause
	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE employee_id = :%d", argCount)
	args = append(args, employeeId)
	log.Println("JAFFA1")

	// Debugging: Log the final query and parameters
	debugQuery := DebugQuery(query, args)
	log.Println("Debug Query Update:", debugQuery)
	// Execute the update
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("Failed to update employee: %v", err)
		return fmt.Errorf("failed to update employee: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		log.Printf("No rows updated, employee with ID %d may not exist or no new data provided", employeeId)
		// Instead of returning an error, you can return a specific message indicating success but no update was needed.
		return fmt.Errorf("no update needed or employee with ID %d not found", employeeId)
	}

	log.Printf("Employee with ID %d updated successfully, %d rows affected", employeeId, rowsAffected)
	return nil
}

func DeleteEmployeeByID(db *sql.DB, employeeId int) error {
	// First, check if the employee exists
	err := checkEmployeeExistence(db, employeeId, "")
	if err != nil {
		return err
	}

	log.Printf("Attempting to delete employee with ID: %d", employeeId)
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	_, err = tx.Exec("DELETE FROM employees WHERE employee_id = :1", employeeId)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Printf("Failed to rollback transaction: %v", rbErr)
		}
		log.Printf("Failed to delete employee: %v", err)
		return fmt.Errorf("failed to delete employee: %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Employee with ID %d deleted successfully", employeeId)
	return nil
}

func PrintEmployees(employees []Employees) {
	for _, emp := range employees {
		fmt.Printf("ID: %d, Name: %s %s, Email: %s, Phone: %s, Hire Date: %s, Job ID: %s, Salary: %.2f, Commission Pct: %v, Manager ID: %d, Department ID: %d\n",
			emp.EmployeeId, emp.FirstName, emp.LastName, emp.Email, emp.Phone, emp.HireDate, emp.JobId, emp.Salary, emp.CommissionPct, emp.ManagerId, emp.DepartmentId)
	}
}

func DebugQuery(query string, params []interface{}) string {
	var buffer bytes.Buffer
	n := 0
	paramIndex := 1 // Starting index for named parameters
	for _, p := range params {
		namedParam := fmt.Sprintf(":%d", paramIndex)
		pos := strings.Index(query[n:], namedParam)
		if pos == -1 {
			break
		}
		buffer.WriteString(query[n : n+pos])
		buffer.WriteString(fmt.Sprintf("'%v'", p))
		n += pos + len(namedParam)
		paramIndex++
	}
	buffer.WriteString(query[n:])
	return buffer.String()
}

// func checkEmployeeExistence(db *sql.DB, employeeId int) error {
// 	var count int
// 	err := db.QueryRow("SELECT COUNT(employee_id) FROM employees WHERE employee_id = :1", employeeId).Scan(&count)
// 	if err != nil {
// 		log.Printf("Error checking employee existence: %v", err)
// 		return fmt.Errorf("error checking employee existence: %v", err)
// 	}
// 	if count == 0 {
// 		return fmt.Errorf("employee with ID %d not found", employeeId)
// 	}
// 	return nil
// }

func checkEmployeeExistence(db *sql.DB, employeeId int, lastName string) error {
	// Initialize the SQL query string and parameters slice
	query := "SELECT COUNT(employee_id) FROM employees WHERE 1=1"
	var params []interface{}

	// Add conditions based on provided arguments
	if employeeId > 0 {
		query += " AND employee_id = :1"
		params = append(params, employeeId)
	}
	if lastName != "" {
		query += " AND last_name = :2"
		params = append(params, lastName)
	}

	// Execute the query
	var count int
	err := db.QueryRow(query, params...).Scan(&count)
	if err != nil {
		log.Printf("Error checking employee existence: %v", err)
		return fmt.Errorf("error checking employee existence: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}
