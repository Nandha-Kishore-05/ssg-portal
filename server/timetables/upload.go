package timetables

import (
	"database/sql"
	"net/http"
	"ssg-portal/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Uploaddetails(c *gin.Context) {
	var records []map[string]interface{}
	if err := c.BindJSON(&records); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data: %s", err.Error())
		return
	}

	for _, row := range records {
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

		// Extract faculty_id from the row (if available)
		facultyIDStr, ok := row["FacultyID"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid faculty ID data")
			return
		}

		status := "1"
		if labSubject, ok := row["Lab-subject"].(string); ok && labSubject == "YES" {
			status = "0"
		}

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

		// Extract academic year from the row
		academicYear, ok := row["AcademicYear"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid academic year data")
			return
		}

		// Get or create department ID
		deptID, err := getDepartmentID(department)
		if err != nil {
			c.String(http.StatusInternalServerError, "Department ID error: %s", err.Error())
			return
		}

		// Get or create faculty ID, now including faculty_id
		facultyID, err := getFacultyID(faculty, facultyIDStr, deptID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Faculty ID error: %s", err.Error())
			return
		}

		// Get or create subject ID
		subjectID, err := getOrCreateSubject(subject, deptID, semester, status, periods)
		if err != nil {
			c.String(http.StatusInternalServerError, "Subject insertion error: %s", err.Error())
			return
		}

		// Get or create academic year ID
		academicYearID, err := getOrCreateAcademicYear(academicYear, deptID, semester)
		if err != nil {
			c.String(http.StatusInternalServerError, "Academic year insertion error: %s", err.Error())
			return
		}

		// Insert into faculty_subjects
		_, err = config.Database.Exec(`INSERT INTO faculty_subjects (faculty_id, subject_id, semester_id, academic_year_id) VALUES (?, ?, ?, ?)`,
			facultyID, subjectID, semester, academicYearID)
		if err != nil {
			c.String(http.StatusInternalServerError, "faculty_subjects insertion error: %s", err.Error())
			return
		}

		// Get or create venue ID
		_, err = getOrCreateVenue(venue, deptID, semester, academicYearID)
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

func getFacultyID(name string, facultyIDStr string, deptID int) (int, error) {
	var id int
	// Check if the faculty ID already exists
	err := config.Database.QueryRow("SELECT id FROM faculty WHERE faculty_id = ?", facultyIDStr).Scan(&id)
	if err == sql.ErrNoRows {
		// If not, insert a new faculty record
		res, err := config.Database.Exec("INSERT INTO faculty (name, department_id, faculty_id) VALUES (?, ?, ?)", name, deptID, facultyIDStr)
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

func getOrCreateVenue(name string, departmentID int, semesterID int, academicYearID int) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM classrooms WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := config.Database.Exec("INSERT INTO classrooms (name, department_id, semester_id,academic_year_id) VALUES (?, ?, ?,?)", name, departmentID, semesterID, academicYearID)
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

func getOrCreateAcademicYear(year string, departmentID int, semesterID int) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM academic_year WHERE academic_year = ? AND department_id = ? AND semester_id = ?", year, departmentID, semesterID).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := config.Database.Exec("INSERT INTO academic_year (academic_year, department_id, semester_id) VALUES (?, ?, ?)", year, departmentID, semesterID)
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
