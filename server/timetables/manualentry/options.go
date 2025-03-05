package manualentry

import (
    "net/http"
    "ssg-portal/config"

    "github.com/gin-gonic/gin"
)

func DayAndTimeOptions(c *gin.Context) {
 
    dayRows, err := config.Database.Query("SELECT working_date FROM master_workingdays")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer dayRows.Close()


    var dayOptions []map[string]string
    for dayRows.Next() {
        var dayName string
        if err := dayRows.Scan(&dayName); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        dayOptions = append(dayOptions, map[string]string{
            "label": dayName,
            "value": dayName,
        })
    }


    if err := dayRows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


    startTimeRows, err := config.Database.Query("SELECT start_time FROM hours")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer startTimeRows.Close()


    var startTimeOptions []map[string]string
    for startTimeRows.Next() {
        var startTime string
        if err := startTimeRows.Scan(&startTime); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        startTimeOptions = append(startTimeOptions, map[string]string{
            "label": startTime,
            "value": startTime,
        })
    }


    if err := startTimeRows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

 
    endTimeRows, err := config.Database.Query("SELECT end_time FROM hours")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer endTimeRows.Close()


    var endTimeOptions []map[string]string
    for endTimeRows.Next() {
        var endTime string
        if err := endTimeRows.Scan(&endTime); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        endTimeOptions = append(endTimeOptions, map[string]string{
            "label": endTime,
            "value": endTime,
        })
    }

    if err := endTimeRows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


    facultyRows, err := config.Database.Query("SELECT name FROM faculty")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer facultyRows.Close()

    var facultyOptions []map[string]string
    for facultyRows.Next() {
        var facultyName string
        if err := facultyRows.Scan(&facultyName); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        facultyOptions = append(facultyOptions, map[string]string{
            "label": facultyName,
            "value": facultyName,
        })
    }


    if err := facultyRows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "dayOptions":       dayOptions,
        "startTimeOptions": startTimeOptions,
        "endTimeOptions":   endTimeOptions,

        "facultyOptions": facultyOptions,
    })
}
