// db/connection.go

package dbs

import (
	schema "autotools-golang-api/kubecloudsinc/schema"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/godror/godror"
	// "github.com/alexbrainman/odbc"
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
	//DB, _ := odbc.Connect("ODBC",dsn)
	if err != nil {
		return err
	}

	// Check the connection
	err = DB.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Database connection established")
	return nil
}

func QueryEmployees(db *sql.DB) ([]Employees, error) {
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
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		// Handle the nullable commission_pct
		// if commissionPct.Valid {
		// 	emp.CommissionPct = commissionPct.Float64
		// } else {
		// 	emp.CommissionPct = nil
		// }
		employees = append(employees, emp)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return employees, nil
}

func QueryEmployee(db *sql.DB, employeeId int, lastName string) ([]Employees, error) {
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
		// Determine the correct placeholder based on existing queryParams
		placeholder := fmt.Sprintf("last_name LIKE :%d", len(queryParams)+1)
		conditions = append(conditions, placeholder)
		queryParams = append(queryParams, "%"+lastName+"%")
	}

	// If there are conditions, append them to the baseQuery
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	fmt.Println("Executing query:", baseQuery)
	fmt.Println("With parameters:", queryParams)

	// Execute the query using the built query string and parameters
	rows, err := db.Query(baseQuery, queryParams...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var employees []Employees
	for rows.Next() {
		var emp Employees
		var commissionPct sql.NullFloat64
		err := rows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Phone, &emp.HireDate, &emp.JobId, &emp.Salary, &commissionPct, &emp.ManagerId, &emp.DepartmentId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return employees, nil
}

func InsertEmployee(db *sql.DB, emp Employees) (int, error) {

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

	_, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert employee: %v", err)
	}

	if emp.EmployeeId != nil {
		return *emp.EmployeeId, nil
	}

	return 0, errors.New("employee ID is nil")
}

func UpdateEmployeeDB(db *sql.DB, employeeId int, emp Employees) error {
	// Initialize the base query
	query := "UPDATE employees SET "
	var args []interface{}
	var updates []string
	var argCount int = 1

	// Dynamically add fields to be updated
	if emp.FirstName != nil {
		updates = append(updates, fmt.Sprintf("first_name = :%d", argCount))
		args = append(args, emp.FirstName)
		argCount++
	}
	if emp.LastName != nil {
		updates = append(updates, fmt.Sprintf("last_name = :%d", argCount))
		args = append(args, emp.LastName)
		argCount++
	}
	// Repeat for each field...

	if emp.Email != nil {
		updates = append(updates, fmt.Sprintf("email = :%d", argCount))
		args = append(args, emp.Email)
		argCount++
	}

	// Continue for other fields...

	// Check if there are updates to make
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	// Finalize the query
	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE employee_id = :%d", argCount)
	args = append(args, employeeId)

	// Execute the update
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update employee: %v", err)
	}

	return nil
}

func DeleteEmployeeByID(db *sql.DB, employeeId int) error {
	// Implement the database call to delete the employee record.
	// Using named placeholder syntax for Oracle
	_, err := db.Exec("DELETE FROM employees WHERE employee_id = :1", employeeId)
	if err != nil {
		return err
	}
	return nil
}

func PrintEmployees(employees []Employees) {
	for _, emp := range employees {
		fmt.Printf("ID: %d, Name: %s %s, Email: %s, Phone: %s, Hire Date: %s, Job ID: %s, Salary: %.2f, Commission Pct: %v, Manager ID: %d, Department ID: %d\n",
			emp.EmployeeId, emp.FirstName, emp.LastName, emp.Email, emp.Phone, emp.HireDate, emp.JobId, emp.Salary, emp.CommissionPct, emp.ManagerId, emp.DepartmentId)
	}
}
