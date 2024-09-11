package function

import (
    "fmt"
    "net/http"
    "ssg-portal/config"

    "github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) (bool, string) {
    authHeader := c.GetHeader("Authorization")

    fmt.Println("Authorization Header:", authHeader)

    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Not Found"})
        return false, ""
    }

    fmt.Println("Extracted Token:", authHeader)

    var userID string
    err := config.Database.QueryRow("SELECT user_id FROM user_login WHERE auth_token = ?", authHeader).Scan(&userID)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Mismatch"})
        return false, ""
    }

    return true, userID
}
