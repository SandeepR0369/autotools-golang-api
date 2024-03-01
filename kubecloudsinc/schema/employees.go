package schema

import (
   // "database/sql"
    // Include other imports as needed
)
// type Employee struct {
// 	EmployeeID    int     `json:"employeeId"`
// 	FirstName     string  `json:"firstName"`
// 	LastName      string  `json:"lastName"`
// 	Email         string  `json:"email"`
// 	Phone         string  `json:"phone"`
// 	HireDate      string  `json:"hireDate"`
// 	JobID         string  `json:"jobId"`
// 	Salary        float64 `json:"salary"`
// 	CommissionPct any     `json:"commissionPct"`
// 	ManagerID     sql.NullInt64     `json:"managerId"`
// 	DepartmentID  sql.NullInt64     `json:"departmentId"`
// }

type Employee struct {
    EmployeeId   int        `json:"employeeId"`
    FirstName    string     `json:"firstName"`
    LastName     string     `json:"lastName"`
    Email        string     `json:"email"`
    Phone        string     `json:"phone"`
    HireDate     string  `json:"hireDate"`
    JobId        string     `json:"jobId"`
    Salary       float64    `json:"salary"`
    CommissionPct *float64  `json:"commissionPct"` // Use pointer for nullable field
    ManagerId    *int       `json:"managerId"`    // Use pointer for nullable field
    DepartmentId *int        `json:"departmentId"`
}