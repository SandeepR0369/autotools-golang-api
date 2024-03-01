// db/connection.go

package dbs

import (
	"database/sql"
	"fmt"
	schema "kubecloudsinc/Schema"

	_ "github.com/godror/godror"
//	"github.com/alexbrainman/odbc"
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

func PrintEmployees(employees []Employees) {
	for _, emp := range employees {
		fmt.Printf("ID: %d, Name: %s %s, Email: %s, Phone: %s, Hire Date: %s, Job ID: %s, Salary: %.2f, Commission Pct: %v, Manager ID: %d, Department ID: %d\n",
			emp.EmployeeId, emp.FirstName, emp.LastName, emp.Email, emp.Phone, emp.HireDate, emp.JobId, emp.Salary, emp.CommissionPct, emp.ManagerId, emp.DepartmentId)
	}
}
