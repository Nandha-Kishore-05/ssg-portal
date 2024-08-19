package timetables

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"ssg-portal/config"
)

func FacultyOptions(c *gin.Context) {
    rows, err := config.Database.Query("SELECT DISTINCT faculty_name FROM timetable")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var facultyOptions []map[string]string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        facultyOptions = append(facultyOptions, map[string]string{
            "label": name,
            "value": name,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, facultyOptions)
}
