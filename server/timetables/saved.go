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

	var entries []models.TimetableEntry
	rows, err := config.Database.Query(`
    SELECT id, day_name, start_time, end_time, subject_name, faculty_name, classroom, semester_id,department_id,status,academic_year
    FROM timetable
    WHERE department_id = ? AND semester_id = ? AND  academic_year = ?`, classroom, sem, AcademicYear)

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

// package timetables

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"

// 	"ssg-portal/config"
// 	"ssg-portal/models"
// )

// func GetTimetable(c *gin.Context) {
// 	classroom := c.Param("departmentID")
// 	sem := c.Param("semesterId")
// 	var name string

// 	// Query to get the classroom name based on department and semester
// 	err := config.Database.QueryRow(
// 		`SELECT name
// 		FROM classrooms
// 		WHERE department_id = ? AND semester_id = ?`, classroom, sem).Scan(&name)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department ID: " + err.Error()})
// 		return
// 	}

// 	var entries []models.TimetableEntry

// 	rows, err := config.Database.Query(
// 		`SELECT t.id, t.day_name, t.start_time, t.end_time, s.id AS subject_id, f.id AS faculty_id, t.classroom, t.semester_id, t.department_id, t.status
// 		FROM timetable t
// 		JOIN subjects s ON t.subject_name = s.name
// 		JOIN faculty f ON t.faculty_name = f.name
// 		WHERE t.classroom = ? AND t.name = ?`, name, classroom)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timetable: " + err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var entry models.TimetableEntry
// 		if err := rows.Scan(&entry.ID, &entry.DayName, &entry.StartTime, &entry.EndTime, &entry.SubjectID, &entry.FacultyID, &entry.Classroom, &entry.SemesterID, &entry.DepartmentID, &entry.Status); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan timetable entry: " + err.Error()})
// 			return
// 		}
// 		entries = append(entries, entry)
// 	}

// 	c.JSON(http.StatusOK, entries)
// }
