// package routes

// import (
// 	"log"
// 	"net/http"
// 	"ssg-portal/timetables"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// func GenerateTimetable(c *gin.Context) {

// 	departmentIDStr := c.Param("departmentID")
// 	semesterIDStr := c.Param("semesterId")
// 	academicYearIDStr := c.Param("academicYearID")
// 	sectionIDStr := c.Param("sectionID")
// 	departmentID, err := strconv.Atoi(departmentIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
// 		return
// 	}
// 	semesterID, err := strconv.Atoi(semesterIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester ID"})
// 		return
// 	}
// 	academicYearID, err := strconv.Atoi(academicYearIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic year ID"})
// 		return
// 	}
// 	sectionID, err := strconv.Atoi(sectionIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
// 		return
// 	}

// 	// Check if the periods are available
// 	// err = timetables.PeriodsAvailable(departmentID, academicYearID, semesterID, sectionID)
// 	// if err != nil {
// 	// 	log.Printf("Error checking available periods: %v", err)
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	// 	return
// 	// }

// 	// Fetch other details required for timetable generation
// 	semesters, err := timetables.GetSemesterDetails(semesterID)
// 	if err != nil {
// 		log.Printf("Error getting semester details for ID %d: %v", semesterID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve semester details"})
// 		return
// 	}

// 	section, err := timetables.GetSectionDetails(sectionID)
// 	if err != nil {
// 		log.Printf("Error getting section details for ID %d: %v", sectionID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve section details"})
// 		return
// 	}

// 	academicYear, err := timetables.GetAcademicDetails(academicYearID)
// 	if err != nil {
// 		log.Printf("Error getting academic year details for ID %d: %v", academicYearID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve academic year details"})
// 		return
// 	}

// 	days, err := timetables.GetAvailableDays()
// 	if err != nil {
// 		log.Printf("Error getting available days: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve available days"})
// 		return
// 	}

// 	hours, err := timetables.GetHours()
// 	if err != nil {
// 		log.Printf("Error getting hours: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve hours"})
// 		return
// 	}

// 	subjects, err := timetables.GetSubjects(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		log.Printf("Error getting subjects for department ID %d: %v", departmentID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve subjects"})
// 		return
// 	}

// 	faculty, err := timetables.GetFaculty(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		log.Printf("Error getting faculty: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty"})
// 		return
// 	}

// 	classrooms, err := timetables.GetClassrooms(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		log.Printf("Error getting classrooms for department ID %d: %v", departmentID, err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve classrooms"})
// 		return
// 	}

// 	facultySubjects, err := timetables.GetFacultySubjects(semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		log.Printf("Error getting faculty subjects: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty subjects"})
// 		return
// 	}

// 	// Generate the timetable
// 	timetable := timetables.GenerateTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, semesters, section, academicYear, departmentID, semesterID, academicYearID, sectionID)
// 	if timetable == nil {
// 		log.Printf("Error generating timetable")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate timetable"})
// 		return
// 	}

// 	// Return the generated timetable
// 	c.Header("Content-Type", "application/json")
// 	c.JSON(http.StatusOK, timetable)
// }
package routes

import (
	"log"
	"net/http"
	"ssg-portal/timetables"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateTimetable(c *gin.Context) {
	departmentIDStr := c.Param("departmentID")
	semesterIDStr := c.Param("semesterId")
	academicYearIDStr := c.Param("academicYearID")
	sectionIDStr := c.Param("sectionID")
	daysCountStr := c.Param("daysCount") // Get days count from query parameters

	// Convert input strings to integers
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

	academicYearID, err := strconv.Atoi(academicYearIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic year ID"})
		return
	}

	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	// Parse and validate daysCount
	daysCount, err := strconv.Atoi(daysCountStr)
	if err != nil || daysCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days count"})
		return
	}

	// Fetch required details for timetable generation
	semesters, err := timetables.GetSemesterDetails(semesterID)
	if err != nil {
		log.Printf("Error getting semester details for ID %d: %v", semesterID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve semester details"})
		return
	}

	section, err := timetables.GetSectionDetails(sectionID)
	if err != nil {
		log.Printf("Error getting section details for ID %d: %v", sectionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve section details"})
		return
	}

	academicYear, err := timetables.GetAcademicDetails(academicYearID)
	if err != nil {
		log.Printf("Error getting academic year details for ID %d: %v", academicYearID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve academic year details"})
		return
	}

	// // Fetch limited days based on daysCount
	// days, err := timetables.GetAvailableDays(daysCount) // Pass daysCount to limit the days fetched
	// if err != nil {
	// 	log.Printf("Error getting available days: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve available days"})
	// 	return
	// }

	hours, err := timetables.GetHours()
	if err != nil {
		log.Printf("Error getting hours: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve hours"})
		return
	}

	subjects, err := timetables.GetSubjects(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		log.Printf("Error getting subjects for department ID %d: %v", departmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve subjects"})
		return
	}

	faculty, err := timetables.GetFaculty(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		log.Printf("Error getting faculty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty"})
		return
	}

	classrooms, err := timetables.GetClassrooms(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		log.Printf("Error getting classrooms for department ID %d: %v", departmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve classrooms"})
		return
	}

	facultySubjects, err := timetables.GetFacultySubjects(semesterID, academicYearID, sectionID)
	if err != nil {
		log.Printf("Error getting faculty subjects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty subjects"})
		return
	}

	// Generate the timetable
	timetable := timetables.GenerateTimetable(daysCount, hours, subjects, faculty, classrooms, facultySubjects, semesters, section, academicYear, departmentID, semesterID, academicYearID, sectionID)
	if timetable == nil {
		log.Printf("Error generating timetable")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate timetable"})
		return
	}

	// Return the generated timetable
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, timetable)
}
