package timetables

import (
	"fmt"
	"log"
	"net/http"
	"ssg-portal/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type Allocation struct {
	Date       string `json:"date"`
	Period     string `json:"period"`
	Venue      string `json:"venue"`
	Subject    string `json:"subject"`
	CourseCode string `json:"course_code"`
	Section    string `json:"section"`
	Faculty    string `json:"faculty"`
}

type DepartmentAllocation struct {
	Department  string       `json:"department"`
	Allocations []Allocation `json:"allocations"`
}

func HandleExcelUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	filePath := fmt.Sprintf("./%s", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Parse the Excel file and generate timetable allocations
	allocations, err := ParseExcelAndGenerateTimetable(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allocations)
}

// Parse the Excel file and generate timetable allocations
func ParseExcelAndGenerateTimetable(filePath string) ([]DepartmentAllocation, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	// Map to store allocations per department
	departmentAllocations := make(map[string][]Allocation)

	// Example: Generate working days, periods, and venues
	workingDays := generateWorkingDays()
	periods := [][2]int{{1, 2}, {3, 4}, {5, 6}}
	venues := []string{"CSE LAB 1", "CSE LAB 2", "CSE LAB 3", "CSE LAB 4", "CSE LAB 5", "CSE LAB 6"}

	// Skip the header row and process data
	for _, row := range rows[1:] {
		if len(row) < 4 {
			continue
		}

		department, subject, faculty, section := row[0], row[1], row[2], row[3] // Extract department, subject, faculty, and section

		// Allocate timetable for each department, subject, faculty, and section
		allocations := allocateTimetable(department, subject, faculty, section, workingDays, periods, venues)
		departmentAllocations[department] = append(departmentAllocations[department], allocations...)
	}

	// Convert map to slice for response
	var result []DepartmentAllocation
	for dept, allocations := range departmentAllocations {
		result = append(result, DepartmentAllocation{
			Department:  dept,
			Allocations: allocations,
		})
	}

	return result, nil
}

// Define restricted dates and periods to skip
var blockedDatesAndPeriods = map[string][]int{
	"2024-12-16": {3, 5},
	"2024-12-17": {3, 5},
	"2024-12-18": {3, 5},
	"2024-12-19": {3, 5},
	"2024-12-20": {3, 5},
	"2024-12-21": {3, 5},
}

func isBlocked(date string, period int) bool {
	if blockedPeriods, exists := blockedDatesAndPeriods[date]; exists {
		for _, blockedPeriod := range blockedPeriods {
			if period == blockedPeriod {
				return true
			}
		}
	}
	return false
}

func allocateTimetable(department, subject, faculty, section string, workingDays []time.Time, periods [][2]int, venues []string) []Allocation {
	var allocations []Allocation
	sectionAllocationsCount := make(map[string]int)

	for _, day := range workingDays {
		dateStr := day.Format("2006-01-02")
		if sectionAllocationsCount[section] >= 6 {
			log.Printf("Section %s has reached maximum allocation on %s", section, dateStr)
			continue
		}

		for _, period := range periods {
			if isBlocked(dateStr, period[0]) || isBlocked(dateStr, period[1]) {
				log.Printf("Skipping blocked period (%d, %d) on %s", period[0], period[1], dateStr)
				continue
			}

			for _, venue := range venues {
				if sectionAllocationsCount[section] >= 6 {
					continue
				}

				allocation := Allocation{
					Date:       dateStr,
					Period:     fmt.Sprintf("(%d,%d)", period[0], period[1]),
					Venue:      venue,
					Subject:    subject,
					CourseCode: "22CS402",
					Section:    section,
					Faculty:    faculty,
				}

				if checkConflict(day, period[0], period[1], venue, faculty, subject, department) {
					log.Printf("Conflict detected for faculty %s or venue %s on %s during (%d, %d)", faculty, venue, dateStr, period[0], period[1])
					continue
				}

				saveAllocation(department, allocation, faculty)
				sectionAllocationsCount[section]++
				allocations = append(allocations, allocation)
				log.Printf("Allocation added: %+v", allocation)
			}
		}
	}

	return allocations
}

// Check if there is any conflict for the given allocation
func checkConflict(day time.Time, start, end int, venue, faculty, subject, department string) bool {
	query := `
    SELECT COUNT(*) FROM allocation 
    WHERE date = ? AND (venue = ? OR faculty_name = ?) AND 
    ((period_start BETWEEN ? AND ?) OR (period_end BETWEEN ? AND ?)) 
`

	var count int
	err := config.Database.QueryRow(query, day, venue, faculty, start, end, start, end).Scan(&count)
	if err != nil {
		log.Println("Error checking conflict:", err)
		return true
	}
	return count > 0
}

// Save the allocation to the database
func saveAllocation(department string, allocation Allocation, faculty string) {
	query := `
    INSERT INTO allocation (department, date, period_start, period_end, venue, course_name, course_code, section, faculty_name)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err := config.Database.Exec(query, department, allocation.Date, allocation.Period[1]-'0', allocation.Period[3]-'0', allocation.Venue, allocation.Subject, allocation.CourseCode, allocation.Section, faculty)
	if err != nil {
		log.Println("Error saving allocation:", err)
	}
}

// Generate the list of working days (for example purposes)
func generateWorkingDays() []time.Time {
	dates := []string{
		"2024-12-16", "2024-12-17", "2024-12-18", "2024-12-19", "2024-12-20",
		"2024-12-21",
	}

	var workingDays []time.Time
	for _, date := range dates {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err == nil {
			workingDays = append(workingDays, parsedDate)
		} else {
			log.Printf("Error parsing date: %v", date)
		}
	}
	return workingDays
}
