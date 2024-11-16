package timetables

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "ssg-portal/config"
)


func TimetableOptions(c *gin.Context) {
    rows, err := config.Database.Query("SELECT id, name FROM departments")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var options []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        options = append(options, map[string]interface{}{
            "label": name,
            "value": id,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, options)
}
func VenueOptions(c *gin.Context) {
    departmentID := c.Query("department_id")
    academicYearID := c.Query("academic_year_id")
    semesterID := c.Query("semester_id")

    // Validate required parameters
    if departmentID == "" || academicYearID == "" || semesterID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "department_id, academic_year_id, and semester_id are required"})
        return
    }

    // Query to get classrooms based on the selected department, academic year, and semester
    query := `
        SELECT id, name 
        FROM classrooms 
        WHERE department_id = ? AND academic_year_id = ? AND semester_id = ?
    `
    rows, err := config.Database.Query(query, departmentID, academicYearID, semesterID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var venues []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        venues = append(venues, map[string]interface{}{
            "label": name,
            "value": id,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, venues)
}
