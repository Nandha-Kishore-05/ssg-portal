package studententry

import (
	"database/sql"
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)

// InsertStudentEntries handles the insertion of student, course, and mapping entries
func InsertStudentEntries(c *gin.Context) {
	var req models.StudentEntryRequest

	// Parse the JSON request body
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Query to check if student already exists
	studentSelectQuery := `SELECT id FROM student WHERE roll_no = ?`

	// Query to insert data into student table
	studentInsertQuery := `
        INSERT INTO student (student_name, roll_no) 
        VALUES (?, ?)
    `

	// Query to select course ID if it exists
	courseSelectQuery := `SELECT id FROM student_courses WHERE course_code = ?`

	// Query to insert data into student_courses table
	courseInsertQuery := `
        INSERT INTO student_courses (course_name, course_code)
        VALUES (?, ?)
    `

	// Query to insert data into student_course_mapping table with multiple course IDs
	mappingInsertQuery := `
        INSERT INTO student_course_mapping (student_id, course_code, department_id, semester_id, academic_year_id)
        VALUES (?, ?, ?, ?, ?)
    `
	for _, student := range req.Students {
		var studentID int64

		// Check if the student already exists based on roll number
		err := config.Database.QueryRow(studentSelectQuery, student.RollNumber).Scan(&studentID)
		if err == sql.ErrNoRows {
			// Insert student data if not exists
			result, err := config.Database.Exec(studentInsertQuery, student.StudentName, student.RollNumber)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert student data", "details": err.Error()})
				return
			}
			studentID, _ = result.LastInsertId()
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student data", "details": err.Error()})
			return
		}

		// Now, use the course details directly from the student struct
		var courseID int64

		// Check if the course already exists based on course code
		err = config.Database.QueryRow(courseSelectQuery, student.CourseCode).Scan(&courseID)
		if err == sql.ErrNoRows {
			// Insert course data if not exists
			result, err := config.Database.Exec(courseInsertQuery, student.CourseName, student.CourseCode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert course data", "details": err.Error()})
				return
			}
			courseID, _ = result.LastInsertId()
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve course data", "details": err.Error()})
			return
		}

		// Insert student-course mapping into student_course_mapping table
		_, err = config.Database.Exec(mappingInsertQuery, studentID, student.CourseCode, req.Department, req.Semester, req.AcademicYear)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert student-course mapping", "details": err.Error()})
			return
		}
	}
	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}
