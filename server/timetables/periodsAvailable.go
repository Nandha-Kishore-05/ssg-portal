package timetables

import (
	"fmt"
	"ssg-portal/config"
)

func PeriodsAvailable(departmentID, academicYearID, semesterID, sectionID int) error {
	var totalPeriods int

	// Check if the timetable already exists for the same department, academic year, semester, and section
	queryExistingTimetable := `
        SELECT COUNT(*) 
        FROM timetable 
        WHERE department_id = ? 
        AND academic_year = ? 
        AND semester_id = ? 
        AND section_id = ?`

	var existingTimetableCount int
	err := config.Database.QueryRow(queryExistingTimetable, departmentID, academicYearID, semesterID, sectionID).Scan(&existingTimetableCount)
	if err != nil {
		return fmt.Errorf("error checking existing timetable: %v", err)
	}

	// If a timetable already exists, return an error
	if existingTimetableCount > 0 {
		return fmt.Errorf("a timetable already exists for this combination of department, academic year, semester, and section")
	}

	// Query to fetch total periods from faculty_subjects
	queryFacultySubjects := `
        SELECT SUM(fs.periods) 
        FROM faculty_subjects fs
        WHERE fs.department_id = ? 
        AND fs.academic_year_id = ? 
        AND fs.semester_id = ? 
        AND fs.section_id = ?`

	var facultySubjectsPeriods int
	err = config.Database.QueryRow(queryFacultySubjects, departmentID, academicYearID, semesterID, sectionID).Scan(&facultySubjectsPeriods)
	if err != nil {
		return fmt.Errorf("error fetching periods from faculty_subjects: %v", err)
	}

	// Query to fetch total periods from manual_timetable
	queryManualTimetable := `
        SELECT COUNT(*) 
        FROM manual_timetable mt
        WHERE mt.department_id = ? 
        AND mt.academic_year = ? 
        AND mt.semester_id = ? 
        AND mt.section_id = ?`

	var manualTimetablePeriods int
	err = config.Database.QueryRow(queryManualTimetable, departmentID, academicYearID, semesterID, sectionID).Scan(&manualTimetablePeriods)
	if err != nil {
		return fmt.Errorf("error fetching periods from manual_timetable: %v", err)
	}

	// Query to fetch total periods from timetable_skips
	queryTimetableSkips := `
        SELECT COUNT(*) 
        FROM timetable_skips ts
        WHERE ts.department_id = ? 
        AND ts.academic_year = ? 
        AND ts.semester_id = ? 
        AND ts.section_id = ?`

	var timetableSkipsPeriods int
	err = config.Database.QueryRow(queryTimetableSkips, departmentID, academicYearID, semesterID, sectionID).Scan(&timetableSkipsPeriods)
	if err != nil {
		return fmt.Errorf("error fetching periods from timetable_skips: %v", err)
	}

	// Calculate the total periods
	totalPeriods = facultySubjectsPeriods + manualTimetablePeriods + timetableSkipsPeriods

	// Validate the total periods
	if totalPeriods < 36 {
		return fmt.Errorf("only %d periods available, 36 periods are required", totalPeriods)
	} else if totalPeriods > 36 {
		return fmt.Errorf("exceeded allowed periods: %d periods found (only 36 allowed)", totalPeriods)
	}
	return nil
}
