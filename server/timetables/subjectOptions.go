package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

type SubjectFilter struct {
	AcademicYearID int `json:"academic_year_id"`
	SemesterID     int `json:"semester_id"`
	DepartmentID   int `json:"department_id"`
	SectionID      int `json:"section_id"`
}
func SubjectOptions(c *gin.Context) {
    // Bind the JSON body to a single SubjectFilter struct
    var filter SubjectFilter
    if err := c.ShouldBindJSON(&filter); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    // Query the database using the current filter
    query := `
        SELECT s.id, s.name 
        FROM subjects s
        JOIN faculty_subjects fs ON fs.subject_id = s.id
        WHERE fs.academic_year_id = ? AND fs.semester_id = ? 
          AND fs.department_id = ? AND fs.section_id = ?
    `
    rows, err := config.Database.Query(query, filter.AcademicYearID, filter.SemesterID, filter.DepartmentID, filter.SectionID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var subjectOptions []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        subjectOptions = append(subjectOptions, map[string]interface{}{
            "label": name,
            "value": id,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, subjectOptions)
}


func CourseCodeOptions(c *gin.Context) {
    // Get the subject name from the query parameter
    subjectName := c.DefaultQuery("subject_name", "")

    if subjectName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Subject name is required"})
        return
    }

    // Query the database to get the course code for the selected subject name
    rows, err := config.Database.Query("SELECT id, course_code FROM subjects WHERE name = ?", subjectName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var courseCodeOptions []map[string]interface{}
    for rows.Next() {
        var id int
        var courseCode string
        if err := rows.Scan(&id, &courseCode); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        courseCodeOptions = append(courseCodeOptions, map[string]interface{}{
            "label": courseCode,
            "value": courseCode,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, courseCodeOptions)
}
