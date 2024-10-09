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

		courseCode, ok := row["CourseCode"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid CourseCode data")
			return
		}

		venue, ok := row["Venue"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid venue data")
			return
		}

		academicYear, ok := row["AcademicYear"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Invalid academic year data")
			return
		}

		deptID, err := getDepartmentID(department)
		if err != nil {
			c.String(http.StatusInternalServerError, "Department ID error: %s", err.Error())
			return
		}


		facultyID, err := getFacultyID(faculty, facultyIDStr, deptID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Faculty ID error: %s", err.Error())
			return
		}


		subjectID, err := getOrCreateSubject(subject, deptID, semester, status, periods,courseCode)
		if err != nil {
			c.String(http.StatusInternalServerError, "Subject insertion error: %s", err.Error())
			return
		}


		academicYearID, err := getOrCreateAcademicYear(academicYear)
		if err != nil {
			c.String(http.StatusInternalServerError, "Academic year insertion error: %s", err.Error())
			return
		}

	
		_, err = config.Database.Exec(`INSERT INTO faculty_subjects (faculty_id, subject_id, semester_id, academic_year_id) VALUES (?, ?, ?, ?)`,
			facultyID, subjectID, semester, academicYearID)
		if err != nil {
			c.String(http.StatusInternalServerError, "faculty_subjects insertion error: %s", err.Error())
			return
		}


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

	err := config.Database.QueryRow("SELECT id FROM faculty WHERE faculty_id = ?", facultyIDStr).Scan(&id)
	if err == sql.ErrNoRows {
	
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

func getOrCreateSubject(name string, deptID int, semester int, status string, periods int, courseCode string) (int, error) {
    var id int

    err := config.Database.QueryRow(
        "SELECT id FROM subjects WHERE name = ? AND department_id = ? AND semester_id = ? AND status = ? AND course_code = ?",
        name, deptID, semester, status, courseCode,
    ).Scan(&id)

    if err == sql.ErrNoRows {

        res, err := config.Database.Exec(
            "INSERT INTO subjects (name, department_id, semester_id, status, periods, course_code) VALUES (?, ?, ?, ?, ?, ?)",
            name, deptID, semester, status, periods, courseCode,
        )
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

func getOrCreateAcademicYear(year string) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM master_academic_year WHERE academic_year = ? ", year).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := config.Database.Exec("INSERT INTO master_academic_year (academic_year) VALUES (?)", year)
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
