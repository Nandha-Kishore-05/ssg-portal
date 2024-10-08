package manualentry

import (
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)

func SubmitManualEntry(c *gin.Context) {
	var request models.ManualEntryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
        INSERT INTO timetable_skips (subject_name, department_id, semester_id, day_name, start_time, end_time, faculty_name,  classroom,academic_year,course_code)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?,?,?)
    `

	_, err := config.Database.Exec(query, request.Subject, request.Department, request.Semester, request.Day, request.StartTime, request.EndTime, request.Faculty, request.Classroom, request.AcademicYear, request.CourseCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manual entry submitted successfully"})
}
