package role

import (
    "net/http"
    "ssg-portal/config"
    function "ssg-portal/functions"
    "ssg-portal/models"

    "github.com/gin-gonic/gin"
)

func GetResources(c *gin.Context) {
    auth, UserId := function.CheckAuth(c)

    if !auth {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var data []models.Resource
    var temp models.Resource

    rows, err := config.Database.Query(`
        SELECT m.id, m.name, m.icon, m.path, m.sort_by 
        FROM master_resource m 
        WHERE FIND_IN_SET(m.id, (
            SELECT r.resources 
            FROM master_roles r
            INNER JOIN master_user mu ON mu.role = r.id
            WHERE mu.id = ?
        )) > 0 
        ORDER BY m.sort_by
    `, UserId)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve resources"})
        return
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&temp.ID, &temp.Name, &temp.Icon, &temp.Path, &temp.SortBy)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process resource data"})
            return
        }
        data = append(data, temp)
    }

    if len(data) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No resources found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": data})
}
