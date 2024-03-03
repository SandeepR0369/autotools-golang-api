package schema

type Employee struct {
	EmployeeId    *int     `json:"employeeId"`
	FirstName     *string  `json:"firstName"`
	LastName      *string  `json:"lastName"`
	Email         *string  `json:"email"`
	Phone         *string  `json:"phone"`
	HireDate      *string  `json:"hireDate"`
	JobId         *string  `json:"jobId"`
	Salary        *float64 `json:"salary"`
	CommissionPct *float64 `json:"commissionPct"`
	ManagerId     *int     `json:"managerId"`
	DepartmentId  *int     `json:"departmentId"`
}
