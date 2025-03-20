package subjectUpload

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

type SubjectData struct {
	CourseCode  string `json:"Course Code"`
	CourseName  string `json:"Course Name"`
	FacultyID   string `json:"Faculty ID"`
	FacultyName string `json:"Faculty NAME"`
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
		status := "1"
		if subject.LabSubject == "YES" {
			status = "0"
		}
		query := "SELECT id FROM subjects WHERE name = ? AND  course_code = ?"
		err := config.Database.QueryRow(query, subject.CourseName, subject.CourseCode).Scan(&subjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Subject not found", "Course Code": subject.CourseCode})
			return
		}

		facultyQuery := "SELECT id FROM faculty WHERE  faculty_id = ?"
		err = config.Database.QueryRow(facultyQuery, subject.FacultyID).Scan(&facultyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Faculty not found", "Faculty Name": subject.FacultyName, "Faculty ID": subject.FacultyID})
			return
		}

		insertQuery := `INSERT INTO faculty_subjects (faculty_id, subject_id, semester_id, academic_year_id, department_id, section_id, status, periods)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = config.Database.Exec(insertQuery, facultyID, subjectID, request.Semester, request.AcademicYear, request.Department, sectionID, status, subject.Periods)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into faculty_subjects", "Subject ID": subjectID})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}
