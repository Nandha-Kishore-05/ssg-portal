package models

type Day struct {
	ID      int    `json:"id"`
	DayName string `json:"day_name"`
}

type Hour struct {
	ID        int    `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Subject struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DepartmentID int    `json:"department_id"`
	Period       int    `json:"periods"`
	Status       int    `json:"status"`
	SemesterID   int    `json:"semester_id"`
}

type Faculty struct {
	ID           int    `json:"id"`
	FacultyName  string `json:"name"`
	DepartmentID int    `json:"department_id"`
}

type Classroom struct {
	ID            int    `json:"id"`
	ClassroomName string `json:"name"`
	DepartmentID  int    `json:"department_id"`
	SemesterID    int    `json:"semester_id"`
}

type FacultySubject struct {
	FacultyID  int `json:"faculty_id"`
	SubjectID  int `json:"subject_id"`
	SemesterID int `json:"semester_id"`
}
type TimetableEntry struct {
	DayName     string `json:"day_name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	SubjectName string `json:"subject_name"`
	FacultyName string `json:"faculty_name"`
	Classroom   string `json:"classroom"`
	Status      int    `json:"status"`
	SemesterID  int    `json:"semester_id"`
	DepartmentID  int    `json:"department_id"`
}
type Class struct {
	Day          string `json:"day_name"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	SubjectName  string `json:"subject_name"`
	FacultyName  string `json:"faculty_name"`
	Classroom    string `json:"classroom"`
	DepartmentID int    `json:"department_id"`
}
type Department struct {
	ID         int    `json:"id"`
	Department string `json:"name"`
}
type FacultyTimetableEntry struct {
	DayName     string `json:"day_name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Classroom   string `json:"classroom"`
	SemesterID  int    `json:"semester_id"`
	SubjectName string `json:"subject_name"`
}
type Semester struct {
	ID           int
	SemesterName string
	YearID       int
}

type LabTimetableEntry struct {
	DayName     string `json:"day_name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Classroom   string `json:"classroom"`
	SemesterID  int    `json:"semester_id"`
	SubjectName string `json:"subject_name"`
	FacultyName string `json:"faculty_name"`
}
