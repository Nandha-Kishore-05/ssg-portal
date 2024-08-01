package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetHours() ([]models.Hour, error) {
	var hours []models.Hour
	rows, err := config.Database.Query("SELECT id, start_time, end_time FROM hours")
	if err != nil {
		return nil, fmt.Errorf("error querying hours: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var hour models.Hour
		if err := rows.Scan(&hour.ID, &hour.StartTime, &hour.EndTime); err != nil {
			return nil, fmt.Errorf("error scanning hour: %v", err)
		}
		hours = append(hours, hour)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return hours, nil
}
