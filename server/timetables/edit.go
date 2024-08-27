package timetables

import (
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)

func GetAvailableFaculty(c *gin.Context) {
	departmentID := c.Param("departmentID")
	semesterID := c.Param("semesterID")
	day := c.Param("day")
	startTime := c.Param("startTime")
	endTime := c.Param("endTime")

	facultyList, err := fetchAvailableFaculty(departmentID, semesterID, day, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, facultyList)
}

// fetchAvailableFaculty queries the database to find available faculty based on parameters.
func fetchAvailableFaculty(departmentID, semesterID, day, startTime, endTime string) ([]models.Faculty, error) {
	var facultyList []models.Faculty

	// Define the SQL query
	query := `
		SELECT f.id, f.name, f.department_id, s.name
		FROM faculty f
		JOIN faculty_subjects fs ON f.id = fs.faculty_id
		JOIN subjects s ON fs.subject_id = s.id
		WHERE f.department_id = ? AND fs.semester_id = ? AND NOT EXISTS (
			SELECT 1 FROM timetable t
			WHERE t.faculty_name = f.name AND t.day_name = ? AND t.start_time <= ? AND t.end_time >= ?
		)`

	// Execute the query
	rows, err := config.Database.Query(query, departmentID, semesterID, day, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results
	for rows.Next() {
		var faculty models.Faculty
		if err := rows.Scan(&faculty.ID, &faculty.FacultyName, &faculty.DepartmentID, &faculty.SubjectName); err != nil {
			return nil, err
		}
		facultyList = append(facultyList, faculty)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return facultyList, nil
}
