package excel

import (
	"fmt"
	"log"
	"net/http"
	"ssg-portal/config"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"ssg-portal/config"

// 	"github.com/gin-gonic/gin"
// 	"github.com/xuri/excelize/v2"
// )

// func DownloadTimetable(c *gin.Context) {

// 	semesterIDParam := c.Param("semesterId")
// 	semesterID, err := strconv.Atoi(semesterIDParam)
// 	if err != nil || semesterID < 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester_id"})
// 		return
// 	}

// 	var academicYear string
// 	yearQuery := `SELECT
//     DISTINCT
//     m.academic_year
// FROM
//     timetable t
// JOIN
//     master_academic_year m
// ON
//     t.academic_year = m.id
// WHERE
//     t.semester_id = ?`
// 	err = config.Database.QueryRow(yearQuery, semesterID).Scan(&academicYear)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	query := `SELECT t.day_name, t.start_time, t.end_time, t.classroom, t.semester_id, t.department_id, t.section_id,
// 	t.subject_name, t.faculty_name, d.name as department_name
// FROM timetable t
// JOIN departments d ON t.department_id = d.id
// WHERE t.semester_id = ?
// ORDER BY t.section_id`

// 	rows, err := config.Database.Query(query, semesterID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	file := excelize.NewFile()
// 	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
// 	timeSlots := []string{
// 		"08:45:00 - 09:35:00",
// 		"09:35:00 - 10:25:00",
// 		"10:40:00 - 11:30:00",
// 		"13:45:00 - 14:35:00",
// 		"14:35:00 - 15:25:00",
// 		"15:40:00 - 16:30:00",
// 	}
// 	dayIndex := map[string]int{
// 		"Monday":    0,
// 		"Tuesday":   1,
// 		"Wednesday": 2,
// 		"Thursday":  3,
// 		"Friday":    4,
// 		"Saturday":  5,
// 	}

// 	timetableData := make(map[string][][]string)
// 	for rows.Next() {
// 		var day, startTime, endTime, classroom, section, subject, faculty string
// 		var semesterID, departmentID int
// 		var departmentName string
// 		if err := rows.Scan(&day, &startTime, &endTime, &classroom, &semesterID, &departmentID, &section, &subject, &faculty, &departmentName); err != nil {
// 			log.Fatal(err)
// 		}

// 		// Create a key for timetable data, truncating the section if necessary.
// 		key := fmt.Sprintf("%s-S%d", departmentName, semesterID)
// 		if section != "" {
// 			key += fmt.Sprintf("-Sec%s", section)
// 		}

// 		if _, exists := timetableData[key]; !exists {
// 			timetableData[key] = make([][]string, len(days))
// 			for i := range timetableData[key] {
// 				timetableData[key][i] = make([]string, len(timeSlots))
// 			}
// 		}

// 		rowIdx := dayIndex[day]

// 		for i, slot := range timeSlots {
// 			if startTime == slot[:8] && endTime == slot[11:] {
// 				// Append the subject and faculty to the existing data for the time slot
// 				if timetableData[key][rowIdx][i] != "" {
// 					timetableData[key][rowIdx][i] += "\n"
// 				}
// 				timetableData[key][rowIdx][i] += fmt.Sprintf("%s\n%s", subject, faculty)
// 				break
// 			}
// 		}
// 	}

// 	centeredStyle, err := file.NewStyle(&excelize.Style{
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
// 		log.Fatal(err)
// 	}

// 	for key, data := range timetableData {
// 		// Truncate the sheet name to ensure it doesn't exceed 31 characters.
// 		sheetName := key
// 		if len(sheetName) > 31 {
// 			sheetName = sheetName[:31]
// 		}

// 		index, err := file.NewSheet(sheetName)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		file.SetColWidth(sheetName, "A", "A", 30)
// 		file.SetColWidth(sheetName, "B", string('B'+len(timeSlots)-1), 30)

// 		file.SetCellValue(sheetName, "A1", "Day/Time")
// 		file.SetCellStyle(sheetName, "A1", "A1", centeredStyle)

// 		for i, timing := range timeSlots {
// 			cell := fmt.Sprintf("%s1", string('B'+i))
// 			file.SetCellValue(sheetName, cell, timing)
// 			file.SetCellStyle(sheetName, cell, cell, centeredStyle)
// 		}

// 		for i, day := range days {
// 			file.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), day)
// 			file.SetCellStyle(sheetName, fmt.Sprintf("A%d", i+2), fmt.Sprintf("A%d", i+2), centeredStyle)
// 		}

// 		for i, dayData := range data {
// 			for j, cellData := range dayData {
// 				if cellData != "" {
// 					cell := fmt.Sprintf("%s%d", string('B'+j), i+2)
// 					file.SetCellValue(sheetName, cell, cellData)
// 					file.SetCellStyle(sheetName, cell, cell, centeredStyle)
// 				}
// 			}
// 		}

// 		for i := 1; i <= len(days)+1; i++ {
// 			file.SetRowHeight(sheetName, i, 65)
// 		}

// 		file.SetActiveSheet(index)
// 	}

// 	file.DeleteSheet("Sheet1")

// 	filename := fmt.Sprintf("%s-Semester-%d.xlsx", academicYear, semesterID)

// 	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// 	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
// 	c.Header("Content-Transfer-Encoding", "binary")
// 	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")

//		if err := file.Write(c.Writer); err != nil {
//			log.Fatal(err)
//		}
//	}
func DownloadTimetable(c *gin.Context) {

	semesterIDParam := c.Param("semesterId")
	semesterID, err := strconv.Atoi(semesterIDParam)
	if err != nil || semesterID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid semester_id"})
		return
	}

	var academicYear string
	yearQuery := `SELECT 
    DISTINCT m.academic_year
	FROM 
	    timetable t
	JOIN 
	    master_academic_year m
	ON 
	    t.academic_year = m.id
	WHERE 
	    t.semester_id = ?`
	err = config.Database.QueryRow(yearQuery, semesterID).Scan(&academicYear)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch unique working dates for the semester.
	datesQuery := `SELECT DISTINCT t.day_name 
	FROM timetable t 
	WHERE t.semester_id = ? 
	ORDER BY t.day_name`
	dateRows, err := config.Database.Query(datesQuery, semesterID)
	if err != nil {
		log.Fatal(err)
	}
	defer dateRows.Close()

	var workingDates []string
	for dateRows.Next() {
		var date string
		if err := dateRows.Scan(&date); err != nil {
			log.Fatal(err)
		}
		workingDates = append(workingDates, date)
	}

	query := `SELECT t.day_name, t.start_time, t.end_time, t.classroom, t.semester_id, t.department_id, t.section_id, 
	t.subject_name, t.faculty_name, d.name as department_name
	FROM timetable t
	JOIN departments d ON t.department_id = d.id
	WHERE t.semester_id = ?
	ORDER BY t.section_id`

	rows, err := config.Database.Query(query, semesterID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	file := excelize.NewFile()
	timeSlots := []string{
		"08:45:00 - 09:35:00",
		"09:35:00 - 10:25:00",
		"10:40:00 - 11:30:00",
		"11:30:00 - 12:20:00",
		"13:30:00 - 14:20:00",
		"14:20:00 - 15:10:00",
		"15:25:00 - 16:30:00",
	}

	timetableData := make(map[string]map[string][]string)
	for rows.Next() {
		var date, startTime, endTime, classroom, section, subject, faculty string
		var semesterID, departmentID int
		var departmentName string
		if err := rows.Scan(&date, &startTime, &endTime, &classroom, &semesterID, &departmentID, &section, &subject, &faculty, &departmentName); err != nil {
			log.Fatal(err)
		}

		key := fmt.Sprintf("%s-S%d", departmentName, semesterID)
		if section != "" {
			key += fmt.Sprintf("-Sec%s", section)
		}

		if _, exists := timetableData[key]; !exists {
			timetableData[key] = make(map[string][]string)
		}
		if _, exists := timetableData[key][date]; !exists {
			timetableData[key][date] = make([]string, len(timeSlots))
		}

		for i, slot := range timeSlots {
			if startTime == slot[:8] && endTime == slot[11:] {
				if timetableData[key][date][i] != "" {
					timetableData[key][date][i] += "\n"
				}
				timetableData[key][date][i] += fmt.Sprintf("%s\n%s", subject, faculty)
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
		if len(sheetName) > 31 {
			sheetName = sheetName[:31]
		}

		index, err := file.NewSheet(sheetName)
		if err != nil {
			log.Fatal(err)
		}

		file.SetColWidth(sheetName, "A", "A", 30)
		file.SetColWidth(sheetName, "B", string('B'+len(timeSlots)-1), 30)

		file.SetCellValue(sheetName, "A1", "Date/Time")
		file.SetCellStyle(sheetName, "A1", "A1", centeredStyle)

		for i, timing := range timeSlots {
			cell := fmt.Sprintf("%s1", string('B'+i))
			file.SetCellValue(sheetName, cell, timing)
			file.SetCellStyle(sheetName, cell, cell, centeredStyle)
		}

		for i, date := range workingDates {
			file.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), date)
			file.SetCellStyle(sheetName, fmt.Sprintf("A%d", i+2), fmt.Sprintf("A%d", i+2), centeredStyle)

			if dayData, ok := data[date]; ok {
				for j, cellData := range dayData {
					if cellData != "" {
						cell := fmt.Sprintf("%s%d", string('B'+j), i+2)
						file.SetCellValue(sheetName, cell, cellData)
						file.SetCellStyle(sheetName, cell, cell, centeredStyle)
					}
				}
			}
		}

		for i := 1; i <= len(workingDates)+1; i++ {
			file.SetRowHeight(sheetName, i, 65)
		}

		file.SetActiveSheet(index)
	}

	file.DeleteSheet("Sheet1")

	filename := fmt.Sprintf("%s-Semester-%d.xlsx", academicYear, semesterID)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")

	if err := file.Write(c.Writer); err != nil {
		log.Fatal(err)
	}
}
