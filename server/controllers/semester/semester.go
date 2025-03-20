package semester

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetSemesterDetails(semesterID int) ([]models.Semester, error) {
	var semesters []models.Semester
	query := `SELECT id, semester_name, year_id FROM semester WHERE id = ?`
	rows, err := config.Database.Query(query, semesterID)
	if err != nil {
		return nil, fmt.Errorf("error querying semester: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var semester models.Semester
		if err := rows.Scan(&semester.ID, &semester.SemesterName, &semester.YearID); err != nil {
			return nil, fmt.Errorf("error scanning semester: %v", err)
		}
		semesters = append(semesters, semester)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return semesters, nil
}
