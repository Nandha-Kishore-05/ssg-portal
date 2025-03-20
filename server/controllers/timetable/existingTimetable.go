package timetable

import (
	"fmt"
	"ssg-portal/config"
	"ssg-portal/models"
)

func FetchExistingTimetable() (map[string]map[string][]models.TimetableEntry, error) {
	existingTimetable := make(map[string]map[string][]models.TimetableEntry)

	rows, err := config.Database.Query(`
        SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom, status, semester_id, department_id, academic_year, course_code ,section_id
        FROM timetable`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dayName, startTime, endTime, subjectName, facultyName, classroom string
		var courseCode []byte
		var status, semesterID, departmentID, academicYearID, sectionID int

		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom, &status, &semesterID, &departmentID, &academicYearID, &courseCode, &sectionID); err != nil {
			return nil, err
		}

		courseCodeStr := string(courseCode)

		if _, exists := existingTimetable[facultyName]; !exists {
			existingTimetable[facultyName] = make(map[string][]models.TimetableEntry)
		}

		entry := models.TimetableEntry{
			DayName:      dayName,
			StartTime:    startTime,
			EndTime:      endTime,
			SubjectName:  subjectName,
			FacultyName:  facultyName,
			Classroom:    classroom,
			Status:       status,
			SemesterID:   semesterID,
			DepartmentID: departmentID,
			AcademicYear: academicYearID,
			CourseCode:   courseCodeStr,
			SectionID:    sectionID,
		}

		existingTimetable[facultyName][dayName] = append(existingTimetable[facultyName][dayName], entry)
	}

	return existingTimetable, nil
}


func FetchManualTimetable(departmentID int, semesterID int, academicYearID int, sectionID int) (map[string]map[string][]models.TimetableEntry, error) {

	manualTimetable := make(map[string]map[string][]models.TimetableEntry)

	query := `
        SELECT day_name, start_time, end_time, classroom, semester_id, department_id, 
               subject_name, faculty_name, status, academic_year, course_code, section_id 
        FROM manual_timetable 
        WHERE department_id = ? AND semester_id = ? AND academic_year = ? AND section_id = ?`

	rows, err := config.Database.Query(query, departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manual timetable: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.TimetableEntry
		if err := rows.Scan(&entry.DayName, &entry.StartTime, &entry.EndTime, &entry.Classroom,
			&entry.SemesterID, &entry.DepartmentID, &entry.SubjectName, &entry.FacultyName,
			&entry.Status, &entry.AcademicYear, &entry.CourseCode, &entry.SectionID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Initialize the nested map for the day if it doesn't exist
		if _, exists := manualTimetable[entry.DayName]; !exists {
			manualTimetable[entry.DayName] = make(map[string][]models.TimetableEntry)
		}

		manualTimetable[entry.DayName][entry.StartTime] = append(manualTimetable[entry.DayName][entry.StartTime], entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return manualTimetable, nil
}