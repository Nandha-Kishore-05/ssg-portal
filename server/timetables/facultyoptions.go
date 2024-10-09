package timetables

import (
	"net/http"

	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func FacultyOptions(c *gin.Context) {
	query := `
    SELECT DISTINCT 
    t.faculty_name, 
    f.faculty_id,  
    may.academic_year,
    t.academic_year
FROM 
    timetable t
JOIN 
    faculty f ON t.faculty_name = f.name

JOIN 
    master_academic_year may ON t.academic_year = may.id
    `
	rows, err := config.Database.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var facultyOptions []map[string]string
	for rows.Next() {
		var name, facultyID, academicYearName, academicYearID string

		if err := rows.Scan(&name, &facultyID, &academicYearName, &academicYearID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		facultyOptions = append(facultyOptions, map[string]string{
			"label":         name,
			"value":         name,
			"academic_year": academicYearName,
			"id":            facultyID,
			"academic_id":   academicYearID,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, facultyOptions)
}
