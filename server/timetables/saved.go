package timetables

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetTimetable(c *gin.Context) {
	classroom := c.Param("departmentID")
	sem := c.Param("semesterId")
	AcademicYear := c.Param("academicYearID")
	sectionID := c.Param("sectionID")

	var entries []models.TimetableEntry
	rows, err := config.Database.Query(`
    SELECT id, day_name, start_time, end_time, subject_name, faculty_name, classroom, semester_id,department_id,status,academic_year
    FROM timetable
    WHERE department_id = ? AND semester_id = ? AND  academic_year = ? AND section_id = ?`, classroom, sem, AcademicYear,sectionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timetable: " + err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.TimetableEntry
		if err := rows.Scan(&entry.ID, &entry.DayName, &entry.StartTime, &entry.EndTime, &entry.SubjectName, &entry.FacultyName, &entry.Classroom, &entry.SemesterID, &entry.DepartmentID, &entry.Status, &entry.AcademicYear); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan timetable entry: " + err.Error()})
			return
		}
		entries = append(entries, entry)
	}

	c.JSON(http.StatusOK, entries)
}
