package groups

import (
	"log"
	"net/http"
	"ssg-portal/controllers/academicYear"
	"ssg-portal/controllers/classroom"
	"ssg-portal/controllers/days"
	"ssg-portal/controllers/facultyDetails"
	"ssg-portal/controllers/hours"
	"ssg-portal/controllers/section"
	"ssg-portal/controllers/semester"
	"ssg-portal/controllers/subject"
	"ssg-portal/controllers/timetable"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateTimetable(c *gin.Context) {
	departmentIDStr := c.Param("departmentID")
	semesterIDStr := c.Param("semesterId")
	academicYearIDStr := c.Param("academicYearID")
	sectionIDStr := c.Param("sectionID")
	// startDateStr := c.Param("startdate")
	// endDateStr := c.Param("enddate")
	log.Println("API ENTERED")
	// Convert input strings to integers
	log.Println("dept",departmentIDStr )
	log.Println("sem",semesterIDStr )
	log.Println("academic ",academicYearIDStr )
	log.Println("section",sectionIDStr )
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

	// Parse start and end dates
	// startDate, err := strconv.Atoi(startDateStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
	// 	return
	// }

	// endDate, err := strconv.Atoi(endDateStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
	// 	return
	// }

	// // Fetch working and holiday days within the date range
	// workingDays, err := timetables.GetWorkingDaysInRange(academicYearID, semesterID, startDate, endDate)
	// if err != nil {
	// 	log.Printf("Error fetching working days: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve working days"})
	// 	return
	// }

	days, err := days.GetAvailableDays()
	if err != nil {
		log.Printf("Error getting available days: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve available days"})
		return
	}

	// Fetch other required details for timetable generation
	semesters, err := semester.GetSemesterDetails(semesterID)
	if err != nil {
		log.Printf("Error getting semester details for ID %d: %v", semesterID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve semester details"})
		return
	}

	section, err :=  section.GetSectionDetails(sectionID)
	if err != nil {
		log.Printf("Error getting section details for ID %d: %v", sectionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve section details"})
		return
	}

	academicYear, err :=  academicYear.GetAcademicDetails(academicYearID)
	if err != nil {
		log.Printf("Error getting academic year details for ID %d: %v", academicYearID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve academic year details"})
		return
	}

	hours, err :=  hours.GetHours()
	if err != nil {
		log.Printf("Error getting hours: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve hours"})
		return
	}

	subjects, err :=  subject.GetSubjects(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		log.Printf("Error getting subjects for department ID %d: %v", departmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve subjects"})
		return
	}

	faculty, err :=  facultyDetails.GetFaculty(departmentID, semesterID, academicYearID)
	if err != nil {
		log.Printf("Error getting faculty: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty"})
		return
	}

	classrooms, err := classroom.GetClassrooms(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		log.Printf("Error getting classrooms for department ID %d: %v", departmentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve classrooms"})
		return
	}

	facultySubjects, err :=  facultyDetails.GetFacultySubjects(semesterID, academicYearID, sectionID)
	if err != nil {
		log.Printf("Error getting faculty subjects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve faculty subjects"})
		return
	}

	// Generate timetable by iterating over working days between start and end date
	timetable :=  timetable.GenerateTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, semesters, section, academicYear, departmentID, semesterID, academicYearID, sectionID)
	if timetable == nil {
		log.Printf("Error generating timetable")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate timetable"})
		return
	}

	// Return the generated timetable
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, timetable)
}
