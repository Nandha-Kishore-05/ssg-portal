package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetClassrooms(departmentID int) ([]models.Classroom, error) {
	var classrooms []models.Classroom
	query := `SELECT id, name FROM classrooms WHERE department_id = ?`
	rows, err := config.Database.Query(query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("error querying classrooms: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var classroom models.Classroom
		if err := rows.Scan(&classroom.ID, &classroom.ClassroomName); err != nil {
			return nil, fmt.Errorf("error scanning classroom: %v", err)
		}
		classrooms = append(classrooms, classroom)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return classrooms, nil
}
