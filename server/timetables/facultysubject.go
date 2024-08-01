package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetFacultySubjects() ([]models.FacultySubject, error) {
	var facultySubjects []models.FacultySubject
	rows, err := config.Database.Query("SELECT faculty_id, subject_id FROM faculty_subjects")
	if err != nil {
		return nil, fmt.Errorf("error querying faculty-subject mappings: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fs models.FacultySubject
		if err := rows.Scan(&fs.FacultyID, &fs.SubjectID); err != nil {
			return nil, fmt.Errorf("error scanning faculty-subject: %v", err)
		}
		facultySubjects = append(facultySubjects, fs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return facultySubjects, nil
}
