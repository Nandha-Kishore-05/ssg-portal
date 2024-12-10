package manualentry

import (
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)

func SubmitManualEntry(c *gin.Context) {
	var requests []models.ManualEntryRequest // Change to a slice to accept multiple entries

	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation query to check for conflicts
	validationQuery := `
		SELECT COUNT(*) 
		FROM timetable_skips 
		WHERE academic_year = ? AND start_time = ? AND end_time = ? AND day_name = ? AND status != ?
	`

	// Insert query
	insertQuery := `
		INSERT INTO timetable_skips (subject_name, department_id, semester_id, day_name, start_time, end_time, faculty_name, classroom, academic_year, course_code, section_id, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, request := range requests {
		// Check for conflicting entries
		var conflictCount int
		err := config.Database.QueryRow(validationQuery, request.AcademicYear, request.StartTime, request.EndTime, request.Day, request.Status).Scan(&conflictCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate data"})
			return
		}

		// If a conflict with a different status is found, return an error message
		if conflictCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Conflict detected: An entry with the same start time, end time, and day but a different status already exists.",
			})
			return
		}

		// Proceed to insert the entry if no conflict
		_, err = config.Database.Exec(insertQuery, request.Subject, request.Department, request.Semester, request.Day, request.StartTime, request.EndTime, request.Faculty, request.Classroom, request.AcademicYear, request.CourseCode, request.SectionID, request.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manual entries submitted successfully"})
}
