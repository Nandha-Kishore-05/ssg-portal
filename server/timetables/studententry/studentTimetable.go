package studententry

import (
	"log"
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)

func StudentTimetable(c *gin.Context) {

	studentID := c.Param("studentID")

	query := `
	SELECT 
    t.id,
    t.day_name,
    t.start_time,
    t.end_time,
    t.classroom,
    t.subject_name,
    t.faculty_name,
    t.status,
    scm.student_id,
    scm.course_code,
    scm.academic_year_id,
    scm.department_id,
    scm.semester_id
FROM 
    timetable t
JOIN 
    student_course_mapping scm ON t.course_code = scm.course_code 
                               AND t.academic_year = scm.academic_year_id
                               AND t.department_id = scm.department_id
                               AND t.semester_id = scm.semester_id

JOIN 
    master_academic_year may ON t.academic_year = may.id
WHERE 
    scm.student_id = ?

	`

	rows, err := config.Database.Query(query, studentID)
	if err != nil {
		log.Println("Error querying student timetable:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch data"})
		return
	}
	defer rows.Close()

	var timetable []models.StudentTimetable
	for rows.Next() {
		var entry models.StudentTimetable
		err := rows.Scan(
			&entry.ID,
			&entry.DayName,
			&entry.StartTime,
			&entry.EndTime,
			&entry.Classroom,
			&entry.SubjectName,
			&entry.FacultyName,
			&entry.Status,
			&entry.StudentID,
			&entry.CourseCode,
			&entry.AcademicYearID,
			&entry.DepartmentID,
			&entry.SemesterID,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing data"})
			return
		}
		timetable = append(timetable, entry)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error after row iteration:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing data"})
		return
	}

	c.JSON(http.StatusOK, timetable)
}
