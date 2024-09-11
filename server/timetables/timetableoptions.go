package timetables

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "ssg-portal/config"
)

func TimetableOptions(c *gin.Context) {
    rows, err := config.Database.Query("SELECT id, name FROM departments")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var options []map[string]interface{}
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        options = append(options, map[string]interface{}{
            "label": name,
            "value": id,
        })
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, options)
}

// package timetables

// import (
//     "net/http"
//     "github.com/gin-gonic/gin"
//     "ssg-portal/config"
// )

// func TimetableOptions(c *gin.Context) {
//     // Query for semester options
//     semRows, err := config.Database.Query("SELECT id, semester_name FROM semester")
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
//     defer semRows.Close()

//     var semOptions []map[string]interface{}
//     for semRows.Next() {
//         var id int
//         var name string
//         if err := semRows.Scan(&id, &name); err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//             return
//         }
//         semOptions = append(semOptions, map[string]interface{}{
//             "label": name,
//             "value": id,
//         })
//     }

//     if err := semRows.Err(); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     // Query for department options
//     deptRows, err := config.Database.Query("SELECT id, name FROM departments")
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
//     defer deptRows.Close()

//     var deptOptions []map[string]interface{}
//     for deptRows.Next() {
//         var id int
//         var name string
//         if err := deptRows.Scan(&id, &name); err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//             return
//         }
//         deptOptions = append(deptOptions, map[string]interface{}{
//             "label": name,
//             "value": id,
//         })
//     }

//     if err := deptRows.Err(); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     // Return both semester and department options
//     c.JSON(http.StatusOK, gin.H{
//         "semesterOptions": semOptions,
//         "departmentOptions": deptOptions,
//     })
// }
