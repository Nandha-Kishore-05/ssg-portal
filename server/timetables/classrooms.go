package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetClassrooms(departmentID int, semesterID int, academicYearID int, sectionID int) ([]models.Classroom, error) {
	var classrooms []models.Classroom
	query := `SELECT id, name,semester_id FROM classrooms WHERE department_id = ? && semester_id = ? && academic_year_id = ? && section_id = ?`
	rows, err := config.Database.Query(query, departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying classrooms: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var classroom models.Classroom
		if err := rows.Scan(&classroom.ID, &classroom.ClassroomName, &classroom.SemesterID); err != nil {
			return nil, fmt.Errorf("error scanning classroom: %v", err)
		}
		classrooms = append(classrooms, classroom)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return classrooms, nil
}
