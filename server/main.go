package main

import (
	"ssg-portal/config"
	"ssg-portal/timetables"
	"ssg-portal/timetables/allocation"
	"ssg-portal/timetables/auth"
	"ssg-portal/timetables/excel"
	"ssg-portal/timetables/labentry"
	"ssg-portal/timetables/manualentry"
	"ssg-portal/timetables/otp"
	"ssg-portal/timetables/role"
	"ssg-portal/timetables/routes"
	"ssg-portal/timetables/studententry"

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
	r.GET("/timetable/:departmentID/:semesterId/:academicYearID/:sectionID", routes.GenerateTimetable)
	r.POST("/timetable/save", timetables.SaveTimetable)
	r.GET("/timetable/saved/:departmentID/:semesterId/:academicYearID/:sectionID", timetables.GetTimetable)
	r.GET("/timetable/faculty/:faculty_name/:academicYearID", timetables.FacultyTimetable)
	r.GET("/timetable/lab/:subject_name", timetables.LabTableTimetable)
	r.GET("/timetable/facultyOptions", timetables.FacultyOptions)
	r.GET("/timetable/labOptions", timetables.LabOptions)
	r.GET("/timetable/options", timetables.TimetableOptions)
	r.GET("/classroomOptions", timetables.VenueOptions)
	r.GET("/timetable/semoptions", timetables.SemOptions)
	r.GET("/timetable/sectionoptions", timetables.SectionOptions)
	r.POST("/upload", timetables.UploadDetails)
	r.GET("faculty/available/:departmentID/:semesterID/:day/:startTime/:endTime/:academicYearID/:sectionID", timetables.GetAvailableFaculty)
	r.GET("/available-timings/:facultyName/:day", timetables.GetAvailableTimingsForFaculty)
	r.GET("/saved/deptoptions", timetables.SavedDepartmentOptions)
	r.PUT("/timetable/update", timetables.UpdateTimetable)
	r.GET("/manual/options", manualentry.DayAndTimeOptions)
	r.POST("/manual/submit", manualentry.SubmitManualEntry)
	r.GET("/classroomDetailsOptions", timetables.ClassroomDetailsOptions)
	r.GET("/acdemicYearOptions", timetables.AcademicYearOptions)
	r.GET("/periodallocation", allocation.Subjectallocation)
	r.PUT("/periodallocationedit", allocation.UpdateAllocation)
	r.GET("/getResource", role.GetResources)
	r.POST("/send-otp", otp.SendOTPHandler)
	r.POST("/verify-otp", otp.VerifyOTPHandler)
	r.GET("/downloadTimetable/:semesterId", excel.DownloadTimetable)
	r.GET("/venueTimetable/:classroom", timetables.VenueTimetable)
	r.GET("/venueTimetableOptions", timetables.ClassroomOptions)
	r.POST("/studententry/upload", studententry.InsertStudentEntries)
	r.GET("/studentTimetable/:studentID", studententry.StudentTimetable)
	r.GET("/download/:academic_year_id", excel.Masterdownload)
	r.GET("/studentoptions", studententry.GetStudentOptions)
	r.GET("/subjectoptions", timetables.SubjectOptions)
	r.GET("/subjectTypeoptions", timetables.SubjectTypeOptions)
	r.GET("/classroomavailabletimings/:academicYearID/:facultyName/:day/:classroomName", timetables.GetAvailableTimingsForFacultyAndClassroom)
	r.GET("/workingDayoptions", timetables.WorkingDayOptions)
	r.GET("/course-code", timetables.CourseCodeOptions)
	r.POST("/labentry", labentry.LabEntry)
	r.POST("/manual/bulksubmit", manualentry.BulkInsert)
	r.POST("/upload-lab", timetables.HandleExcelUpload)
	r.Run(":8080")

}
