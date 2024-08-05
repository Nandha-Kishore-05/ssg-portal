// package main

// import (
// 	"log"
// 	"net/http"
// 	"strconv"
//
//
//
//
//
//
//
//
//
//
//

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"

// 	"ssg-portal/config"
// 	"ssg-portal/timetables"

// )

// func main() {
// 	config.ConnectDB()

// 	r := gin.Default()
// 	r.Use(cors.Default())

// 	r.GET("/timetable/:departmentID", func(c *gin.Context) {
// 		departmentIDStr := c.Param("departmentID")
// 		if departmentIDStr == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Department ID is required"})
// 			return
// 		}

// 		departmentID, err := strconv.Atoi(departmentIDStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
// 			return
// 		}
// 		// Get available days
// 		days, err := timetables.GetAvailableDays()
// 		if err != nil {
// 			log.Printf("Error getting available days: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve available days"})
// 			return
// 		}

// 		// Get available hours
// 		hours, err := timetables.GetHours()
// 		if err != nil {
// 			log.Printf("Error getting hours: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve hours"})
// 			return
// 		}

// 		subjects, err := timetables.GetSubjects(departmentID)
// 		if err != nil {
// 			log.Printf("Error getting subjects: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve subjects"})
// 			return
// 		}

// 		faculty, err := timetables.GetFaculty()
// 		if err != nil {
// 			log.Printf("Error getting faculty: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty"})
// 			return
// 		}

// 		classrooms, err := timetables.GetClassrooms(departmentID)
// 		if err != nil {
// 			log.Printf("Error getting classrooms: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve classrooms"})
// 			return
// 		}

// 		facultySubjects, err := timetables.GetFacultySubjects()
// 		if err != nil {
// 			log.Printf("Error getting faculty subjects: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty subjects"})
// 			return
// 		}
// 		timetable := timetables.GenerateTimetable(days, hours, subjects, faculty, classrooms, facultySubjects)
// 		if err != nil {
// 			log.Printf("Error generating timetable: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate timetable"})
// 			return
// 		}
// 		c.Header("Content-Type", "application/json")
// 		// Use the generated timetable
// 		//c.JSON(http.StatusOK, gin.H{"departmentID": departmentID})
// 		c.JSON(http.StatusOK, timetable)
// 	})

//		//r.POST("/timetable/save", timetables.SaveTimetable)
//		r.GET("/timetable/saved/:departmentID", timetables.GetTimetable)
//		r.Run(":8080")
//	}
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
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/timetable/:departmentID", func(c *gin.Context) {
		departmentIDStr := c.Param("departmentID")
		if departmentIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Department ID is required"})
			return
		}

		departmentID, err := strconv.Atoi(departmentIDStr)
		if err != nil {
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

		// Generate the timetable
		timetable := timetables.GenerateTimetable(days, hours, subjects, faculty, classrooms, facultySubjects)
		if err != nil {
			log.Printf("Error generating timetable: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate timetable"})
			return
		}
		// Send the timetable as JSON response
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, timetable)
	})

	// Define route for saved timetable (if needed)
	r.POST("/timetable/save", timetables.SaveTimetable)
	r.GET("/timetable/saved/:departmentID", timetables.GetTimetable)
	r.Run(":8080")
}
