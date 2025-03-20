package semester

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"ssg-portal/config"
)

func SemOptions(c *gin.Context) {
    rows, err := config.Database.Query("SELECT id, semester_name FROM semester")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var semoptions []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        semoptions = append(semoptions, map[string]interface{}{
            "label": name,
            "value": id,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, semoptions)
}
