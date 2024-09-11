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
	SubjectName  string `json:"subject_name"`
}

type Classroom struct {
	ID            int    `json:"id"`
	ClassroomName string `json:"name"`
	DepartmentID  int    `json:"department_id"`
	SemesterID    int    `json:"semester_id"`
}

type FacultySubject struct {
	FacultyID   int    `json:"faculty_id"`
	SubjectID   int    `json:"subject_id"`
	SemesterID  int    `json:"semester_id"`
	SubjectName string `json:"subject_name"`
}
type TimetableEntry struct {
	ID              int    `json:"id"`
	DayName         string `json:"day_name"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	SubjectName     string `json:"subject_name"`
	FacultyName     string `json:"faculty_name"`
	Classroom       string `json:"classroom"`
	Status          int    `json:"status"`
	SemesterID      int    `json:"semester_id"`
	DepartmentID    int    `json:"department_id"`
	SubjectID       int    `json:"subject_id"`
	FacultyID       int    `json:"faculty_id"`
	TimetableStatus int    `json:timetable_status`
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
type Record struct {
	Department string `json:"Department"`
	Faculty    string `json:"Faculty"`
	LabSubject string `json:"Lab-subject"`
	Periods    int    `json:"Periods"`
	Semester   int    `json:"Semester"`
	Subject    string `json:"Subject"`
	Venue      string `json:"Venue"`
}
type ManualEntryRequest struct {
	Subject        string `json:"subject_name"`
	Department     int    `json:"department_id"`
	Semester       int    `json:"semester_id"`
	Day            string `json:"day_name"`
	DepartmentName string `json:"name"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Faculty        string `json:"faculty_name"`
	Classroom      string `json:"classroom"`
	Status         string `json:"status"`
	// Lab            int    `json:"status"`
}
type BulkManualEntryRequest struct {
	Entries []ManualEntryRequest `json:"entries"`
}
type MenuItem struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Path  string `json:"path"`
	Icon  string `json:"icon"`
}
type SubjectInfo struct {
	SubjectID      int    `json:"id"`
	SubjectName    string `json:"subject_name"`
	DepartmentName string `json:"department_name"`
	SemesterName   string `json:"semester_name"`
	Periods        int    `json:"periods"`
	Status         string `json:"status"`
	FacultyName    string `json:"faculty_name"`
	FacultyID      int   `json:"faculty_id"`
    DepartmentID    int    `json:"department_id"`
    SemesterID      int    `json:"semester_id"`
}
type Resource struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Icon   string  `json:"icon"`
	Path   string  `json:"path"`
	SortBy int     `json:"sort_by"`
}
type EditRequest struct {
	ID              int    `json:"id"`
	SubjectName     string `json:"subject_name"`
	OldSubjectName  string `json:"old_subject_name"`
	DepartmentID    int    `json:"department_id"`
	SemesterID      int    `json:"semester_id"`
}
type UpdateRequest struct {
	SubjectID     int    `json:"id"`
	Periods       int    `json:"periods"`
	SubjectName   string `json:"subject_name"`
	OldSubjectName string `json:"old_subject_name"`
	DepartmentID  int    `json:"department_id"`
	SemesterID    int    `json:"semester_id"`
	FacultyName    string `json:"faculty_name"`
	OldFacultyName  string `json:"old_faculty_name"`
	 FacultyID      *int   `json:"faculty_id"`
}
