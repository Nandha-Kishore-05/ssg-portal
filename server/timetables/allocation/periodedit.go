package allocation

import (
	"log"
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"

	"github.com/gin-gonic/gin"
)

func UpdateAllocation(c *gin.Context) {
	var req models.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Periods > 0 {
		stmt := "UPDATE subjects SET periods = ? WHERE id = ?"
		_, err := config.Database.Exec(stmt, req.Periods, req.SubjectID)
		if err != nil {
			log.Printf("Error updating periods: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating periods"})
			return
		}
	}

	if req.SubjectName != "" && req.OldSubjectName != "" {

		_, err := config.Database.Exec(`
			UPDATE subjects 
			SET name = ? 
			WHERE id = ?`, req.SubjectName, req.SubjectID)
		if err != nil {
			log.Printf("Error updating subjects: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating subjects"})
			return
		}

		_, err = config.Database.Exec(`
			UPDATE timetable 
			SET subject_name = ? 
			WHERE subject_name = ? AND department_id = ? AND semester_id = ?`,
			req.SubjectName, req.OldSubjectName, req.DepartmentID, req.SemesterID)
		if err != nil {
			log.Printf("Error updating timetable: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating timetable"})
			return
		}

		_, err = config.Database.Exec(`
			UPDATE timetable_skips 
			SET subject_name = ? 
			WHERE subject_name = ? AND department_id = ? AND semester_id = ?`,
			req.SubjectName, req.OldSubjectName, req.DepartmentID, req.SemesterID)
		if err != nil {
			log.Printf("Error updating timetable_skips: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating timetable_skips"})
			return
		}
	}

	

	if req.FacultyName != "" && req.OldFacultyName != "" {
		log.Printf("Updating faculty: name=%s, id=%d", req.FacultyName, req.FacultyID)

		_, err := config.Database.Exec(`
			UPDATE faculty 
			SET name = ? 
			WHERE id = ?`, req.FacultyName, req.FacultyID)
		if err != nil {
			log.Printf("Error updating faculty: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating faculty"})
			return
		}

		result, err := config.Database.Exec(`
	UPDATE timetable 
	SET faculty_name = ? 
	WHERE faculty_name = ? AND department_id = ? AND semester_id = ?`,
			req.FacultyName, req.OldFacultyName, req.DepartmentID, req.SemesterID)
		if err != nil {
			log.Printf("Error updating timetable: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating timetable"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error fetching rows affected: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching rows affected"})
			return
		}

		if rowsAffected == 0 {
			log.Printf("No rows were updated in the timetable. Check if the old faculty name and other details match exactly.")
		}

		log.Printf("Updating timetable_skips faculty: name=%s, old_name=%s, department_id=%d, semester_id=%d", req.FacultyName, req.OldFacultyName, req.DepartmentID, req.SemesterID)
		_, err = config.Database.Exec(`
			UPDATE timetable_skips 
			SET faculty_name = ? 
			WHERE faculty_name = ? AND department_id = ? AND semester_id = ?`,
			req.FacultyName, req.OldFacultyName, req.DepartmentID, req.SemesterID)
		if err != nil {
			log.Printf("Error updating timetable_skips: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating timetable_skips"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
}
