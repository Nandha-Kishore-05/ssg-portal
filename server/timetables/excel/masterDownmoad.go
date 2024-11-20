// package excel

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"ssg-portal/config"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/xuri/excelize/v2"
// )

// func Masterdownload(c *gin.Context) {
// 	academicYearID := c.Param("academic_year_id")

// 	if academicYearID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic year parameter"})
// 		return
// 	}

// 	// Updated query to fetch all semesters belonging to the given academic_year_id
// 	query := `
// 		SELECT t.day_name, t.start_time, t.end_time, t.classroom,
//        t.subject_name, t.faculty_name, d.name AS department_name, may.academic_year, t.semester_id
// FROM timetable t
// JOIN departments d ON t.department_id = d.id

// JOIN master_academic_year may ON t.academic_year = may.id
// WHERE may.id = ?
// 	`

// 	rows, err := config.Database.Query(query, academicYearID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error: " + err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	f := excelize.NewFile()

// 	centeredStyle, err := f.NewStyle(&excelize.Style{
// 		Font: &excelize.Font{
// 			Family: "Segoe UI Variable Display Semib",
// 			Size:   12,
// 			Bold:   true,
// 		},
// 		Alignment: &excelize.Alignment{
// 			Horizontal: "center",
// 			Vertical:   "center",
// 			WrapText:   true,
// 		},
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cell style: " + err.Error()})
// 		return
// 	}

// 	timeSlots := []string{
// 		"08:45:00 - 09:35:00",
// 		"09:35:00 - 10:25:00",
// 		"10:40:00 - 11:30:00",
// 		"13:45:00 - 14:35:00",
// 		"14:35:00 - 15:25:00",
// 		"15:40:00 - 16:30:00",
// 	}
// 	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// 	departmentData := make(map[string]map[string]map[string]map[string]string)
// 	var academicYear string
// 	for rows.Next() {
// 		var dayName, startTime, endTime, classroom, subjectName, facultyName, deptName, semesterID string
// 		if err := rows.Scan(&dayName, &startTime, &endTime, &classroom, &subjectName, &facultyName, &deptName, &academicYear, &semesterID); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row scan error: " + err.Error()})
// 			return
// 		}

// 		if _, exists := departmentData[deptName]; !exists {
// 			departmentData[deptName] = make(map[string]map[string]map[string]string)
// 		}
// 		if _, exists := departmentData[deptName][semesterID]; !exists {
// 			departmentData[deptName][semesterID] = make(map[string]map[string]string)
// 		}
// 		if _, exists := departmentData[deptName][semesterID][dayName]; !exists {
// 			departmentData[deptName][semesterID][dayName] = make(map[string]string)
// 		}

// 		timeSlot := fmt.Sprintf("%s - %s", startTime, endTime)
// 		departmentData[deptName][semesterID][dayName][timeSlot] = fmt.Sprintf("%s\n%s", subjectName, facultyName)
// 	}

// 	for dept, semesterData := range departmentData {
// 		sheetName := dept
// 		f.NewSheet(sheetName)

// 		f.SetCellValue(sheetName, "A1", "BANNARI AMMAN INSTITUTE OF TECHNOLOGY")
// 		log.Println(academicYear)
// 		f.SetCellValue(sheetName, "A2", fmt.Sprintf("MASTER TIME TABLE - %s", academicYear))
// 		f.SetCellValue(sheetName, "A3", fmt.Sprintf("DEPARTMENT OF %s", dept))

// 		f.MergeCell(sheetName, "A1", "G1")
// 		f.MergeCell(sheetName, "A2", "G2")
// 		f.MergeCell(sheetName, "A3", "G3")
// 		f.SetCellStyle(sheetName, "A1", "G1", centeredStyle)
// 		f.SetCellStyle(sheetName, "A2", "G2", centeredStyle)
// 		f.SetCellStyle(sheetName, "A3", "G3", centeredStyle)

// 		f.SetColWidth(sheetName, "A", "A", 30)
// 		f.SetColWidth(sheetName, "B", string('B'+len(timeSlots)-1), 30)

// 		rowOffset := 5

// 		for semesterID, dayData := range semesterData {
// 			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), fmt.Sprintf("Semester %s", semesterID))
// 			rowOffset++

// 			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), "Day/Time")
// 			for i, timeSlot := range timeSlots {
// 				col := string('B' + i)
// 				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, rowOffset), timeSlot)
// 				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, rowOffset), fmt.Sprintf("%s%d", col, rowOffset), centeredStyle)
// 			}
// 			rowOffset++

// 			for i, day := range days {
// 				row := rowOffset + i
// 				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), day)
// 				f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), centeredStyle)

// 				for j, timeSlot := range timeSlots {
// 					col := string('B' + j)
// 					cellValue := dayData[day][timeSlot]
// 					f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, row), cellValue)
// 					f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), centeredStyle)
// 				}
// 			}

// 			rowOffset += len(days) + 2
// 		}
// 	}

// 	f.SetActiveSheet(0)
// 	f.DeleteSheet("Sheet1")

// 	filename := fmt.Sprintf("timetable_%s.xlsx", time.Now().Format("20060102150405"))
// 	c.Header("Content-Type", "application/octet-stream")
// 	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
// 	c.Header("Content-Transfer-Encoding", "binary")

//		if err := f.Write(c.Writer); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file: " + err.Error()})
//			return
//		}
//	}
package excel

import (
	"fmt"
	"net/http"
	"ssg-portal/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func Masterdownload(c *gin.Context) {
	academicYearID := c.Param("academic_year_id")

	if academicYearID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic year parameter"})
		return
	}

	// Updated query to fetch all semesters belonging to the given academic_year_id
	query := `
		SELECT t.day_name, t.start_time, t.end_time, t.classroom,
		t.subject_name, t.faculty_name, d.name AS department_name, may.academic_year, t.semester_id, t.section_id
		FROM timetable t
		JOIN departments d ON t.department_id = d.id
		JOIN master_academic_year may ON t.academic_year = may.id
		WHERE may.id = ?
	`

	rows, err := config.Database.Query(query, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error: " + err.Error()})
		return
	}
	defer rows.Close()

	f := excelize.NewFile()

	// Define a style for the header
	centeredStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Segoe UI Variable Display Semib",
			Size:   12,
			Bold:   true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cell style: " + err.Error()})
		return
	}

	timeSlots := []string{
		"08:45:00 - 09:35:00",
		"09:35:00 - 10:25:00",
		"10:40:00 - 11:30:00",
		"13:45:00 - 14:35:00",
		"14:35:00 - 15:25:00",
		"15:40:00 - 16:30:00",
	}
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	departmentData := make(map[string]map[string]map[string]map[string]map[string]string)
	var academicYear string
	for rows.Next() {
		var dayName, startTime, endTime, classroom, subjectName, facultyName, deptName, semesterID, section string
		if err := rows.Scan(&dayName, &startTime, &endTime, &classroom, &subjectName, &facultyName, &deptName, &academicYear, &semesterID, &section); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row scan error: " + err.Error()})
			return
		}

		// Organize the data into the departmentData map
		if _, exists := departmentData[deptName]; !exists {
			departmentData[deptName] = make(map[string]map[string]map[string]map[string]string)
		}
		if _, exists := departmentData[deptName][semesterID]; !exists {
			departmentData[deptName][semesterID] = make(map[string]map[string]map[string]string)
		}
		if _, exists := departmentData[deptName][semesterID][section]; !exists {
			departmentData[deptName][semesterID][section] = make(map[string]map[string]string)
		}
		if _, exists := departmentData[deptName][semesterID][section][dayName]; !exists {
			departmentData[deptName][semesterID][section][dayName] = make(map[string]string)
		}
		timeSlot := fmt.Sprintf("%s - %s", startTime, endTime)
		departmentData[deptName][semesterID][section][dayName][timeSlot] = fmt.Sprintf("%s\n%s", subjectName, facultyName)
	}

	// Populate Excel sheet with timetable data
	for dept, semesterData := range departmentData {
		sheetName := dept
		f.NewSheet(sheetName)

		f.SetCellValue(sheetName, "A1", "BANNARI AMMAN INSTITUTE OF TECHNOLOGY")
		f.SetCellValue(sheetName, "A2", fmt.Sprintf("MASTER TIME TABLE - %s", academicYear))
		f.SetCellValue(sheetName, "A3", fmt.Sprintf("DEPARTMENT OF %s", dept))

		// Merge header cells
		f.MergeCell(sheetName, "A1", "G1")
		f.MergeCell(sheetName, "A2", "G2")
		f.MergeCell(sheetName, "A3", "G3")
		f.SetCellStyle(sheetName, "A1", "G1", centeredStyle)
		f.SetCellStyle(sheetName, "A2", "G2", centeredStyle)
		f.SetCellStyle(sheetName, "A3", "G3", centeredStyle)

		// Set column width
		f.SetColWidth(sheetName, "A", "A", 30)
		f.SetColWidth(sheetName, "B", "Z", 30)

		rowOffset := 5

		for semesterID, sectionData := range semesterData {
			// Write Semester Header
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), fmt.Sprintf("Semester %s", semesterID))
			rowOffset++

			// Prepare columns for each section
			columns := make([]string, 0)
			for section := range sectionData {
				columns = append(columns, section)
			}

			// Write section headers
			for sectionIndex, section := range columns {
				startCol := int('B') + sectionIndex*len(timeSlots) // Starting column for this section
				endCol := startCol + len(timeSlots) - 1            // Ending column for this section

				// Convert startCol and endCol to column letters
				startColLetter := string(rune(startCol))
				endColLetter := string(rune(endCol))

				// Set cell value and merge cells for the section header
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", startColLetter, rowOffset), fmt.Sprintf("Section %s", section))
				f.MergeCell(sheetName, fmt.Sprintf("%s%d", startColLetter, rowOffset), fmt.Sprintf("%s%d", endColLetter, rowOffset))
				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", startColLetter, rowOffset), fmt.Sprintf("%s%d", endColLetter, rowOffset), centeredStyle)
			}
			rowOffset++

			// Write Time Slot Header for Each Section
			for sectionIndex := range columns {
				startCol := int('B') + sectionIndex*len(timeSlots)
				for j, timeSlot := range timeSlots {
					col := string(rune(startCol + j))
					f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, rowOffset), timeSlot)
					f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, rowOffset), fmt.Sprintf("%c%d", col, rowOffset), centeredStyle)
				}
			}
			rowOffset++

			// Fill in the timetable data
			for _, day := range days {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), day)
				f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowOffset), fmt.Sprintf("A%d", rowOffset), centeredStyle)

				// Fill the data for each section
				for sectionIndex, section := range columns {
					startCol := int('B') + sectionIndex*len(timeSlots)
					for j, timeSlot := range timeSlots {
						col := string(rune(startCol + j))
						cellValue := sectionData[section][day][timeSlot] // Fetching the subject and faculty

						f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, rowOffset), cellValue)
						f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, rowOffset), fmt.Sprintf("%c%d", col, rowOffset), centeredStyle)
					}
				}
				rowOffset++
			}

			rowOffset += 2 // Add space between semesters
		}
	}

	f.SetActiveSheet(0)
	f.DeleteSheet("Sheet1")

	// Set the filename for the Excel file
	filename := fmt.Sprintf("timetable_%s_%s.xlsx", academicYear, time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Transfer-Encoding", "binary")

	// Write the file to the response
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file: " + err.Error()})
	}
}
