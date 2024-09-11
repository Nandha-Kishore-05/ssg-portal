package auth

import (
    "crypto/rand"
    "database/sql"
    "fmt"
    "net/http"
    "ssg-portal/config"

    "github.com/gin-gonic/gin"
)

type UserLogin struct {
    Email string `json:"email"`
}

type UserDetails struct {
    Id        int    `json:"id"`
    Name      string `json:"name"`
    AuthToken string `json:"auth_token"`
}

func Login(c *gin.Context) {
    var input UserLogin
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user UserDetails

    err := config.Database.QueryRow(
        `SELECT mu.id, mu.name 
         FROM master_user mu 
         WHERE mu.email = ? AND mu.status = 1`,
        input.Email,
    ).Scan(&user.Id, &user.Name)

    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }


    tokenBytes := make([]byte, 16)
    _, err = rand.Read(tokenBytes)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
        return
    }
    authToken := fmt.Sprintf("%x", tokenBytes)


    _, err = config.Database.Exec("UPDATE user_login SET auth_token = ? WHERE user_id = ?", authToken, user.Id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


    user.AuthToken = authToken
    c.JSON(http.StatusOK, gin.H{
        "user":       user,
        "auth_token": authToken,
    })
}
