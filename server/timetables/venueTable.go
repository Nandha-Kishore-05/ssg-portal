package timetables

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func VenueTimetable(c *gin.Context) {
	venueTable := c.Param("classroom")

	timetableEntries, err := getvenueTimetable(venueTable)
	if err != nil {
		fmt.Println("Error fetching timetable:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetableEntries)
}

func getvenueTimetable(classroom string) ([]models.VenueTimetable, error) {
	query := `
		SELECT day_name, start_time, end_time,semester_id,subject_name,faculty_name
		FROM timetable
		WHERE classroom = ? 
	`

	rows, err := config.Database.Query(query, classroom)
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, err
	}
	defer rows.Close()

	var timetableEntries []models.VenueTimetable

	for rows.Next() {
		var entry models.VenueTimetable
		if err := rows.Scan(&entry.DayName, &entry.StartTime, &entry.EndTime, &entry.SemesterID, &entry.SubjectName, &entry.FacultyName); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		timetableEntries = append(timetableEntries, entry)
	}

	if len(timetableEntries) == 0 {
		fmt.Println("No entries found for lab:", classroom)
	}

	return timetableEntries, nil
}