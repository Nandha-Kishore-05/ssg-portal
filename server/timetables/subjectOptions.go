package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func SubjectOptions(c *gin.Context) {
	rows, err := config.Database.Query("SELECT id, name FROM subjects")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var Subjectoptions []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		Subjectoptions = append(Subjectoptions, map[string]interface{}{
			"label": name,
			"value": name,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Subjectoptions)
}
func CourseCodeOptions(c *gin.Context) {
    // Get the subject name from the query parameter
    subjectName := c.DefaultQuery("subject_name", "")

    if subjectName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Subject name is required"})
        return
    }

    // Query the database to get the course code for the selected subject name
    rows, err := config.Database.Query("SELECT id, course_code FROM subjects WHERE name = ?", subjectName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var courseCodeOptions []map[string]interface{}
    for rows.Next() {
        var id int
        var courseCode string
        if err := rows.Scan(&id, &courseCode); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        courseCodeOptions = append(courseCodeOptions, map[string]interface{}{
            "label": courseCode,
            "value": courseCode,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, courseCodeOptions)
}
