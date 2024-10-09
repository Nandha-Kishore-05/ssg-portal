package excel

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"ssg-portal/config"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func DownloadTimetable(c *gin.Context) {

	semesterIDParam := c.Param("semesterId")
	semesterID, err := strconv.Atoi(semesterIDParam)
	if err != nil || semesterID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester_id"})
		return
	}

	var academicYear, semesterName string
	yearQuery := `SELECT may.academic_year, s.semester_name
FROM academic_year ay
JOIN semester s ON ay.semester_id = s.id
JOIN master_academic_year may ON ay.academic_year = may.id
WHERE ay.semester_id = ?
`
	err = config.Database.QueryRow(yearQuery, semesterID).Scan(&academicYear, &semesterName)
	if err != nil {
		log.Fatal(err)
	}

	var semesterType string
	if semesterID%2 == 0 {
		semesterType = "Even"
	} else {
		semesterType = "Odd"
	}

	query := `
        SELECT day_name, start_time, end_time, classroom, semester_id, department_id, subject_name, faculty_name
        FROM timetable
        WHERE semester_id = ?
        ORDER BY department_id, semester_id
    `
	rows, err := config.Database.Query(query, semesterID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	file := excelize.NewFile()

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	timeSlots := []string{
		"08:45:00 - 09:35:00",
		"09:35:00 - 10:25:00",
		"10:40:00 - 11:30:00",
		"13:45:00 - 14:35:00",
		"14:35:00 - 15:25:00",
		"15:40:00 - 16:30:00",
	}
	dayIndex := map[string]int{
		"Monday":    0,
		"Tuesday":   1,
		"Wednesday": 2,
		"Thursday":  3,
		"Friday":    4,
		"Saturday":  5,
	}

	timetableData := make(map[string][][]string)

	for rows.Next() {
		var day, startTime, endTime, classroom, subject, faculty string
		var semesterID, departmentID int
		if err := rows.Scan(&day, &startTime, &endTime, &classroom, &semesterID, &departmentID, &subject, &faculty); err != nil {
			log.Fatal(err)
		}

		key := fmt.Sprintf("Department %d - Semester %d", departmentID, semesterID)

		if _, exists := timetableData[key]; !exists {
			timetableData[key] = make([][]string, len(days))
			for i := range timetableData[key] {
				timetableData[key][i] = make([]string, len(timeSlots))
			}
		}

		rowIdx := dayIndex[day]

		for i, slot := range timeSlots {
			if startTime == slot[:8] && endTime == slot[11:] {
				timetableData[key][rowIdx][i] = fmt.Sprintf("%s\n%s", subject, faculty)
				break
			}
		}
	}

	centeredStyle, err := file.NewStyle(&excelize.Style{
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
		log.Fatal(err)
	}

	for key, data := range timetableData {
		sheetName := key
		index, err := file.NewSheet(sheetName)
		if err != nil {
			log.Fatal(err)
		}

		file.SetColWidth(sheetName, "A", "A", 30)
		file.SetColWidth(sheetName, "B", string('B'+len(timeSlots)-1), 30)

		file.SetCellValue(sheetName, "A1", "Day/Time")
		file.SetCellStyle(sheetName, "A1", "A1", centeredStyle)

		for i, timing := range timeSlots {
			cell := fmt.Sprintf("%s1", string('B'+i))
			file.SetCellValue(sheetName, cell, timing)
			file.SetCellStyle(sheetName, cell, cell, centeredStyle)
		}

		for i, day := range days {
			file.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), day)
			file.SetCellStyle(sheetName, fmt.Sprintf("A%d", i+2), fmt.Sprintf("A%d", i+2), centeredStyle)
		}

		for i, dayData := range data {
			for j, cellData := range dayData {
				if cellData != "" {
					cell := fmt.Sprintf("%s%d", string('B'+j), i+2)
					file.SetCellValue(sheetName, cell, cellData)
					file.SetCellStyle(sheetName, cell, cell, centeredStyle)
				}
			}
		}

		for i := 1; i <= len(days)+1; i++ {
			file.SetRowHeight(sheetName, i, 65)
		}

		file.SetActiveSheet(index)
	}

	file.DeleteSheet("Sheet1")

	filename := fmt.Sprintf("%s-Semester-%d-(%s sem).xlsx", academicYear, semesterID, semesterType)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")

	if err := file.Write(c.Writer); err != nil {
		log.Fatal(err)
	}
}
