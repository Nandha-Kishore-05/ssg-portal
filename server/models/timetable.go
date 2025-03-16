package models

import "time"

type WorkingDay struct {
	WorkingDate time.Time `json:"working_date"`
	Day string `json:"day"`
}

type FacultyAssignment struct {
    FacultyID   int
    FacultyName string
    SectionIDs  map[int]bool
}

type Day struct {
	ID      int    `json:"id"`
	DayName string `json:"day_name"`
}

type Hour struct {
	ID        int    `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Value     string
}

type Subject struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DepartmentID int    `json:"department_id"`
	Period       int    `json:"periods"`
	Status       int    `json:"status"`
	SemesterID   int    `json:"semester_id"`
	CourseCode   string `json:"course_code"`
	SectionID       int    `json:"section_id"`
}
type Section struct {
	ID          int    `json:"id"`
	SectionName string `json:"section_name"`
	SemesterID   int    `json:"semester_id"`
	DepartmentID    int    `json:"department_id"`
	AcademicYear  int    `json:"academic_year_id"`
}
type Faculty struct {
	ID           int    `json:"id"`
	FacultyName  string `json:"name"`
	DepartmentID int    `json:"department_id"`
	SubjectName  string `json:"subject_name"`
	AcademicYear  int    `json:"academic_year_id"`
	SemesterID   int    `json:"semester_id"`
}

type Classroom struct {
	ID            int    `json:"id"`
	ClassroomName string `json:"name"`
	DepartmentID  int    `json:"department_id"`
	SemesterID    int    `json:"semester_id"`
	AcademicYear  int    `json:"academic_year_id"`
	SectionID     int    `json:"section_id"`
	Status        int    `json:"status"`
	LabVenue  string `json:"lab_name"`
}

type FacultySubject struct {
	FacultyID    int    `json:"faculty_id"`
	SubjectID    int    `json:"subject_id"`
	SemesterID   int    `json:"semester_id"`
	SubjectName  string `json:"subject_name"`
	SectionID    int    `json:"section_id"`
	AcademicYear int    `json:"academic_year_id"`
	DepartmentID int    `json:"department_id"`
	ClassroomID  int    `json:"classroom_id`
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
	TimetableStatus int    `json:"timetable_status"`
	AcademicYear    int    `json:"academic_year_id"`
	CourseCode      string `json:"course_code"`
	SectionID       int    `json:"section_id"`
	LabVenue  string `json:"lab_name"`
	
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
	Status         int    `json:"status"`
	AcademicYear   int    `json:"academic_year"`
	CourseCode     string `json:"course_code"`
	SectionID      int    `json:"section_id"`
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
	FacultyID      int    `json:"faculty_id"`
	DepartmentID   int    `json:"department_id"`
	SemesterID     int    `json:"semester_id"`
}
type Resource struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Path   string `json:"path"`
	SortBy int    `json:"sort_by"`
}
type EditRequest struct {
	ID             int    `json:"id"`
	SubjectName    string `json:"subject_name"`
	OldSubjectName string `json:"old_subject_name"`
	DepartmentID   int    `json:"department_id"`
	SemesterID     int    `json:"semester_id"`
}
type UpdateRequest struct {
	SubjectID      int    `json:"id"`
	Periods        int    `json:"periods"`
	SubjectName    string `json:"subject_name"`
	OldSubjectName string `json:"old_subject_name"`
	DepartmentID   int    `json:"department_id"`
	SemesterID     int    `json:"semester_id"`
	FacultyName    string `json:"faculty_name"`
	OldFacultyName string `json:"old_faculty_name"`
	FacultyID      *int   `json:"faculty_id"`
}
type AcademicYear struct {
	AcademicYear     int    `json:"academic_year_id"`
	AcademicYearName string `json:"academic_year"`
}
type VenueTimetable struct {
	DayName     string `json:"day_name"`
	StartTime   string `json:"start_time"` // Scanning start_time as string
	EndTime     string `json:"end_time"`
	SubjectName string `json:"subject_name"`
	FacultyName string `json:"faculty_name"`
	SemesterID  int    `json:"semester_id"`
	Department  string `json:"department_name"`
	SectionName string `json:"section_name"`
}
type Student struct {
	StudentName string `json:"Student Name"`
	RollNumber  string `json:"Roll Number"`
	CourseName  string `json:"Course Name"`
	CourseCode  string `json:"Course Code"`
}

type StudentEntryRequest struct {
	Students     []Student `json:"students"`
	Department   int       `json:"department"`
	Semester     int       `json:"semester"`
	AcademicYear int       `json:"academicYear"`
}
type StudentOptions struct {
	StudentID      int    `json:"student_id"`
	StudentName    string `json:"student_name"`
	StudentRollNo  string `json:"roll_no"`
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	SemesterID     int    `json:"semester_id"`
	SemesterName   string `json:"semester_name"`
	AcademicYearID int    `json:"academic_year_id"`
	AcademicYear   string `json:"academic_year"`
}
type StudentTimetable struct {
	ID             int    `json:"id"`
	DayName        string `json:"day_name"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Classroom      string `json:"classroom"`
	SubjectName    string `json:"subject_name"`
	FacultyName    string `json:"faculty_name"`
	Status         int    `json:"status"`
	StudentID      int    `json:"student_id"`
	CourseCode     string `json:"course_code"`
	AcademicYearID int    `json:"academic_year_id"`
	DepartmentID   int    `json:"department_id"`
	SemesterID     int    `json:"semester_id"`
}
type LabVenue struct {
	ID        int    `json:"id"`
	LabVenue  string `json:"lab_name"`
	SubjectID int    `json:"subject_id"`
	SemesterID    int    `json:"semester_id"`
	DepartmentID   int    `json:"department_id"`
}
