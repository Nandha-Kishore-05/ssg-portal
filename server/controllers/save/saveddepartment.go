package save

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

// SavedDepartmentOptions retrieves distinct department options related to the timetable
func SavedDepartmentOptions(c *gin.Context) {
	rows, err := config.Database.Query(`
	SELECT 
		DISTINCT d.name AS department_name,
		d.id AS department_id,
		se.id AS section_id,
		se.section_name AS section,
		s.semester_name AS semester_name,
		s.id AS semester_id,
		may.id AS academic_year_id,
		may.academic_year AS academic_year_name
	FROM 
		timetable t
	JOIN 
		departments d ON t.department_id = d.id
	JOIN 
		semester s ON t.semester_id = s.id
	JOIN
		master_section se ON t.section_id = se.id
	JOIN 
		master_academic_year may ON t.academic_year = may.id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var SavedDeptOptions []map[string]string
	for rows.Next() {
		var departmentName, departmentID, section, sectionID, semesterName, semesterID, academicYearID, academicYearName string

		// Correctly match the order of fields in the SELECT statement
		if err := rows.Scan(&departmentName, &departmentID, &sectionID, &section, &semesterName, &semesterID, &academicYearID, &academicYearName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		SavedDeptOptions = append(SavedDeptOptions, map[string]string{
			"department_name":    departmentName,
			"department_id":      departmentID,
			"section_id":         sectionID,
			"classroom":          section, // Changed to section for clarity
			"semester_name":      semesterName,
			"semester_id":        semesterID,
			"academic_year_id":   academicYearID,
			"academic_year_name": academicYearName,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SavedDeptOptions)
}
