package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)


func GetSubjects(departmentID, semesterID int) ([]models.Subject, error) {
	var subjects []models.Subject

	
	query := `
		SELECT id, name, status, periods
		FROM subjects
		WHERE department_id = ? AND semester_id = ?
	`


	rows, err := config.Database.Query(query, departmentID, semesterID)
	if err != nil {
		return nil, fmt.Errorf("error querying subjects: %v", err)
	}
	defer rows.Close()


	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Status, &subject.Period); err != nil {
			return nil, fmt.Errorf("error scanning subject: %v", err)
		}
		subjects = append(subjects, subject)
	}

	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return subjects, nil
}
