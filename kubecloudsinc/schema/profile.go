package schema

import "database/sql"

type EmployeeProfile struct {
	EmployeeId    *int             `json:"employeeId"`
	FirstName     *string          `json:"firstName"`
	LastName      *string          `json:"lastName"`
	Email         *string          `json:"email"`
	Phone         *string          `json:"phone"`
	Salary        *float64         `json:"salary"`
	CommissionPct *sql.NullFloat64 `json:"commissionPct,omitempty"`
	ManagerId     *int             `json:"managerId"`
	JobDetails    *JobDetails      `json:"job_details,omitempty"`
}

type JobDetails struct {
	Jobs       []*Job      `json:"jobs"`
	Manager    *Manager    `json:"manager"`
	Department *Department `json:"department"`
}

type Job struct {
	JobId          *string       `json:"jobId"`
	JobTitle       *string       `json:"jobTitle"`
	HireDate       *string       `json:"hireDate"`
	Salary         *float64      `json:"salary"`
	DepartmentId   *int          `json:"departmentId"`
	DepartmentName *string       `json:"departmentName"`
	JobHistory     []*JobHistory `json:"job_history"`
}

type JobHistory struct {
	JobId     *string `json:"jobId"`
	JobTitle  *string `json:"jobTitle"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
}

type Manager struct {
	ManagerId    *int    `json:"managerId"`
	ManagerFirst *string `json:"managerFirst"`
	ManagerLast  *string `json:"managerLast"`
}

type Department struct {
	DepartmentId   *int      `json:"departmentId"`
	DepartmentName *string   `json:"departmentName"`
	Location       *Location `json:"location"`
}

type Location struct {
	LocationId    *int     `json:"locationId"`
	StreetAddress *string  `json:"streetAddress"`
	PostalCode    *string  `json:"postalCode"`
	City          *string  `json:"city"`
	StateProvince *string  `json:"stateProvince"`
	Country       *Country `json:"country"`
}

type Country struct {
	CountryId   *string `json:"countryId"`
	CountryName *string `json:"countryName"`
	Region      *Region `json:"region"`
}

type Region struct {
	RegionId   *int    `json:"regionId"`
	RegionName *string `json:"regionName"`
}
