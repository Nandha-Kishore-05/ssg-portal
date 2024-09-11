package timetables

import (
	"net/http"
	"ssg-portal/config" 
	"ssg-portal/models"  
	"github.com/gin-gonic/gin"
)

func GetMenuItems(c *gin.Context) {
	
	var menuItems []models.MenuItem


	rows, err := config.Database.Query("SELECT id, label, path, icon FROM menu_items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve menu items"})
		return
	}
	defer rows.Close()

	
	for rows.Next() {
		var menuItem models.MenuItem
		if err := rows.Scan(&menuItem.ID, &menuItem.Label, &menuItem.Path, &menuItem.Icon); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan menu item"})
			return
		}
		menuItems = append(menuItems, menuItem)
	}

	
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating rows"})
		return
	}


	c.JSON(http.StatusOK, menuItems)
}
