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

	// Parse the Excel file
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
    workingDays := generateWorkingDays()
    periods := [][2]int{{1, 2}, {3, 4}, {5, 6}, {6, 7}}
    venues := []string{"CSE LAB 1", "CSE LAB 2", "CSE LAB 3", "CSE LAB 4", "CSE LAB 5"}

    // Skip the header row and process data
    for _, row := range rows[1:] {
        if len(row) < 3 {
            continue
        }

        department, subject, faculty := row[0], row[1], row[2]  // Extract faculty name
        sectionCount := getSectionCount(department)  // Get the section count
        allocations := allocateTimetable(department, subject, faculty, workingDays, periods, venues, sectionCount)
        departmentAllocations[department] = append(departmentAllocations[department], allocations...)
    }

    // Convert map to slice
    var result []DepartmentAllocation
    for dept, allocations := range departmentAllocations {
        result = append(result, DepartmentAllocation{
            Department:  dept,
            Allocations: allocations,
        })
    }

    return result, nil
}

// Get the number of sections for the department
func getSectionCount(department string) int {
    departmentSections := map[string]int{
        "CSE": 4,
        "ISE": 1,
        "CSD": 1,
        "IT": 3,
        "FT": 1,
        "ECE": 4,
        "CT": 1,
        "AIDS": 3,
        "FD": 1,
        "AGRI": 1,
        "EIE": 1,
        "BME": 1,
        "AIML": 2,
        "EEE": 1,
        "CSBS": 1,
        "MTRS": 1,
        "MECH": 1,
        "BT": 2,
        "CIVIL": 1,
    }

    if sections, exists := departmentSections[department]; exists {
        return sections
    }
    return 1 // Default: 1 section if the department is not found
}

func allocateTimetable(department, subject, faculty string, workingDays []time.Time, periods [][2]int, venues []string, sectionCount int) []Allocation {
    var allocations []Allocation

    for _, day := range workingDays {
        for _, period := range periods {
            for _, venue := range venues {
                for section := 1; section <= sectionCount; section++ {
                    if checkConflict(day, period[0], period[1], venue, faculty, subject, department) {
                        continue
                    }

                    allocation := Allocation{
                        Date:       day.Format("2006-01-02"),
                        Period:     fmt.Sprintf("(%d,%d)", period[0], period[1]),
                        Venue:      venue,
                        Subject:    subject,
                        CourseCode: "22CS402", // Example hardcoded, replace with logic
                        Section:    fmt.Sprintf("Section %d", section),
                        Faculty:    faculty,  // Set faculty here
                    }

                    saveAllocation(department, allocation, faculty)
                    allocations = append(allocations, allocation)
                    break
                }
                if len(allocations) > 0 {
                    break
                }
            }
        }
    }

    return allocations
}

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

func saveAllocation(department string, allocation Allocation, faculty string) {
	query := `
		INSERT INTO  allocation (department, date, period_start, period_end, venue, course_name, course_code, section, faculty_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := config.Database.Exec(query, department, allocation.Date, allocation.Period[1]-'0', allocation.Period[3]-'0', allocation.Venue, allocation.Subject, allocation.CourseCode, allocation.Section, faculty)
	if err != nil {
		log.Println("Error saving allocation:", err)
	}
}

func generateWorkingDays() []time.Time {
	dates := []string{
		"2024-12-16", "2024-12-17", "2024-12-18", "2024-12-19", "2024-12-20", 
		"2024-12-21", "2024-12-23", "2024-12-24", "2024-12-26", "2024-12-27", 
		"2024-12-28", "2025-01-03", "2025-01-04", "2025-01-06", "2025-01-07", 
		"2025-01-08", "2025-01-09", "2025-01-11", "2025-01-20", "2025-01-21", 
		"2025-01-22",
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
