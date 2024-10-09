package studententry

import (
	"log"
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// GetStudentCourseMappingsHandler handles the API request and returns student course mappings
func GetStudentOptions(c *gin.Context) {
	query := `
	SELECT DISTINCT scm.student_id, s.student_name,S.roll_no, scm.department_id, d.name as department_name,
		   scm.semester_id, sem.semester_name, scm.academic_year_id, ay.academic_year
	FROM student_course_mapping scm
	JOIN student s ON scm.student_id = s.id
	JOIN departments d ON scm.department_id = d.id
	JOIN semester sem ON scm.semester_id = sem.id
	JOIN master_academic_year ay ON scm.academic_year_id = ay.id
	`

	// Execute the query
	rows, err := config.Database.Query(query)
	if err != nil {
		log.Println("Error querying student_course_mapping:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch data"})
		return
	}
	defer rows.Close()

	// Create a slice to hold the mappings
	var mappings []models.StudentOptions

	// Iterate over the result rows and scan the data into the slice
	for rows.Next() {
		var mapping models.StudentOptions
		err := rows.Scan(
			&mapping.StudentID, &mapping.StudentName, &mapping.StudentRollNo, &mapping.DepartmentID, &mapping.DepartmentName,
			&mapping.SemesterID, &mapping.SemesterName,
			&mapping.AcademicYearID, &mapping.AcademicYear,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing data"})
			return
		}
		mappings = append(mappings, mapping)
	}

	// Check for any error after iterating through rows
	if err = rows.Err(); err != nil {
		log.Println("Error after row iteration:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing data"})
		return
	}

	// Return the response with the data
	c.JSON(http.StatusOK, mappings)
}
