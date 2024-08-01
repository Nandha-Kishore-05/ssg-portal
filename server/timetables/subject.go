package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetSubjects(departmentID int) ([]models.Subject, error) {
	var subjects []models.Subject

	// Define the query to select subjects by department ID
	query := `
		SELECT id, name
		FROM subjects
		WHERE department_id = ?
	`

	// Execute the query
	rows, err := config.Database.Query(query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("error querying subjects: %v", err)
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			return nil, fmt.Errorf("error scanning subject: %v", err)
		}
		subjects = append(subjects, subject)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return subjects, nil
}
