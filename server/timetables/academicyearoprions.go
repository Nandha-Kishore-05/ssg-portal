package timetables

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"ssg-portal/config"
)

func AcademicYearOptions(c *gin.Context) {
    rows, err := config.Database.Query("SELECT id, academic_year FROM academic_year ")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var Academicsoptions []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        Academicsoptions  = append(Academicsoptions , map[string]interface{}{
            "label": name,
            "value": id,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, Academicsoptions )
}
