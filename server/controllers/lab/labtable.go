package lab

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func LabTableTimetable(c *gin.Context) {
	LabTable := c.Param("subject_name")

	timetableEntries, err := getLabTimetable(LabTable)
	if err != nil {
		fmt.Println("Error fetching timetable:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetableEntries)
}

func getLabTimetable(subjectName string) ([]models.LabTimetableEntry, error) {
	query := `
		SELECT day_name, start_time, end_time,semester_id,subject_name,faculty_name
		FROM timetable
		WHERE subject_name = ? AND status = 0
	`

	rows, err := config.Database.Query(query, subjectName)
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, err
	}
	defer rows.Close()

	var timetableEntries []models.LabTimetableEntry

	for rows.Next() {
		var entry models.LabTimetableEntry
		if err := rows.Scan(&entry.DayName, &entry.StartTime, &entry.EndTime, &entry.SemesterID, &entry.SubjectName, &entry.FacultyName); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		timetableEntries = append(timetableEntries, entry)
	}

	if len(timetableEntries) == 0 {
		fmt.Println("No entries found for lab:", subjectName)
	}

	return timetableEntries, nil
}
