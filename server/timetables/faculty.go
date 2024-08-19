package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetFaculty() ([]models.Faculty, error) {
	var faculty []models.Faculty
	rows, err := config.Database.Query("SELECT id, name FROM faculty")
	if err != nil {
		return nil, fmt.Errorf("error querying faculty: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fac models.Faculty
		if err := rows.Scan(&fac.ID, &fac.FacultyName); err != nil {
			return nil, fmt.Errorf("error scanning faculty: %v", err)
		}
		faculty = append(faculty, fac)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return faculty, nil
}
