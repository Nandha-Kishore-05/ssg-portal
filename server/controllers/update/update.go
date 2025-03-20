package update

import (
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)

func UpdateTimetable(c *gin.Context) {
	var timetableEntries []models.TimetableEntry

	if err := c.ShouldBindJSON(&timetableEntries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tx, err := config.Database.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		}
	}()

	for _, entry := range timetableEntries {

		_, err := tx.Exec(`
			UPDATE timetable 
			SET day_name = ?, start_time = ?, end_time = ?, subject_name = ?, faculty_name = ?, classroom = ?
			WHERE id = ?
		`, entry.DayName, entry.StartTime, entry.EndTime, entry.SubjectName, entry.FacultyName, entry.Classroom, entry.ID)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timetable entry"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Timetable updated successfully"})
}
