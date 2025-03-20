package timetable

import (
	"ssg-portal/config"
	"ssg-portal/models"
)

func FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID int) (map[string]map[string][]models.TimetableEntry, error) {
	skipEntries := make(map[string]map[string][]models.TimetableEntry)

	query := `
	SELECT day_name, start_time, end_time, subject_name, faculty_name, semester_id, department_id, classroom, status, academic_year, course_code, section_id
	FROM timetable_skips 
	WHERE department_id = ? AND semester_id = ? AND academic_year = ? AND section_id = ?`

	rows, err := config.Database.Query(query, departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dayName, startTime, endTime, subjectName, facultyName, classroom, courseCode string
		var semesterID, departmentID, status, academicYear, sectionID int

		// Scan values from the row
		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &semesterID, &departmentID, &classroom, &status, &academicYear, &courseCode, &sectionID); err != nil {
			return nil, err
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
			AcademicYear: academicYear,
			CourseCode:   courseCode,
			SectionID:    sectionID,
		}
		if skipEntries[dayName] == nil {
			skipEntries[dayName] = make(map[string][]models.TimetableEntry)
		}
		// Append the entry to the corresponding day and time
		skipEntries[dayName][startTime] = append(skipEntries[dayName][startTime], entry)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return skipEntries, nil
}