package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

type SubjectData struct {
	CourseCode  string `json:"Course Code"`
	CourseName  string `json:"Course Name"`
	FacultyID   string `json:"Faculty ID"`
	FacultyName string `json:"Faculty Name"`
	LabSubject  string `json:"Lab-Subject"`
	Periods     int    `json:"Periods"`
}

type FacultySubjectsRequest struct {
	AcademicYear int           `json:"academicYear"`
	Classroom    int           `json:"classroom"`
	Department   int           `json:"department"`
	Semester     int           `json:"semester"`
	SubjectData  []SubjectData `json:"subjectData"`
}

func UploadDetails(c *gin.Context) {
	var request FacultySubjectsRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Retrieve the section_id from classrooms table based on the provided classroom_id
	var sectionID int
	sectionQuery := "SELECT section_id FROM classrooms WHERE id = ?"
	err := config.Database.QueryRow(sectionQuery, request.Classroom).Scan(&sectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Section not found for the given classroom"})
		return
	}

	for _, subject := range request.SubjectData {
		var subjectID int
		var facultyID int

		// Step 2: Retrieve the subject_id from subjects table
		query := "SELECT id FROM subjects WHERE course_code = ? AND name = ?"
		err := config.Database.QueryRow(query, subject.CourseCode, subject.CourseName).Scan(&subjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Subject not found"})
			return
		}

		// Step 3: Retrieve the faculty_id from faculty table
		facultyQuery := "SELECT id FROM faculty WHERE name = ? AND faculty_id = ?"
		err = config.Database.QueryRow(facultyQuery, subject.FacultyName, subject.FacultyID).Scan(&facultyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Faculty not found"})
			return
		}

		// Step 4: Determine the status based on Lab-Subject
		status := "1" // Default status is "1" for "NO"
		if subject.LabSubject == "YES" {
			status = "0" // Change status to "0" for "YES"
		}

		// Step 5: Update the status and periods in the subjects table
		updateQuery := "UPDATE subjects SET status = ?, periods = ? WHERE id = ?"
		if _, err := config.Database.Exec(updateQuery, status, subject.Periods, subjectID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subjects table"})
			return
		}

		// Step 6: Insert data into faculty_subjects table
		insertQuery := `INSERT INTO faculty_subjects (faculty_id, subject_id, semester_id, academic_year_id, department_id, section_id, status, periods)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = config.Database.Exec(insertQuery, facultyID, subjectID, request.Semester, request.AcademicYear, request.Department, sectionID, status, subject.Periods)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into faculty_subjects"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}
