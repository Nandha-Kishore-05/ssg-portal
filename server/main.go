package main

import (
	"ssg-portal/config"
	"ssg-portal/timetables"
	"ssg-portal/timetables/allocation"
	"ssg-portal/timetables/auth"
	"ssg-portal/timetables/excel"
	"ssg-portal/timetables/manualentry"
	"ssg-portal/timetables/otp"
	"ssg-portal/timetables/role"
	"ssg-portal/timetables/routes"
    "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	r.GET("/timetable/:departmentID/:semesterId/:academicYearID", routes.GenerateTimetable)
    r.POST("/timetable/save", timetables.SaveTimetable)
	r.GET("/timetable/saved/:departmentID/:semesterId/:academicYearID", timetables.GetTimetable)
	r.GET("/timetable/faculty/:faculty_name/:academicYearID", timetables.FacultyTimetable)
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
	r.GET("/acdemicYearOptions", timetables.AcademicYearOptions)
	r.GET("/periodallocation", allocation.Subjectallocation)
	r.PUT("/periodallocationedit", allocation.UpdateAllocation)
    r.GET("/getResource",role.GetResources)
	r.POST("/send-otp",otp.SendOTPHandler)
	r.POST("/verify-otp",otp.VerifyOTPHandler)
	r.GET("/downloadTimetable/:semesterId", excel.DownloadTimetable)
	r.Run(":8080")
}
