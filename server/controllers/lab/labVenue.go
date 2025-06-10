// package lab

// import (
// 	"fmt"
// 	//"log"

// 	"ssg-portal/config"
// 	"ssg-portal/models"
// )

// func GetLabVenue() ([]models.LabVenue, error) {
// 	//log.Println("Academic YearDSCFEF ID:")
// 	var Lab []models.LabVenue
// 	rows, err := config.Database.Query("SELECT id, lab_name,subject_id FROM lab_venue")
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying faculty: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var Venue models.LabVenue
// 		if err := rows.Scan(&Venue.ID, &Venue.LabVenue, &Venue.SubjectID); err != nil {
// 			return nil, fmt.Errorf("error scanning faculty: %v", err)
// 		}
// 		Lab = append(Lab, Venue)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating rows: %v", err)
// 	}
// 	//log.Println("Academic YearDSEFFEFEF3EBUCFEF ID:")
// 	return Lab, nil

// }

package lab

import (
	"fmt"
	"log"
	//"log"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetLabVenue() ([]models.LabVenue, error) {
	//log.Println("Academic YearDSCFEF ID:")
	var Lab []models.LabVenue
	rows, err := config.Database.Query("SELECT id, lab_name, subject_id, max_sections FROM lab_venue")
	if err != nil {
		return nil, fmt.Errorf("error querying lab venues: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var Venue models.LabVenue
		if err := rows.Scan(&Venue.ID, &Venue.LabVenue, &Venue.SubjectID, &Venue.MaxSections); err != nil {
			return nil, fmt.Errorf("error scanning lab venue: %v", err)
		}
		Lab = append(Lab, Venue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	log.Println("Academic YearDSEFFEFEF3EBUCFEF ID:",Lab)
	return Lab, nil
}



func GetLabOccupancyFromDB(day, startTime string, academicYearID int) (map[string]int, error) {
	labOccupancy := make(map[string]int) // lab_name -> current_sections_count
	
	query := `
		SELECT classroom, COUNT(*) as section_count 
		FROM timetable 
		WHERE day_name = ? AND start_time = ? AND academic_year = ? AND status = 0
		GROUP BY classroom`
	
	rows, err := config.Database.Query(query, day, startTime, academicYearID)
	if err != nil {
		return nil, fmt.Errorf("error querying lab occupancy: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var labName string
		var sectionCount int
		if err := rows.Scan(&labName, &sectionCount); err != nil {
			return nil, fmt.Errorf("error scanning lab occupancy: %v", err)
		}
		labOccupancy[labName] = sectionCount
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	
	return labOccupancy, nil
}