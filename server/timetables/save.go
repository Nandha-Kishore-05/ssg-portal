package timetables

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func SaveTimetable(c *gin.Context) {
	var entries []models.TimetableEntry
	if err := c.BindJSON(&entries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.Database

	stmt, err := db.Prepare(`
		INSERT INTO timetable (day_name, start_time, end_time, subject_name, faculty_name, classroom,status,semester_id,department_id,academic_year,course_code,section_id)
		VALUES (?, ?, ?, ?, ?, ?,?,?,?,?,?,?)
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement: " + err.Error()})
		return
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err := stmt.Exec(
			entry.DayName,
			entry.StartTime,
			entry.EndTime,
			entry.SubjectName,
			entry.FacultyName,
			entry.Classroom,
			entry.Status,
			entry.SemesterID,
			entry.DepartmentID,
			entry.AcademicYear,
			entry.CourseCode,
			entry.SectionID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save timetable entry: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable saved successfully!"})
}
