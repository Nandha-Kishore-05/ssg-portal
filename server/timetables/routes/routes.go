package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ssg-portal/timetables"
)

func GenerateTimetable(c *gin.Context) {
	// Extract departmentID and semester_id from URL
	departmentIDStr := c.Param("departmentID")
	semesterIDStr := c.Param("semesterId")

	// Correct logging for string values
	log.Printf("Department ID: %s", departmentIDStr)
	log.Printf("Semester ID: %s", semesterIDStr)

	// Convert departmentID and semesterID from string to int
	departmentID, err := strconv.Atoi(departmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}
	semesterID, err := strconv.Atoi(semesterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
		return
	}

	// Fetch semester details
	semesters, err := timetables.GetSemesterDetails(semesterID)
	if err != nil {
		log.Printf("Error getting semester details for ID %d: %v", semesterID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve semester details"})
		return
	}

	// Fetch available days
	days, err := timetables.GetAvailableDays()
	if err != nil {
		log.Printf("Error getting available days: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve available days"})
		return
	}

	// Fetch available hours
	hours, err := timetables.GetHours()
	if err != nil {
		log.Printf("Error getting hours: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve hours"})
		return
	}

	// Fetch subjects
	subjects, err := timetables.GetSubjects(departmentID, semesterID)
	if err != nil {
		log.Printf("Error getting subjects for department ID %d: %v", departmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve subjects"})
		return
	}

	// Fetch faculty
	faculty, err := timetables.GetFaculty()
	if err != nil {
		log.Printf("Error getting faculty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty"})
		return
	}

	// Fetch classrooms
	classrooms, err := timetables.GetClassrooms(departmentID, semesterID)
	if err != nil {
		log.Printf("Error getting classrooms for department ID %d: %v", departmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve classrooms"})
		return
	}

	// Fetch faculty subjects
	facultySubjects, err := timetables.GetFacultySubjects(semesterID)
	if err != nil {
		log.Printf("Error getting faculty subjects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty subjects"})
		return
	}

	// Generate timetable
	timetable := timetables.GenerateTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, semesters)
	if timetable == nil {
		log.Printf("Error generating timetable")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate timetable"})
		return
	}

	// Send the timetable as JSON response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, timetable)
}
