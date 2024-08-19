package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetFacultySubjects(semesterID int) ([]models.FacultySubject, error) {
	var facultySubjects []models.FacultySubject
	rows, err := config.Database.Query("SELECT faculty_id, subject_id,semester_id FROM faculty_subjects WHERE semester_id = ? ", semesterID)
	if err != nil {
		return nil, fmt.Errorf("error querying faculty-subject mappings: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fs models.FacultySubject
		if err := rows.Scan(&fs.FacultyID, &fs.SubjectID, &fs.SemesterID); err != nil {
			return nil, fmt.Errorf("error scanning faculty-subject: %v", err)
		}
		facultySubjects = append(facultySubjects, fs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return facultySubjects, nil
}
