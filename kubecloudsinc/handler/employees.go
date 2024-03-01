package handler

import (
	"net/http"
	"kubecloudsinc/dbs"
	//"kubecloudsinc/Schema"
	"encoding/json"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := dbs.QueryEmployees(dbs.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//dbs.PrintEmployees(employees)
	// Instead of printing, send the employees back as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(employees)
}