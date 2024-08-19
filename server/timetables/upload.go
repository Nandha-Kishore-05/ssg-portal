package timetables

import (
	"database/sql"
	"net/http"
	"ssg-portal/config"
	"ssg-portal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func Uploaddetails(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Upload file error: %s", err.Error())
		return
	}

	f, err := file.Open()
	if err != nil {
		c.String(http.StatusBadRequest, "File open error: %s", err.Error())
		return
	}
	defer f.Close()

	xlFile, err := excelize.OpenReader(f)
	if err != nil {
		c.String(http.StatusBadRequest, "Excel parsing error: %s", err.Error())
		return
	}

	records, err := parseExcel(xlFile)
	if err != nil {
		c.String(http.StatusBadRequest, "Data parsing error: %s", err.Error())
		return
	}

	err = insertRecords(records)
	if err != nil {
		c.String(http.StatusInternalServerError, "Database insertion error: %s", err.Error())
		return
	}

	c.String(http.StatusOK, "File processed and data inserted successfully!")
}

func parseExcel(xlFile *excelize.File) ([]models.Record, error) {
	var records []models.Record

	rows, err := xlFile.GetRows(xlFile.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		if i == 0 {
			// Skip header row
			continue
		}

		// Map "Lab-subject" to Status (0 for "yes", 1 for "no")
		status := 1
		if row[3] == "yes" {
			status = 0
		}

		periods, _ := strconv.Atoi(row[4])
		semester, _ := strconv.Atoi(row[5])

		record := models.Record{
			Department: row[1], // "Department"
			Faculty:    row[2], // "Faculty"
			Status:     status, // "Lab-subject" mapped to Status
			Periods:    periods,
			Semester:   semester,
			Subject:    row[6], // "Subject"
			Venue:      row[7], // "Venue"
		}

		records = append(records, record)
	}

	return records, nil
}

func insertRecords(records []models.Record) error {
	for _, record := range records {
		deptID, err := getDepartmentID(record.Department)
		if err != nil {
			return err
		}

		facultyID, err := getFacultyID(record.Faculty, deptID)
		if err != nil {
			return err
		}

		subjectID, err := getSubjectID(record.Subject, deptID, record.Semester)
		if err != nil {
			return err
		}

		// Insert the record into the faculty_subjects table
		query := `INSERT INTO faculty_subjects (faculty_id, subject_id, semester_id) VALUES (?, ?, ?)`
		_, err = config.Database.Exec(query, facultyID, subjectID, record.Semester)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDepartmentID(name string) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM departments WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		// Insert department if not exists
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
		// Insert faculty if not exists
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

func getSubjectID(name string, deptID int, semester int) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM subjects WHERE name = ? AND department_id = ? AND semester_id = ?", name, deptID, semester).Scan(&id)
	if err == sql.ErrNoRows {
		// Insert subject if not exists
		res, err := config.Database.Exec("INSERT INTO subjects (name, department_id, semester_id) VALUES (?, ?, ?)", name, deptID, semester)
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

func getVenueID(name string) (int, error) {
	var id int
	err := config.Database.QueryRow("SELECT id FROM classrooms WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		// Insert venue if not exists
		res, err := config.Database.Exec("INSERT INTO classrooms (name, department_id) VALUES (?, 0)", name)
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
