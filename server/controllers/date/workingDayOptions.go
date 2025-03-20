package date

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func WorkingDayOptions(c *gin.Context) {
	rows, err := config.Database.Query("SELECT id,working_date FROM master_workingdays")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var WorkingDayoptions []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		WorkingDayoptions = append(WorkingDayoptions, map[string]interface{}{
			"label": name,
			"value": id,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, WorkingDayoptions)
}
