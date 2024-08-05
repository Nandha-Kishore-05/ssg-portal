package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetAvailableDays() ([]models.Day, error) {
	var days []models.Day
	rows, err := config.Database.Query("SELECT id, day_name FROM days ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("error querying days: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var day models.Day
		if err := rows.Scan(&day.ID, &day.DayName); err != nil {
			return nil, fmt.Errorf("error scanning day: %v", err)
		}
		days = append(days, day)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return days, nil
}

// Define similar functions for GetHours, GetSubjects, GetFaculty, GetClassrooms, and GetFacultySubjects
