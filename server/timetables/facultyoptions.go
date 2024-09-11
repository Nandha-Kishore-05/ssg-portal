package timetables

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "ssg-portal/config"
)

func FacultyOptions(c *gin.Context) {
    query := `
        SELECT DISTINCT t.faculty_name, f.faculty_id 
        FROM timetable t
        JOIN faculty f ON t.faculty_name = f.name
    `
    rows, err := config.Database.Query(query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var facultyOptions []map[string]string
    for rows.Next() {
        var name, facultyID string
        if err := rows.Scan(&name, &facultyID); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        facultyOptions = append(facultyOptions, map[string]string{
            "label": name,
            "value": name ,
            "id": facultyID,  
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, facultyOptions)
}
