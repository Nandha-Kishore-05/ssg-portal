
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ssg-portal/config"
	"ssg-portal/timetables"
	"ssg-portal/timetables/routes"
)

func main() {
	config.ConnectDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/timetable/:departmentID/:semesterId", routes.GenerateTimetable)

	r.POST("/timetable/save", timetables.SaveTimetable)
	r.GET("/timetable/saved/:departmentID/:semesterId", timetables.GetTimetable)
	r.GET("/timetable/faculty/:faculty_name", timetables.FacultyTimetable)
	r.GET("/timetable/lab/:subject_name", timetables.LabTableTimetable)
	r.GET("/timetable/facultyOptions", timetables.FacultyOptions)
	r.GET("/timetable/labOptions", timetables.LabOptions)
	r.GET("/timetable/options", timetables.TimetableOptions)
	r.GET("/timetable/semoptions", timetables.SemOptions)
	r.POST("/upload", timetables.Uploaddetails)
	r.GET("faculty/available/:departmentID/:semesterID/:day/:startTime/:endTime", timetables.GetAvailableFaculty)
	r.GET("/saved/deptoptions", timetables.SavedDepartmentOptions)
	r.Run(":8080")
}
