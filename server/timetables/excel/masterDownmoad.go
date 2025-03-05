package excel

import (
	"fmt"
	"net/http"
	"ssg-portal/config"
	"strings"
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

	// Fetch distinct days dynamically
	dayQuery := `SELECT DISTINCT day_name FROM timetable WHERE academic_year = ? ORDER BY STR_TO_DATE(day_name, '%Y-%m-%d')`
	dayRows, err := config.Database.Query(dayQuery, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timetable dates: " + err.Error()})
		return
	}
	defer dayRows.Close()

	var days []string
	for dayRows.Next() {
		var day string
		if err := dayRows.Scan(&day); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning timetable dates: " + err.Error()})
			return
		}
		days = append(days, day)
	}

	// Ensure that the `days` array is not empty
	if len(days) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No timetable data found for the given academic year"})
		return
	}

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

	// Define a style for headers
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
		"11:30:00 - 12:20:00",
		"13:30:00 - 14:20:00",
		"14:20:00 - 15:10:00",
		"15:25:00 - 16:30:00",
	}

	departmentData := make(map[string]map[string]map[string]map[string]map[string][]string)

	var academicYear string
	for rows.Next() {
		var dayName, startTime, endTime, classroom, subjectName, facultyName, deptName, semesterID, section string
		if err := rows.Scan(&dayName, &startTime, &endTime, &classroom, &subjectName, &facultyName, &deptName, &academicYear, &semesterID, &section); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row scan error: " + err.Error()})
			return
		}
		timeSlot := fmt.Sprintf("%s - %s", startTime, endTime)

		if _, exists := departmentData[deptName]; !exists {
			departmentData[deptName] = make(map[string]map[string]map[string]map[string][]string)
		}
		if _, exists := departmentData[deptName][semesterID]; !exists {
			departmentData[deptName][semesterID] = make(map[string]map[string]map[string][]string)
		}
		if _, exists := departmentData[deptName][semesterID][section]; !exists {
			departmentData[deptName][semesterID][section] = make(map[string]map[string][]string)
		}
		if _, exists := departmentData[deptName][semesterID][section][dayName]; !exists {
			departmentData[deptName][semesterID][section][dayName] = make(map[string][]string)
		}
		departmentData[deptName][semesterID][section][dayName][timeSlot] = append(
			departmentData[deptName][semesterID][section][dayName][timeSlot],
			fmt.Sprintf("%s\n%s", subjectName, facultyName),
		)
	}

	for dept, semesterData := range departmentData {
		sheetName := dept
		f.NewSheet(sheetName)

		f.SetCellValue(sheetName, "A1", "BANNARI AMMAN INSTITUTE OF TECHNOLOGY")
		f.SetCellValue(sheetName, "A2", fmt.Sprintf("MASTER TIME TABLE - %s", academicYear))
		f.SetCellValue(sheetName, "A3", fmt.Sprintf("DEPARTMENT OF %s", dept))

		f.MergeCell(sheetName, "A1", "G1")
		f.MergeCell(sheetName, "A2", "G2")
		f.MergeCell(sheetName, "A3", "G3")
		f.SetCellStyle(sheetName, "A1", "G1", centeredStyle)
		f.SetCellStyle(sheetName, "A2", "G2", centeredStyle)
		f.SetCellStyle(sheetName, "A3", "G3", centeredStyle)

		f.SetColWidth(sheetName, "A", "A", 30)
		f.SetColWidth(sheetName, "B", "Z", 30)

		rowOffset := 5

		for semesterID, sectionData := range semesterData {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), fmt.Sprintf("Semester %s", semesterID))
			rowOffset++

			columns := make([]string, 0)
			for section := range sectionData {
				columns = append(columns, section)
			}

			for sectionIndex, section := range columns {
				startCol := 2 + sectionIndex*len(timeSlots)
				startColLetter, _ := excelize.ColumnNumberToName(startCol)
				endColLetter, _ := excelize.ColumnNumberToName(startCol + len(timeSlots) - 1)

				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", startColLetter, rowOffset), fmt.Sprintf("Section %s", section))
				f.MergeCell(sheetName, fmt.Sprintf("%s%d", startColLetter, rowOffset), fmt.Sprintf("%s%d", endColLetter, rowOffset))
				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", startColLetter, rowOffset), fmt.Sprintf("%s%d", endColLetter, rowOffset), centeredStyle)
			}
			rowOffset++

			for _, day := range days {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), day)
				f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowOffset), fmt.Sprintf("A%d", rowOffset), centeredStyle)

				for sectionIndex, section := range columns {
					startCol := 2 + sectionIndex*len(timeSlots)
					for j, timeSlot := range timeSlots {
						col, _ := excelize.ColumnNumberToName(startCol + j)
						cellValues := sectionData[section][day][timeSlot]
						combinedValue := strings.Join(cellValues, "\n-----\n")
						f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, rowOffset), combinedValue)
						f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, rowOffset), fmt.Sprintf("%s%d", col, rowOffset), centeredStyle)
					}
				}

				rowOffset++
			}

			rowOffset += 2
		}
	}

	f.SetActiveSheet(0)
	f.DeleteSheet("Sheet1")

	filename := fmt.Sprintf("timetable_%s_%s.xlsx", academicYear, time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file: " + err.Error()})
	}
}
