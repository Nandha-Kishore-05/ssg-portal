package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/timetables"
)

func main() {

	config.ConnectDB()

	r := gin.Default()

	r.Use(cors.Default())

	// Define route for timetable with departmentID parameter
	r.GET("/timetable/:departmentID", func(c *gin.Context) {
		// Get departmentID from URL parameter
		departmentIDStr := c.Param("departmentID")
		departmentID, err := strconv.Atoi(departmentIDStr)
		if err != nil {
			log.Printf("Error converting departmentID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
			return
		}

		// Get available days
		days, err := timetables.GetAvailableDays()
		if err != nil {
			log.Printf("Error getting available days: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve available days"})
			return
		}

		// Get available hours
		hours, err := timetables.GetHours()
		if err != nil {
			log.Printf("Error getting hours: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve hours"})
			return
		}

		subjects, err := timetables.GetSubjects(departmentID)
		if err != nil {
			log.Printf("Error getting subjects: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve subjects"})
			return
		}

		faculty, err := timetables.GetFaculty()
		if err != nil {
			log.Printf("Error getting faculty: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty"})
			return
		}

		classrooms, err := timetables.GetClassrooms(departmentID)
		if err != nil {
			log.Printf("Error getting classrooms: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve classrooms"})
			return
		}

		facultySubjects, err := timetables.GetFacultySubjects()
		if err != nil {
			log.Printf("Error getting faculty subjects: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty subjects"})
			return
		}

		timetable := timetables.GenerateTimetable(days, hours, subjects, faculty, classrooms, facultySubjects)
		c.JSON(http.StatusOK, timetable)
	})
	r.POST("/timetable/save", timetables.SaveTimetable)
	r.Run(":8080")
}
