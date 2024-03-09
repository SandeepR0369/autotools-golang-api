package dbs

import (
	schema "autotools-golang-api/kubecloudsinc/schema"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/godror/godror"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	// Database connection pool
	DB *sql.DB
)

type Employees schema.Employee
type EmployeeProfile schema.EmployeeProfile

type NullableFloat64 struct {
	Float64 float64
	Valid   bool
}

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

func QueryEmployees(txn *newrelic.Transaction, db *sql.DB) ([]Employees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	segment := newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL, // Adjust for your database
		Collection: "employees",
		Operation:  "SELECT",
	}
	defer segment.End()

	log.Println("Making a DB call to get all employees")
	// Define the SQL query
	query := `SELECT employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, commission_pct, manager_id, department_id FROM employees`

	// Execute the query
	rows, err := db.QueryContext(ctx, query)
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

func QueryEmployee(txn *newrelic.Transaction, db *sql.DB, employeeId int, lastName string) ([]Employees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	segment := newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL, // Adjust for your database
		Collection: "employees",
		Operation:  "SELECT",
	}
	defer segment.End()

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
		conditions = append(conditions, "last_name LIKE :2")
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
	rows, err := db.QueryContext(ctx, baseQuery, queryParams...)
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

func InsertEmployee(txn *newrelic.Transaction, db *sql.DB, emp Employees) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	segment := newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL, // Adjust based on your actual database
		Collection: "employees",
		Operation:  "INSERT",
	}
	defer segment.End()

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

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		log.Printf("Failed to insert employee: %v", err)
		return 0, fmt.Errorf("failed to insert employee: %v", err)
	}

	if emp.EmployeeId != nil {
		return *emp.EmployeeId, nil
	}

	return 0, errors.New("failed to insert employee")
}

func UpdateEmployeeDB(txn *newrelic.Transaction, db *sql.DB, employeeId int, emp Employees) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	segment := newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL, // Adjust for your database
		Collection: "employees",
		Operation:  "UPDATE",
	}
	defer segment.End()

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

	// Debugging: Log the final query and parameters
	debugQuery := DebugQuery(query, args)
	log.Println("Debug Query Update:", debugQuery)
	// Execute the update
	result, err := db.ExecContext(ctx, query, args...)
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

func DeleteEmployeeByID(txn *newrelic.Transaction, db *sql.DB, employeeId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	segment := newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL,
		Collection: "employees",
		Operation:  "DELETE",
	}
	defer segment.End()

	// First, check if the employee exists
	err := checkEmployeeExistence(db, employeeId, "")
	if err != nil {
		return err
	}

	log.Printf("Making a DB call to delete employee with ID: %d", employeeId)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM employees WHERE employee_id = :1", employeeId)
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

func GetEmployeeProfile(txn *newrelic.Transaction, db *sql.DB, employeeId int) (*EmployeeProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	segment := newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(txn),
		Product:    newrelic.DatastoreMySQL, // Adjust for your database
		Collection: "employees",
		Operation:  "SELECT",
	}
	defer segment.End()

	log.Printf("Making a DB call to fetch employee profile with ID: %d", employeeId)

	// First, check if the employee exists
	err := checkEmployeeExistence(db, employeeId, "")
	if err != nil {
		return nil, err
	}

	query := `
	SELECT
    e.employee_id,
    e.first_name,
    e.last_name,
    e.email,
    e.phone_number,
    e.salary,
    e.commission_pct,
    e.job_id,
    j.job_title,
    e.hire_date,
    d.department_id,
    d.department_name,
    e.manager_id,
    m.first_name AS manager_firstName,
    m.last_name AS manager_lastName,
    l.location_id,
    l.street_address,
    l.postal_code,
    l.city,
    l.state_province,
    c.country_id,
    c.country_name,
    r.region_id,
    r.region_name,
    jh.job_title AS job_history_title,
    jh.start_date,
    jh.end_date,
    jh.job_id AS job_history_id
FROM
    EMPLOYEES e
JOIN
    EMPLOYEES m ON e.manager_id = m.employee_id
JOIN
    JOBS j ON j.job_id = e.job_id
JOIN
    DEPARTMENTS d ON d.department_id = e.department_id
JOIN
    LOCATIONS l ON d.location_id = l.location_id
JOIN
    COUNTRIES c ON c.country_id = l.country_id
JOIN
    REGIONS r ON c.region_id = r.region_id
LEFT JOIN (
    SELECT
        jh.employee_id,
        j.job_title,
        jh.start_date,
        jh.end_date,
        jh.job_id
    FROM
        JOB_HISTORY jh
    JOIN
        JOBS j ON j.job_id = jh.job_id
) jh ON jh.employee_id = e.employee_id
WHERE
    e.employee_id = :1
    `
	rows, err := db.QueryContext(ctx, query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employeeProfile := &EmployeeProfile{
		JobDetails: &schema.JobDetails{},
	}

	for rows.Next() {
		var (
			managerFirstName, managerLastName, departmentName, countryName, regionName, jobHistoryTitle string
			hireDate, jobTitle                                                                          string
			startDate, endDate                                                                          sql.NullString
			departmentId, locationId, regionId                                                          int
			streetAddress, postalCode, city, stateProvince, countryId, jobId, jobHistoryId              string
			commissionPct                                                                               sql.NullFloat64
			salary                                                                                      sql.NullFloat64 // Change to sql.NullFloat64
		)

		err := rows.Scan(
			&employeeProfile.EmployeeId,
			&employeeProfile.FirstName,
			&employeeProfile.LastName,
			&employeeProfile.Email,
			&employeeProfile.Phone,
			&salary,
			&commissionPct,
			&jobId,
			&jobTitle,
			&hireDate,
			&departmentId,
			&departmentName,
			&employeeProfile.ManagerId,
			&managerFirstName,
			&managerLastName,
			&locationId,
			&streetAddress,
			&postalCode,
			&city,
			&stateProvince,
			&countryId,
			&countryName,
			&regionId,
			&regionName,
			&jobHistoryTitle,
			&startDate,
			&endDate,
			&jobHistoryId,
		)
		if err != nil {
			return nil, err
		}

		// Set the commission percentage
		employeeProfile.CommissionPct = &commissionPct.Float64

		// Set the commission percentage
		var commission NullableFloat64
		if commissionPct.Valid {
			commission.Float64 = commissionPct.Float64
			commission.Valid = true
		}

		var startDatePtr *string
		if startDate.Valid {
			startDateValue := startDate.String
			startDatePtr = &startDateValue
		}

		var endDatePtr *string
		if endDate.Valid {
			endDateValue := endDate.String
			endDatePtr = &endDateValue
		}

		// Set the salary
		if salary.Valid {
			employeeProfile.Salary = &salary.Float64
		}

		// Create a new job
		job := &schema.Job{
			JobId:          &jobId,
			JobTitle:       &jobTitle,
			HireDate:       &hireDate,
			DepartmentName: &departmentName,
			Salary:         employeeProfile.Salary, // Use the Salary from EmployeeProfile
			DepartmentId:   &departmentId,
		}

		// Check if the job exists in the jobDetails.Jobs slice
		var existingJob *schema.Job
		for _, j := range employeeProfile.JobDetails.Jobs {
			if *j.JobId == jobId {
				existingJob = j
				break
			}
		}

		// If the job doesn't exist, append it to the jobDetails.Jobs slice
		if existingJob == nil {
			employeeProfile.JobDetails.Jobs = append(employeeProfile.JobDetails.Jobs, job)
			existingJob = job
		}

		// Append job history to the existing job
		existingJob.JobHistory = append(existingJob.JobHistory, &schema.JobHistory{
			JobTitle:  &jobHistoryTitle,
			StartDate: startDatePtr,
			EndDate:   endDatePtr,
			JobId:     &jobHistoryId,
		})

		// Set manager details
		employeeProfile.JobDetails.Manager = &schema.Manager{
			ManagerId:    employeeProfile.ManagerId,
			ManagerFirst: &managerFirstName,
			ManagerLast:  &managerLastName,
		}

		// Set department details
		employeeProfile.JobDetails.Department = &schema.Department{
			DepartmentId:   &departmentId,
			DepartmentName: &departmentName,
			Location: &schema.Location{
				LocationId:    &locationId,
				StreetAddress: &streetAddress,
				PostalCode:    &postalCode,
				City:          &city,
				StateProvince: &stateProvince,
				Country: &schema.Country{
					CountryId:   &countryId,
					CountryName: &countryName,
					Region: &schema.Region{
						RegionId:   &regionId,
						RegionName: &regionName,
					},
				},
			},
		}
	}

	return employeeProfile, nil
}

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
