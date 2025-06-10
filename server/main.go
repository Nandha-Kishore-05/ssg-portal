package main

import (
	"ssg-portal/config"
	"ssg-portal/controllers"
	"ssg-portal/controllers/academicYear"
	"ssg-portal/controllers/allocation"
	"ssg-portal/controllers/auth"
	"ssg-portal/controllers/classroom"
	"ssg-portal/controllers/date"
	"ssg-portal/controllers/edit"
	"ssg-portal/controllers/excel"
	"ssg-portal/controllers/facultyDetails"
	groups "ssg-portal/controllers/group"
	"ssg-portal/controllers/lab"
	"ssg-portal/controllers/labentry"
	"ssg-portal/controllers/manualentry"
	"ssg-portal/controllers/otp"
	"ssg-portal/controllers/role"
	"ssg-portal/controllers/save"
	"ssg-portal/controllers/section"
	"ssg-portal/controllers/semester"
	"ssg-portal/controllers/studententry"
	"ssg-portal/controllers/subject"
	"ssg-portal/controllers/subjectUpload"
	"ssg-portal/controllers/update"
	"ssg-portal/controllers/venue"

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
	r.GET("/timetable/:departmentID/:semesterId/:academicYearID/:sectionID", groups.GenerateTimetable)
	r.POST("/timetable/save", save.SaveTimetable)
	r.GET("/timetable/saved/:departmentID/:semesterId/:academicYearID/:sectionID", save.GetTimetable)
	r.GET("/timetable/faculty/:faculty_name/:academicYearID", facultyDetails.FacultyTimetable)
	r.GET("/timetable/lab/:subject_name", lab.LabTableTimetable)
	r.GET("/timetable/facultyOptions", facultyDetails.FacultyOptions)
	r.GET("/timetable/labOptions", lab.LabOptions)
	r.GET("/timetable/options", controllers.TimetableOptions)
	r.GET("/classroomOptions", controllers.VenueOptions)
	r.GET("/timetable/semoptions", semester.SemOptions)
	r.GET("/timetable/sectionoptions", section.SectionOptions)
	r.POST("/upload", subjectUpload.UploadDetails)
	r.GET("faculty/available/:departmentID/:semesterID/:day/:startTime/:endTime/:academicYearID/:sectionID", edit.GetAvailableFaculty)
	r.GET("/available-timings/:facultyName/:day", facultyDetails.GetAvailableTimingsForFaculty)
	r.GET("/saved/deptoptions", save.SavedDepartmentOptions)
	r.PUT("/timetable/update", update.UpdateTimetable)
	r.GET("/manual/options", manualentry.DayAndTimeOptions)
	r.POST("/manual/submit", manualentry.SubmitManualEntry)
	r.GET("/classroomDetailsOptions", classroom.ClassroomDetailsOptions)
	r.GET("/acdemicYearOptions", academicYear.AcademicYearOptions)
	r.GET("/periodallocation", allocation.Subjectallocation)
	r.PUT("/periodallocationedit", allocation.UpdateAllocation)
	r.GET("/getResource", role.GetResources)
	r.POST("/send-otp", otp.SendOTPHandler)
	r.POST("/verify-otp", otp.VerifyOTPHandler)
	r.GET("/downloadTimetable/:semesterId", excel.DownloadTimetable)
	r.GET("/venueTimetable/:classroom", venue.VenueTimetable)
	r.GET("/venueTimetableOptions", venue.ClassroomOptions)
	r.POST("/studententry/upload", studententry.InsertStudentEntries)
	r.GET("/studentTimetable/:studentID", studententry.StudentTimetable)
	r.GET("/download/:academic_year_id", excel.Masterdownload)
	r.GET("/studentoptions", studententry.GetStudentOptions)
	r.GET("/subjectoptions", subject.SubjectOptions)
	r.GET("/subjectTypeoptions", subject.SubjectTypeOptions)
	r.GET("/classroomavailabletimings/:academicYearID/:facultyName/:day/:classroomName", classroom.GetAvailableTimingsForFacultyAndClassroom)
	r.GET("/workingDayoptions", date.WorkingDayOptions)
	r.GET("/course-code", subject.CourseCodeOptions)
	r.POST("/labentry", labentry.LabEntry)
	r.POST("/manual/bulksubmit", manualentry.BulkInsert)
	r.POST("/upload-lab", lab.HandleExcelUpload)
	r.Run(":8080")

}
	