package section

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "ssg-portal/config"
)

func SectionOptions(c *gin.Context) {
    rows, err := config.Database.Query("SELECT id, section_name FROM master_section")
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
