package excel

import (
	"fmt"
	"log"
	"net/http"
	"ssg-portal/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func Masterdownload(c *gin.Context) {
	academicYearID := c.Param("academic_year_id")
	semesterType := c.Param("type")

	if academicYearID == "" || (semesterType != "odd" && semesterType != "even") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	semesterIDs := []string{"1", "3", "5", "7"}
	if semesterType == "even" {
		semesterIDs = []string{"2", "4", "6", "8"}
	}
	semesterIDsStr := strings.Join(semesterIDs, ",")

	query := fmt.Sprintf(`
		SELECT t.day_name, t.start_time, t.end_time, t.classroom,
		       t.subject_name, t.faculty_name, d.name AS department_name, ay.academic_year, t.semester_id
		FROM timetable t
		JOIN departments d ON t.department_id = d.id
		JOIN academic_year ay ON t.academic_year = ay.id
		WHERE ay.id = ? AND t.semester_id IN (%s)
	`, semesterIDsStr)

	rows, err := config.Database.Query(query, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error: " + err.Error()})
		return
	}
	defer rows.Close()

	f := excelize.NewFile()

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

	departmentData := make(map[string]map[string]map[string]map[string]string)
	var academicYear string
	for rows.Next() {
		var dayName, startTime, endTime, classroom, subjectName, facultyName, deptName, semesterID string
		if err := rows.Scan(&dayName, &startTime, &endTime, &classroom, &subjectName, &facultyName, &deptName, &academicYear, &semesterID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row scan error: " + err.Error()})
			return
		}

		if _, exists := departmentData[deptName]; !exists {
			departmentData[deptName] = make(map[string]map[string]map[string]string)
		}
		if _, exists := departmentData[deptName][semesterID]; !exists {
			departmentData[deptName][semesterID] = make(map[string]map[string]string)
		}
		if _, exists := departmentData[deptName][semesterID][dayName]; !exists {
			departmentData[deptName][semesterID][dayName] = make(map[string]string)
		}

		timeSlot := fmt.Sprintf("%s - %s", startTime, endTime)
		departmentData[deptName][semesterID][dayName][timeSlot] = fmt.Sprintf("%s\n%s", subjectName, facultyName)
	}

	for dept, semesterData := range departmentData {
		sheetName := dept
		f.NewSheet(sheetName)

		f.SetCellValue(sheetName, "A1", "BANNARI AMMAN INSTITUTE OF TECHNOLOGY")
		log.Println(academicYear)
		f.SetCellValue(sheetName, "A2", fmt.Sprintf("MASTER TIME TABLE - %s (%s SEMESTER)", academicYear, (semesterType)))
		f.SetCellValue(sheetName, "A3", fmt.Sprintf("DEPARTMENT OF %s", dept))

		f.MergeCell(sheetName, "A1", "G1")
		f.MergeCell(sheetName, "A2", "G2")
		f.MergeCell(sheetName, "A3", "G3")
		f.SetCellStyle(sheetName, "A1", "G1", centeredStyle)
		f.SetCellStyle(sheetName, "A2", "G2", centeredStyle)
		f.SetCellStyle(sheetName, "A3", "G3", centeredStyle)

		f.SetColWidth(sheetName, "A", "A", 30)
		f.SetColWidth(sheetName, "B", string('B'+len(timeSlots)-1), 30)

		rowOffset := 5

		for semesterID, dayData := range semesterData {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), fmt.Sprintf("Semester %s", semesterID))
			rowOffset++

			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowOffset), "Day/Time")
			for i, timeSlot := range timeSlots {
				col := string('B' + i)
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, rowOffset), timeSlot)
				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, rowOffset), fmt.Sprintf("%s%d", col, rowOffset), centeredStyle)
			}
			rowOffset++

			for i, day := range days {
				row := rowOffset + i
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), day)
				f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), centeredStyle)

				for j, timeSlot := range timeSlots {
					col := string('B' + j)
					cellValue := dayData[day][timeSlot]
					f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, row), cellValue)
					f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), centeredStyle)
				}
			}

			rowOffset += len(days) + 2
		}
	}

	f.SetActiveSheet(0)
	f.DeleteSheet("Sheet1")

	filename := fmt.Sprintf("timetable_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file: " + err.Error()})
		return
	}
}
