package timetables

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetTimetable(c *gin.Context) {
	classroom := c.Param("departmentID")

	var name string
	err := config.Database.QueryRow(`
		SELECT  name
		FROM classrooms
		WHERE department_id = ?`, classroom).Scan(&name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department ID: " + err.Error()})
		return
	}

	var entries []models.TimetableEntry
	rows, err := config.Database.Query(`
    SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom
    FROM timetable
    WHERE classroom = ?`, name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timetable: " + err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.TimetableEntry
		if err := rows.Scan(&entry.DayName, &entry.StartTime, &entry.EndTime, &entry.SubjectName, &entry.FacultyName, &entry.Classroom); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan timetable entry: " + err.Error()})
			return
		}
		entries = append(entries, entry)
	}

	c.JSON(http.StatusOK, entries)
}
