package timetables

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetLabVenue() ([]models.LabVenue, error) {
	var Lab []models.LabVenue
	rows, err := config.Database.Query("SELECT id, lab_name,subject_id FROM lab_venue")
	if err != nil {
		return nil, fmt.Errorf("error querying faculty: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var Venue models.LabVenue
		if err := rows.Scan(&Venue.ID, &Venue.LabVenue, &Venue.SubjectID); err != nil {
			return nil, fmt.Errorf("error scanning faculty: %v", err)
		}
		Lab = append(Lab, Venue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return Lab, nil
}