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

	// Prepare the SQL query
	query := `
		INSERT INTO timetable_skips (subject_name, department_id, semester_id, day_name, start_time, end_time, faculty_name, classroom, academic_year, course_code, section_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Loop through each request and insert into the database
	for _, request := range requests {
		_, err := config.Database.Exec(query, request.Subject, request.Department, request.Semester, request.Day, request.StartTime, request.EndTime, request.Faculty, request.Classroom, request.AcademicYear, request.CourseCode, request.SectionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manual entries submitted successfully"})
}
