package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func LabOptions(c *gin.Context) {
	rows, err := config.Database.Query(" SELECT DISTINCT t.subject_name, a.academic_year,t.academic_year FROM timetable t JOIN master_academic_year a ON t.academic_year = a.id WHERE t.status = 0 ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var labOptions []map[string]string
	for rows.Next() {
		var name, academicYearName, academicYearID string
		if err := rows.Scan(&name, &academicYearName, &academicYearID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		labOptions = append(labOptions, map[string]string{
			"label":         name,
			"value":         name,
			"academic_year": academicYearName,
			"academic_id":   academicYearID,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, labOptions)
}
