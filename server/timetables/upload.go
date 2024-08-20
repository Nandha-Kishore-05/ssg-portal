
package timetables

import (
	"database/sql"
	"net/http"
	"ssg-portal/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Uploaddetails(c *gin.Context) {
	// Assuming the data is received as a JSON array
	var records []map[string]interface{}
	if err := c.BindJSON(&records); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data: %s", err.Error())
		return
	}

	for _, row := range records {
		// Extract data from the row

		department, ok := row["Department"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid department data")
			return
		}

		faculty, ok := row["Faculty"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid faculty data")
			return
		}

		// Handle the "Lab-subject" field
		status := "1" // Default status to "no"
		if labSubject, ok := row["Lab-subject"].(string); ok && labSubject == "YES" {
			status = "0"
		}

		// Handle "Periods" field, which might be a number (float64)
		var periods int
		switch v := row["Periods"].(type) {
		case string:
			periods, _ = strconv.Atoi(v)
		case float64:
			periods = int(v)
		default:
			c.String(http.StatusBadRequest, "Invalid periods data")
			return
		}

		// Handle "Semester" field, which might also be a number (float64)
		var semester int
		switch v := row["Semester"].(type) {
		case string:
			semester, _ = strconv.Atoi(v)
		case float64:
			semester = int(v)
		default:
			c.String(http.StatusBadRequest, "Invalid semester data")
			return
		}

		subject, ok := row["Subject"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid subject data")
			return
		}

		venue, ok := row["Venue"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid venue data")
			return
		}

		// Get Department ID
		deptID, err := getDepartmentID(department)
		if err != nil {
			c.String(http.StatusInternalServerError, "Department ID error: %s", err.Error())
			return
		}

		// Get Faculty ID
		facultyID, err := getFacultyID(faculty, deptID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Faculty ID error: %s", err.Error())
			return
		}

		// Insert into subjects table
		subjectID, err := getOrCreateSubject(subject, deptID, semester, status, periods)
		if err != nil {
			c.String(http.StatusInternalServerError, "Subject insertion error: %s", err.Error())
			return
		}

		// Insert into faculty_subjects table
		_, err = config.Database.Exec(`INSERT INTO faculty_subjects (faculty_id, subject_id, semester_id) VALUES (?, ?, ?)`,
			facultyID, subjectID, semester)
		if err != nil {
			c.String(http.StatusInternalServerError, "faculty_subjects insertion error: %s", err.Error())
			return
		}

		// Insert into classrooms (venue)
		_, err = getOrCreateVenue(venue, deptID,semester)
		if err != nil {
			c.String(http.StatusInternalServerError, "Classroom insertion error: %s", err.Error())
			return
		}
	}

	c.String(http.StatusOK, "File processed and data inserted successfully!")
}

func getDepartmentID(name string) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM departments WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := config.Database.Exec("INSERT INTO departments (name) VALUES (?)", name)
		if err != nil {
			return 0, err
		}
		id64, _ := res.LastInsertId()
		id = int(id64)
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func getFacultyID(name string, deptID int) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM faculty WHERE name = ? AND department_id = ?", name, deptID).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := config.Database.Exec("INSERT INTO faculty (name, department_id) VALUES (?, ?)", name, deptID)
		if err != nil {
			return 0, err
		}
		id64, _ := res.LastInsertId()
		id = int(id64)
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func getOrCreateSubject(name string, deptID int, semester int, status string, periods int) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM subjects WHERE name = ? AND department_id = ? AND semester_id = ?", name, deptID, semester).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := config.Database.Exec("INSERT INTO subjects (name, department_id, semester_id, status, periods) VALUES (?, ?, ?, ?, ?)",
			name, deptID, semester, status, periods)
		if err != nil {
			return 0, err
		}
		id64, _ := res.LastInsertId()
		id = int(id64)
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func getOrCreateVenue(name string, departmentID int,semesterID int) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM classrooms WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		// Classroom does not exist, so create a new one with the correct department ID
		res, err := config.Database.Exec("INSERT INTO classrooms (name, department_id,semester_id) VALUES (?, ?,?)", name, departmentID,semesterID)
		if err != nil {
			return 0, err
		}
		id64, _ := res.LastInsertId()
		id = int(id64)
	} else if err != nil {
		return 0, err
	}
	return id, nil
}
