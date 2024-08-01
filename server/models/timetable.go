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
	LabCourse    int    `json:"lab_course"`
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
}

type FacultySubject struct {
	FacultyID int `json:"faculty_id"`
	SubjectID int `json:"subject_id"`
}
type TimetableEntry struct {
	Day       string `json:"day_name"` // Use DayName for consistency
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Subject   string `json:"subject_name"` // Use Subject for consistency
	Faculty   string `json:"faculty_name"` // Use Faculty for consistency
	Classroom string `json:"classroom"`
}
type Class struct {
	Day         string `json:"day_name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	SubjectName string `json:"subject_name"`
	FacultyName string `json:"faculty_name"`
	Classroom   string `json:"classroom"`
}
