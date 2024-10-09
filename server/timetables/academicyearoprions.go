package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func AcademicYearOptions(c *gin.Context) {
	// Query only distinct academic years
	rows, err := config.Database.Query("SELECT  id, academic_year FROM  master_academic_year")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var academicOptions []map[string]interface{}
	for rows.Next() {
		var id int
		var year string
		if err := rows.Scan(&id, &year); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		academicOptions = append(academicOptions, map[string]interface{}{
			"label": year,
			"value": id,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, academicOptions)
}
