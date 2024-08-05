package timetables

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func FacultyTimetable(c *gin.Context) {
	facultyName := c.Param("faculty_name")

	timetableEntries, err := getFacultyTimetable(facultyName)
	if err != nil {
		fmt.Println("Error fetching timetable:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetableEntries)
}

func getFacultyTimetable(facultyName string) ([]models.FacultyTimetableEntry, error) {
	query := `
		SELECT day_name, start_time, end_time, classroom
		FROM timetable
		WHERE faculty_name = ? 
	`

	rows, err := config.Database.Query(query, facultyName)
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, err
	}
	defer rows.Close()

	var timetableEntries []models.FacultyTimetableEntry

	for rows.Next() {
		var entry models.FacultyTimetableEntry
		if err := rows.Scan(&entry.DayName, &entry.StartTime, &entry.EndTime, &entry.Classroom); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		timetableEntries = append(timetableEntries, entry)
	}

	if len(timetableEntries) == 0 {
		fmt.Println("No entries found for faculty:", facultyName)
	}

	return timetableEntries, nil
}
