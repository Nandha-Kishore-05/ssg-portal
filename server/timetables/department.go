package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetDepartments() ([]models.Department, error) {
	var dept []models.Department
	rows, err := config.Database.Query("SELECT id, name FROM departments")
	if err != nil {
		return nil, fmt.Errorf("error querying faculty: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dep models.Department
		if err := rows.Scan(&dep.ID, &dep.Department); err != nil {
			return nil, fmt.Errorf("error scanning faculty: %v", err)
		}
		dept = append(dept, dep)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return dept, nil
}
