package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func SavedDepartmentOptions(c *gin.Context) {
	rows, err := config.Database.Query(`
	SELECT 
		DISTINCT d.name AS department_name,
		d.id AS department_id,
		t.classroom,
		t.section_id,
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
		master_academic_year may ON t.academic_year = may.id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var SavedDeptOptions []map[string]string
	for rows.Next() {
		var departmentName string
		var departmentID string
		var classroom string
		var sectionID string
		var semesterName string
		var semesterID string
		var academicYearID string
		var academicYearName string

		// The order here should match the SELECT fields' order
		if err := rows.Scan(&departmentName, &departmentID, &classroom, &sectionID,&semesterName, &semesterID, &academicYearID, &academicYearName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		SavedDeptOptions = append(SavedDeptOptions, map[string]string{
			"department_name":   departmentName,
			"department_id":     departmentID,
			"classroom":         classroom,
			"section_id":        sectionID,
			"semester_name":     semesterName,
			"semester_id":       semesterID,
			"academic_year_id":  academicYearID,
			"academic_year_name": academicYearName,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SavedDeptOptions)
}
