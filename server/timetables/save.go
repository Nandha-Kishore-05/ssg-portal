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
	for _, entry := range entries {
		_, err := db.Exec(`
            INSERT INTO timetable (day_name, start_time, end_time, subject_name, faculty_name, classroom)
            VALUES (?, ?, ?, ?, ?, ?)`,
			entry.Day,
			entry.StartTime,
			entry.EndTime,
			entry.Subject,
			entry.Faculty,
			entry.Classroom,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save timetable: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable saved successfully!"})
}
