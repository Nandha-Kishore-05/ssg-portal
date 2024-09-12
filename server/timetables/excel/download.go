package excel

import (
    "fmt"
    "log"
    

    "ssg-portal/config"
    "github.com/gin-gonic/gin"
    "github.com/xuri/excelize/v2"
)

func DownloadTimetable(c *gin.Context) {


    query := `
        SELECT day_name, start_time, end_time, classroom, semester_id, department_id, subject_name, faculty_name
        FROM timetable ORDER BY department_id, semester_id
    `
    rows, err := config.Database.Query(query)
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

  
    for key, data := range timetableData {
        sheetName := key
        index, err := file.NewSheet(sheetName) 
        if err != nil {
            log.Fatal(err)
        }

      
        for i, timing := range timeSlots {
            cell := fmt.Sprintf("%s1", string('B'+i))
            file.SetCellValue(sheetName, cell, timing)
        }

    
        for i, day := range days {
            file.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), day)
        }

      
        for i, dayData := range data {
            for j, cellData := range dayData {
                if cellData != "" {
                    file.SetCellValue(sheetName, fmt.Sprintf("%s%d", string('B'+j), i+2), cellData)
                }
            }
        }

        file.SetActiveSheet(index)
    }

 
    file.DeleteSheet("Sheet1")

  
    c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    c.Header("Content-Disposition", "attachment; filename=department_timetable.xlsx")
    c.Header("Content-Transfer-Encoding", "binary")


    if err := file.Write(c.Writer); err != nil {
        log.Fatal(err)
    }
}
