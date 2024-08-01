package timetables

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetExistingTimetables(c *gin.Context) {
	rows, err := config.Database.Query(`
        SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom 
        FROM timetable`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing timetables: " + err.Error()})
		return
	}
	defer rows.Close()

	timetables := make(map[string][]models.TimetableEntry)
	for rows.Next() {
		var entry models.TimetableEntry
		if err := rows.Scan(&entry.Day, &entry.StartTime, &entry.EndTime, &entry.Subject, &entry.Faculty, &entry.Classroom); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan timetable entry: " + err.Error()})
			return
		}

		class := models.TimetableEntry{
			Day:       entry.Day,
			StartTime: entry.StartTime,
			EndTime:   entry.EndTime,
			Subject:   entry.Subject,
			Faculty:   entry.Faculty,
			Classroom: entry.Classroom,
		}
		timetables[entry.Day] = append(timetables[entry.Day], class)
	}

	c.JSON(http.StatusOK, timetables)
}
