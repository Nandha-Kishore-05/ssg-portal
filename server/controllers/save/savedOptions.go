package save

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"ssg-portal/config"
)

func SaveTimetableOptions(c *gin.Context) {
    rows, err := config.Database.Query("SELECT DISTINCT subject_name FROM timetable where status = 0")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var labOptions []map[string]string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        labOptions = append(labOptions, map[string]string{
            "label": name,
            "value": name,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, labOptions)
}
