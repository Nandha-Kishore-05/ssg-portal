package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

// ClassroomOptions fetches classroom names along with their academic years and IDs
func ClassroomOptions(c *gin.Context) {
	// SQL query to fetch distinct classroom names, academic year names, and academic year IDs
	query := `
		SELECT DISTINCT t.classroom, ay.academic_year, ay.id 
		FROM timetable t 
		JOIN academic_year ay ON t.academic_year = ay.id
	`

	rows, err := config.Database.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Prepare a slice to hold the venue options
	var venueOptions []map[string]interface{}

	// Iterate over the rows and populate the venueOptions slice
	for rows.Next() {
		var classroomName, academicYearName string
		var academicYearID int
		if err := rows.Scan(&classroomName, &academicYearName, &academicYearID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		venueOptions = append(venueOptions, map[string]interface{}{
			"classroom":      classroomName,
			"classroomValue": classroomName,
			"academicyear":   academicYearName,
			"year_id":        academicYearID,
		})
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the venue options as JSON
	c.JSON(http.StatusOK, venueOptions)
}
