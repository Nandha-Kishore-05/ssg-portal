package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func VenueOptions(c *gin.Context) {
	rows, err := config.Database.Query("SELECT id, name FROM classrooms")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var venueoptions []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		venueoptions = append(venueoptions, map[string]interface{}{
			"label": name,
			"value": name,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, venueoptions)
}
