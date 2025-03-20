package facultyDetails

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func FacultyTimetable(c *gin.Context) {
	facultyName := c.Param("faculty_name")
	AcademicYear := c.Param("academicYearID")
	academicYearID, err := strconv.Atoi(AcademicYear)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
		return
	}
	timetableEntries, err := getFacultyTimetable(facultyName,academicYearID)
	if err != nil {
		fmt.Println("Error fetching timetable:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetableEntries)
}

func getFacultyTimetable(facultyName string,AcademicYear int) ([]models.FacultyTimetableEntry, error) {
	query := `
		SELECT day_name, start_time, end_time, classroom,semester_id,subject_name
		FROM timetable
		WHERE faculty_name = ? AND academic_year = ?
	`

	rows, err := config.Database.Query(query, facultyName,AcademicYear)
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, err
	}
	defer rows.Close()

	var timetableEntries []models.FacultyTimetableEntry

	for rows.Next() {
		var entry models.FacultyTimetableEntry
		if err := rows.Scan(&entry.DayName, &entry.StartTime, &entry.EndTime, &entry.Classroom, &entry.SemesterID, &entry.SubjectName); err != nil {
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
