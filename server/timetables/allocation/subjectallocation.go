package allocation

import (
	"database/sql"
	"log"
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)
func Subjectallocation(c *gin.Context) {
    query := `
    (
        SELECT 
            DISTINCT 
            subjects.id AS subject_id,
            subjects.name AS subject_name,
            departments.name AS department_name,
            semester.semester_name AS semester_name,
            subjects.periods,
            CASE 
                WHEN subjects.status = '0' THEN 'Lab Subject'
                WHEN subjects.status = '1' THEN 'Non-Lab Subject'
            END AS STATUS,
            faculty.name AS faculty_name,
            faculty.id AS faculty_id,  
            subjects.department_id, 
            subjects.semester_id
        FROM 
            subjects
        LEFT JOIN 
            faculty_subjects ON subjects.id = faculty_subjects.subject_id
        LEFT JOIN 
            faculty ON faculty_subjects.faculty_id = faculty.id
        LEFT JOIN 
            departments ON subjects.department_id = departments.id
        LEFT JOIN 
            semester ON subjects.semester_id = semester.id
        WHERE 
            subjects.name IS NOT NULL
            AND departments.name IS NOT NULL
            AND semester.semester_name IS NOT NULL
            AND subjects.periods IS NOT NULL
            AND subjects.status IS NOT NULL
    )
    UNION ALL
    (
        SELECT 
            DISTINCT 
            timetable_skips.id AS subject_id,  
            timetable_skips.subject_name AS subject_name,
            departments.name AS department_name,
            semester.semester_name AS semester_name,
            '1' AS periods,  
            'Non-Lab Subject' AS STATUS,  
            timetable_skips.faculty_name AS faculty_name,
            NULL AS faculty_id,  
            timetable_skips.department_id, 
            timetable_skips.semester_id
        FROM 
            timetable_skips
        LEFT JOIN 
            departments ON timetable_skips.department_id = departments.id
        LEFT JOIN 
            semester ON timetable_skips.semester_id = semester.id
        WHERE 
            timetable_skips.subject_name IS NOT NULL
            AND departments.name IS NOT NULL
            AND semester.semester_name IS NOT NULL
    )
    ORDER BY 
        department_id, semester_id;
    `

    rows, err := config.Database.Query(query)
    if err != nil {
        log.Println("Error executing query: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
        return
    }
    defer rows.Close()

    var subjects []models.SubjectInfo

    for rows.Next() {
        var subject models.SubjectInfo
        var facultyID sql.NullInt64

        err := rows.Scan(
            &subject.SubjectID,
            &subject.SubjectName,
            &subject.DepartmentName,
            &subject.SemesterName,
            &subject.Periods,
            &subject.Status,
            &subject.FacultyName,
            &facultyID,           
            &subject.DepartmentID,
            &subject.SemesterID,
        )
        if err != nil {
            log.Println("Error scanning row: ", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Data parsing error"})
            return
        }

        if facultyID.Valid {
            subject.FacultyID = int(facultyID.Int64)  
        } else {
            subject.FacultyID = 0  
        }

        subjects = append(subjects, subject)
    }

    if err = rows.Err(); err != nil {
        log.Println("Row iteration error: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Data retrieval error"})
        return
    }

    c.JSON(http.StatusOK, subjects)
}