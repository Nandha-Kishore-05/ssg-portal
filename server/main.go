package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
    "ssg-portal/config"
	"ssg-portal/timetables"
	"ssg-portal/timetables/allocation"
	"ssg-portal/timetables/auth"
	"ssg-portal/timetables/manualentry"
	"ssg-portal/timetables/otp"
	"ssg-portal/timetables/role"
	"ssg-portal/timetables/routes"
)

func main() {
	config.ConnectDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    r.POST("/login", auth.Login)
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
	r.PUT("/timetable/update", timetables.UpdateTimetable)
	r.GET("/manual/options", manualentry.DayAndTimeOptions)
    r.POST("/manual/submit", manualentry.SubmitManualEntry)
	r.GET("/menuitems", timetables.GetMenuItems)
	r.GET("/periodallocation", allocation.Subjectallocation)
	r.PUT("/periodallocationedit", allocation.UpdateAllocation)
	
    r.GET("/getResource",role.GetResources)
	r.POST("/send-otp",otp.SendOTPHandler)
	r.POST("/verify-otp",otp.VerifyOTPHandler)
	r.Run(":8080")
}
