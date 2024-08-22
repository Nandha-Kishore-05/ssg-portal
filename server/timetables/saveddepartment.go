package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func SavedDepartmentOptions(c *gin.Context) {
	rows, err := config.Database.Query(`
	SELECT 
	DISTINCT	d.name AS department_name,
		s.semester_name AS semester_name,
		t.classroom AS classroom
	FROM 
		timetable t
	JOIN 
		departments d ON t.department_id = d.id
	JOIN 
		semester s ON t.semester_id = s.id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var SavedDeptOptions []map[string]string
	for rows.Next() {
		var departmentName string
		var semesterName string
		var classroom string

		if err := rows.Scan(&departmentName, &semesterName, &classroom); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		SavedDeptOptions = append(SavedDeptOptions, map[string]string{
			"department_name": departmentName,
			"semester_name":   semesterName,
			"classroom":       classroom,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SavedDeptOptions)
}
