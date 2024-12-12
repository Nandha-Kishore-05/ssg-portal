package manualentry

import (
	"database/sql"
	"fmt"
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type BulkEntry struct {
	DayName        string   `json:"day_name"`
	Period         []int    `json:"period"`
	Classroom      string   `json:"classroom"`
	SemesterID     int      `json:"semester_id"`
	DepartmentName []string `json:"department_name"`
	SubjectName    string   `json:"subject_name"`
	FacultyName    string   `json:"faculty_name"`
	SubjectType    string   `json:"subject_type"`
	AcademicYear   int      `json:"academic_year"`
	CourseCode     string   `json:"course_code"`
	SectionName    string   `json:"section"`
}

func BulkInsert(c *gin.Context) {
	var payload struct {
		Entries []BulkEntry `json:"entries"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `INSERT INTO timetable_skips (day_name, start_time, end_time, classroom, semester_id, department_id, subject_name, faculty_name, status, academic_year, course_code, section_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := config.Database.Prepare(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare SQL statement"})
		return
	}
	defer stmt.Close()

	for _, entry := range payload.Entries {
		for _, department := range entry.DepartmentName {
			departmentID, err := getDepartmentID(department)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department ID for department: " + department})
				return
			}

			sectionID, err := getSectionID(entry.SectionName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch section ID for section: " + entry.SectionName})
				return
			}

			status, err := getStatusFromSubjectType(entry.SubjectType)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject type: " + entry.SubjectType})
				return
			}

			for _, period := range entry.Period {
				periodSlots, err := getDynamicPeriodSlots(period)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid period: %d", period)})
					return
				}

				for _, slot := range periodSlots {
					_, err = stmt.Exec(
						entry.DayName,
						slot.StartTime,
						slot.EndTime,
						entry.Classroom,
						entry.SemesterID,
						departmentID,
						entry.SubjectName,
						entry.FacultyName,
						status,
						entry.AcademicYear,
						entry.CourseCode,
						sectionID,
					)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into the database"})
						return
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bulk insert successful"})
}



// getDynamicPeriodSlots dynamically generates start and end times for the given period(s).
func getDynamicPeriodSlots(periods int) ([]struct{ StartTime, EndTime string }, error) {
	periodTimes := map[int][2]string{
		1: {"08:45:00", "09:35:00"},
		2: {"09:35:00", "10:25:00"},
		3: {"10:40:00", "11:30:00"},
		4: {"13:45:00", "14:35:00"},
		5: {"14:35:00", "15:25:00"},
		6: {"15:40:00", "16:30:00"},
	}

	// Assume "periods" is a bitmask or list of selected periods (e.g., [3, 4])
	var slots []struct{ StartTime, EndTime string }

	// Example: Handle single periods or ranges dynamically
	if times, exists := periodTimes[periods]; exists {
		// Single period
		slots = append(slots, struct{ StartTime, EndTime string }{times[0], times[1]})
	} else {
		// Handle multiple or range periods
		for period, times := range periodTimes {
			if period <= periods {
				slots = append(slots, struct{ StartTime, EndTime string }{times[0], times[1]})
			}
		}
	}

	if len(slots) == 0 {
		return nil, fmt.Errorf("no valid time slots found for periods: %d", periods)
	}

	return slots, nil
}

func getDepartmentID(departmentName string) (int, error) {
	var departmentID int
	query := "SELECT id FROM departments WHERE name = ?"
	err := config.Database.QueryRow(query, departmentName).Scan(&departmentID)
	if err == sql.ErrNoRows {
		return 0, sql.ErrNoRows
	} else if err != nil {
		return 0, err
	}
	return departmentID, nil
}

func getSectionID(sectionName string) (int, error) {
	var sectionID int
	query := "SELECT id FROM master_section WHERE section_name = ?"
	err := config.Database.QueryRow(query, sectionName).Scan(&sectionID)
	if err != nil {
		return 0, err
	}
	return sectionID, nil
}
func getStatusFromSubjectType(subjectType string) (int, error) {
	subjectTypeMap := map[string]int{
		"Lab Subject":     1,
		"Non Lab Subject": 2,
		"Elective 3":      3,
		"Elective 4":      4,
		"Elective 5":      5,
		"Open Elective":   6,
		"Add On Course":   7,
		"Honor":           8,
		"Minor":           9,
		"Elective 1":      10,
	}

	status, exists := subjectTypeMap[subjectType]
	if !exists {
		return 0, fmt.Errorf("unknown subject type: %s", subjectType)
	}
	return status, nil
}
