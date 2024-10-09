package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetAcademicDetails(academicYearID int) ([]models.AcademicYear, error) {
	var academic []models.AcademicYear
	query := `SELECT id, academic_year FROM master_academic_year WHERE id = ?`
	rows, err := config.Database.Query(query, academicYearID)
	if err != nil {
		return nil, fmt.Errorf("error querying Academics: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var academicYear models.AcademicYear
		if err := rows.Scan(&academicYear.AcademicYear, &academicYear.AcademicYearName); err != nil {
			return nil, fmt.Errorf("error scanning academics: %v", err)
		}
		academic = append(academic, academicYear)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return academic, nil
}
