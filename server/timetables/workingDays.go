package timetables

import (
	"log"
	"ssg-portal/config"
	"ssg-portal/models"
	"time"
)

// GetWorkingDaysInRange retrieves working days for a given academic year, semester, and ID range
func GetWorkingDaysInRange(academicYearID int, semesterID int, startID int, endID int) ([]models.WorkingDay, error) {
	// Define the query to get working days for the given academic year, semester, and id range
	query := `
        SELECT working_date
        FROM master_workingdays
        WHERE academic_year_id = ? 
        AND semester_id = ? 
        AND id BETWEEN ? AND ?
        ORDER BY working_date;
    `

	// Prepare the statement
	rows, err := config.Database.Query(query, academicYearID, semesterID, startID, endID)
	if err != nil {
		log.Printf("Error fetching working days: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the working days
	var workingDays []models.WorkingDay

	// Iterate through the result set and populate the workingDays slice
	for rows.Next() {
		var wd models.WorkingDay
		var workingDate []byte // Temporarily scan as []byte
		if err := rows.Scan(&workingDate); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Convert []byte to time.Time
		parsedDate, err := time.Parse("2006-01-02", string(workingDate))
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			return nil, err
		}

		wd.WorkingDate = parsedDate
		workingDays = append(workingDays, wd)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error during iteration: %v", err)
		return nil, err
	}

	// Return the list of working days
	return workingDays, nil
}
