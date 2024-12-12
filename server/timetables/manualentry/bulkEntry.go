package manualentry

import (
	"fmt"
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

type TimetableEntry struct {
	DepartmentID int    `json:"department_id"`
	FacultyName  string `json:"faculty_name"`
	SubjectName  string `json:"subject_name"`
	CourseCode   string `json:"course_code"`
	SectionID    int    `json:"section_id"`
	Classroom    string `json:"classroom"`
	SemesterID   int    `json:"semester_id"`
	AcademicYear int    `json:"academic_year"`
	Hour         string `json:"hour"`
	// ... other fields if needed
}

func BulkEntry(c *gin.Context) {
	var entries []TimetableEntry
	err := c.BindJSON(&entries)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := config.Database.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", r)})
		} else if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			tx.Commit()
			c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully!"})
		}
	}()

	for _, entry := range entries {
		startTime, endTime := parseHour(entry.Hour)

		// Prepare and execute the SQL statement within the transaction
		_, err := tx.Exec("INSERT INTO manual_timetable (department_id, faculty_name, subject_name, course_code, section_id, classroom, semester_id, academic_year, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			entry.DepartmentID, entry.FacultyName, entry.SubjectName, entry.CourseCode, entry.SectionID, entry.Classroom, entry.SemesterID, entry.AcademicYear, startTime, endTime)
		if err != nil {
			return
		}
	}
}
func parseHour(hour string) (string, string) {
	switch hour {
	case "PERIOD 1":
		return "08:45:00", "09:35:00"
	case "PERIOD 2":
		return "09:35:00", "10:25:00"
	case "PERIOD 3":
		return "10:40:00", "11:30:00"
	case "PERIOD 4":
		return "13:45:00", "14:35:00"
	case "PERIOD 5":
		return "14:35:00", "15:25:00"
	case "PERIOD 6":
		return "15:40:00", "16:30:00"
	// Default case to handle periods outside the defined range
	default:
		return "", "" // Or return a specific error message
	}
}
