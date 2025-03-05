package timetables

import (
	"fmt"
	"log"
	"math/rand"
	"ssg-portal/config"
	"ssg-portal/models"
	"time"
)

type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

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

// func GenerateTimetable(
// 	workingDays []models.WorkingDay,
// 	hours []models.Hour,
// 	subjects []models.Subject,
// 	faculty []models.Faculty,
// 	classrooms []models.Classroom,
// 	facultySubjects []models.FacultySubject,
// 	semesters []models.Semester,
// 	sections []models.Section,
// 	academicYear []models.AcademicYear,
// 	departmentID, semesterID, academicYearID, sectionID int) map[string]map[string][]models.TimetableEntry {

// 	// Initialize data structures
// 	var days []string
// 	for _, wd := range workingDays {
// 		days = append(days, wd.WorkingDate.Format("2006-01-02"))
// 	}

// 	// Check if this is a multi-section department for the same year and semester
// 	var sectionsInSameSemester []models.Section
// 	for _, section := range sections {
// 		if section.DepartmentID == departmentID && section.SemesterID == semesterID && section.AcademicYear == academicYearID {
// 			sectionsInSameSemester = append(sectionsInSameSemester, section)
// 		}
// 	}

// 	var count int
// 	query := `
// 		SELECT COUNT(*) FROM timetable
// 		WHERE department_id = ? AND semester_id = ? AND academic_year = ?
// 	`
// 	err := config.Database.QueryRow(query, departmentID, semesterID, academicYearID).Scan(&count)
// 	if err != nil {
// 		log.Println("Error executing query:", err)
// 		return nil
// 	}

// 	if count == 0 {
// 		return generateRandomTimetable(workingDays, hours, subjects, faculty, classrooms, facultySubjects, sections, semesters, departmentID, semesterID, academicYearID, sectionID)
// 	}

// 	// Fetch any one section ID
// 	var existingSectionID int
// 	sectionQuery := `
// 		SELECT section_id FROM timetable
// 		WHERE department_id = ? AND semester_id = ? AND academic_year = ?
// 		LIMIT 1
// 	`
// 	err = config.Database.QueryRow(sectionQuery, departmentID, semesterID, academicYearID).Scan(&existingSectionID)
// 	if err != nil {
// 		log.Println("Error fetching section:", err)
// 		return nil
// 	}

// 	// Check if sectionsInSameSemester has elements before accessing index 0
// 	if len(sectionsInSameSemester) > 0 && existingSectionID == sectionsInSameSemester[0].ID {
// 		// Fetch the existing timetable and handle errors
// 		existingTimetable, err := FetchExistingTimetable()
// 		if err != nil {
// 			fmt.Println("Error fetching existing timetable:", err)
// 			return nil
// 		}

// 		// Fetch lab venues, skips, and manual timetable entries
// 		labVenues, err := GetLabVenue()
// 		if err != nil {
// 			fmt.Println("Error fetching lab venues:", err)
// 			return nil
// 		}

// 		skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
// 		if err != nil {
// 			fmt.Println("Error fetching timetable skips:", err)
// 			return nil
// 		}

// 		// Map subject to classrooms and faculty to subjects
// 		subjectClassrooms := map[string][]models.Classroom{}
// 		for _, subject := range subjects {
// 			for _, cls := range classrooms {
// 				if cls.DepartmentID == subject.DepartmentID {
// 					subjectClassrooms[subject.Name] = append(subjectClassrooms[subject.Name], cls)
// 				}
// 			}
// 		}

// 		facultySubjectMap := map[int]map[int]bool{}
// 		for _, fs := range facultySubjects {
// 			if facultySubjectMap[fs.FacultyID] == nil {
// 				facultySubjectMap[fs.FacultyID] = map[int]bool{}
// 			}
// 			facultySubjectMap[fs.FacultyID][fs.SubjectID] = true
// 		}

// 		// Generate timetable using a randomized approach with max attempts
// 		for {
// 			timetable := make(map[string]map[string][]models.TimetableEntry)
// 			subjectsAssigned := make(map[string]map[string]bool)
// 			periodsLeft := make(map[string]int)
// 			status0Assignments := make(map[string]map[string]bool)
// 			facultyAssignments := make(map[string]map[string]string)
// 			 facultyDailyCount := make(map[string]map[string]int)
// 			labAssigned := make(map[string]bool)

// 			for _, subject := range subjects {
// 				periodsLeft[subject.Name] = subject.Period
// 				if subject.Status == 0 {
// 					status0Assignments[subject.Name] = make(map[string]bool)
// 				}
// 			}

// 			for _, day := range days {
// 				timetable[day] = make(map[string][]models.TimetableEntry)
// 				subjectsAssigned[day] = make(map[string]bool)
// 				facultyAssignments[day] = make(map[string]string)
// 				facultyDailyCount[day] = make(map[string]int)
// 				labAssigned[day] = false
// 				if skips, ok := skipTimetable[day]; ok {
// 					for startTime := range skips {
// 						timetable[day][startTime] = append(timetable[day][startTime])
// 					}
// 				}
// 			}

// 			rand.Seed(time.Now().UnixNano())

// 			for _, day := range days {
// 				for i := 0; i < len(hours); i++ {
// 					assigned := false
// 					for attempts := 0; attempts < 1000; attempts++ {
// 						var filteredSubjects []models.Subject
// 						for _, subject := range subjects {
// 							if periodsLeft[subject.Name] > 0 && (!subjectsAssigned[day][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) {

// 								if subject.Status == 0 && labAssigned[day] {
// 									continue
// 								}
// 								if subject.Status == 1 && subjectsAssigned[day][subject.Name] {
// 									continue
// 								}

// 								var validClassrooms []models.Classroom
// 								for _, cls := range classrooms {
// 									if cls.DepartmentID == subject.DepartmentID {
// 										for _, semester := range semesters {
// 											if semester.ID == cls.SemesterID {
// 												validClassrooms = append(validClassrooms, cls)
// 												break
// 											}
// 										}
// 									}
// 								}

// 								if len(validClassrooms) > 0 {
// 									filteredSubjects = append(filteredSubjects, subject)
// 								}
// 							}
// 						}

// 						if len(filteredSubjects) == 0 {
// 							continue
// 						}

// 						subjectIndex := rand.Intn(len(filteredSubjects))
// 						subject := filteredSubjects[subjectIndex]

// 						hourIndex := i % len(hours)
// 						startTime := hours[hourIndex].StartTime
// 						endTime := hours[hourIndex].EndTime

// 						if _, ok := timetable[day][startTime]; ok {
// 							if len(timetable[day][startTime]) > 0 {
// 								continue
// 							}
// 						}

// 						var availableFaculty []models.Faculty
// 						for _, fac := range faculty {
// 							for _, fs := range facultySubjects {
// 								if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID &&
// 									fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
// 									fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
// 									availableFaculty = append(availableFaculty, fac)
// 									break
// 								}
// 							}
// 						}

// 						if len(availableFaculty) == 0 {
// 							continue
// 						}

// 						facultyIndex := rand.Intn(len(availableFaculty))
// 						selectedFaculty := availableFaculty[facultyIndex]

// 						if facultyDailyCount[day][selectedFaculty.FacultyName] >= 2 {
// 							continue
// 						}

// 						if assignedClassroom, exists := facultyAssignments[day][selectedFaculty.FacultyName]; exists && assignedClassroom == startTime {
// 							continue
// 						}
// 						var selectedClassroom models.Classroom
// 						if subject.Status == 0 && len(labVenues) > 0 { // Check lab venues for lab subjects
// 							// Use the LabVenue for lab subjects
// 							selectedLabVenue := labVenues[rand.Intn(len(labVenues))]
// 							selectedClassroom = models.Classroom{
// 								ID:            selectedLabVenue.ID,
// 								ClassroomName: selectedLabVenue.LabVenue, // Map LabVenue fields to Classroom-like structure
// 								DepartmentID:  selectedLabVenue.DepartmentID,
// 								SemesterID:    selectedLabVenue.SemesterID,
// 							}
// 						} else {
// 							// Regular classroom assignment
// 							for _, cls := range classrooms {
// 								if cls.DepartmentID == subject.DepartmentID {
// 									for _, semester := range semesters {
// 										if semester.ID == cls.SemesterID {
// 											selectedClassroom = cls
// 											break
// 										}
// 									}
// 									break
// 								}
// 							}
// 						}

// 						entry := models.TimetableEntry{
// 							DayName:      day,
// 							StartTime:    startTime,
// 							EndTime:      endTime,
// 							SubjectName:  subject.Name,
// 							FacultyName:  selectedFaculty.FacultyName,
// 							Classroom:    selectedClassroom.ClassroomName,
// 							Status:       subject.Status,
// 							SemesterID:   semesterID,
// 							DepartmentID: departmentID,
// 							AcademicYear: academicYearID,
// 							CourseCode:   subject.CourseCode,
// 							SectionID:    sectionID,
// 						}

// 						if _, ok := timetable[day]; !ok {
// 							timetable[day] = make(map[string][]models.TimetableEntry)
// 						}
// 						if _, ok := timetable[day][startTime]; !ok {
// 							timetable[day][startTime] = []models.TimetableEntry{}
// 						}

// 						if subject.Status == 0 {
// 							if _, ok := status0Assignments[subject.Name][startTime]; !ok {
// 								if i < len(hours)-1 {
// 									nextHourIndex := (hourIndex + 1) % len(hours)
// 									nextStartTime := hours[nextHourIndex].StartTime
// 									if IsPeriodAvailable(existingTimetable, day, nextStartTime, "") {
// 										entry2 := entry
// 										entry2.StartTime = nextStartTime
// 										entry2.EndTime = hours[nextHourIndex].EndTime

// 										timetable[day][startTime] = append(timetable[day][startTime], entry)
// 										timetable[day][nextStartTime] = append(timetable[day][nextStartTime], entry2)
// 										periodsLeft[subject.Name] -= 2
// 										subjectsAssigned[day][subject.Name] = true
// 										status0Assignments[subject.Name][startTime] = true
// 										status0Assignments[subject.Name][nextStartTime] = true
// 										facultyAssignments[day][selectedFaculty.FacultyName] = nextStartTime
// 										facultyDailyCount[day][selectedFaculty.FacultyName] += 2
// 										labAssigned[day] = true
// 										assigned = true
// 										break
// 									}
// 								}
// 								continue
// 							}
// 						}

// 						timetable[day][startTime] = append(timetable[day][startTime], entry)
// 						periodsLeft[subject.Name]--
// 						facultyDailyCount[day][selectedFaculty.FacultyName]++
// 						subjectsAssigned[day][subject.Name] = true
// 						if subject.Status == 0 {
// 							status0Assignments[subject.Name][startTime] = true
// 						}
// 						facultyAssignments[day][selectedFaculty.FacultyName] = startTime
// 						if subject.Status == 0 { // Mark as lab subject assigned
// 							labAssigned[day] = true
// 						}
// 						assigned = true
// 						break
// 					}
// 					if !assigned {
// 						// fmt.Printf("Warning: Could not assign a subject for %s during period %d\n", day, i+1)
// 					}
// 				}
// 			}

// 			allAssigned := true
// 			for subjectName, remainingPeriods := range periodsLeft {
// 				if remainingPeriods > 0 {
// 					fmt.Printf("Warning: Subject %s has %d periods left unassigned.\n", subjectName, remainingPeriods)
// 					allAssigned = false
// 				}
// 			}

// 			// Check if all periods are filled
// 			periodsFilled := true
// 			for _, day := range days {
// 				for _, hour := range hours {
// 					startTime := hour.StartTime
// 					if len(timetable[day][startTime]) == 0 {
// 						periodsFilled = false
// 						break
// 					}
// 				}
// 				if !periodsFilled {
// 					break
// 				}
// 			}

// 			// Return timetable if all constraints are satisfied
// 			if allAssigned && periodsFilled {
// 				return timetable
// 			}
// 		}
// 	}

func GenerateTimetable(
	workingDays []models.WorkingDay,
	hours []models.Hour,
	subjects []models.Subject,
	faculty []models.Faculty,
	classrooms []models.Classroom,
	facultySubjects []models.FacultySubject,
	semesters []models.Semester,
	sections []models.Section,
	academicYear []models.AcademicYear,
	departmentID, semesterID, academicYearID, sectionID int) map[string]map[string][]models.TimetableEntry {

	// Initialize data structures
	var days []string
	for _, wd := range workingDays {
		days = append(days, wd.WorkingDate.Format("2006-01-02"))
	}

	// Check if this is a multi-section department for the same year and semester
	var sectionsInSameSemester []models.Section
	for _, section := range sections {
		if section.DepartmentID == departmentID && section.SemesterID == semesterID && section.AcademicYear == academicYearID {
			sectionsInSameSemester = append(sectionsInSameSemester, section)
		}
	}

	var count int
	query := `
		SELECT COUNT(*) FROM timetable 
		WHERE department_id = ? AND semester_id = ? AND academic_year = ?
	`
	err := config.Database.QueryRow(query, departmentID, semesterID, academicYearID).Scan(&count)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil
	}

	if count == 0 {
		return generateRandomTimetable(workingDays, hours, subjects, faculty, classrooms, facultySubjects, sections, semesters, departmentID, semesterID, academicYearID, sectionID)
	}

	// Fetch any one section ID
	var existingSectionID int
	sectionQuery := `
		SELECT section_id FROM timetable 
		WHERE department_id = ? AND semester_id = ? AND academic_year = ?
		LIMIT 1
	`
	err = config.Database.QueryRow(sectionQuery, departmentID, semesterID, academicYearID).Scan(&existingSectionID)
	if err != nil {
		log.Println("Error fetching section:", err)
		return nil
	}

	// Check if sectionsInSameSemester has elements before accessing index 0
	if len(sectionsInSameSemester) > 0 && existingSectionID == sectionsInSameSemester[0].ID {
		// Fetch the existing timetable and handle errors
		existingTimetable, err := FetchExistingTimetable()
		if err != nil {
			fmt.Println("Error fetching existing timetable:", err)
			return nil
		}

		// Fetch lab venues, skips, and manual timetable entries
		labVenues, err := GetLabVenue()
		if err != nil {
			fmt.Println("Error fetching lab venues:", err)
			return nil
		}

		skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
		if err != nil {
			fmt.Println("Error fetching timetable skips:", err)
			return nil
		}

		// Map subject to classrooms and faculty to subjects
		subjectClassrooms := map[string][]models.Classroom{}
		for _, subject := range subjects {
			for _, cls := range classrooms {
				if cls.DepartmentID == subject.DepartmentID {
					subjectClassrooms[subject.Name] = append(subjectClassrooms[subject.Name], cls)
				}
			}
		}

		facultySubjectMap := map[int]map[int]bool{}
		for _, fs := range facultySubjects {
			if facultySubjectMap[fs.FacultyID] == nil {
				facultySubjectMap[fs.FacultyID] = map[int]bool{}
			}
			facultySubjectMap[fs.FacultyID][fs.SubjectID] = true
		}

		// Generate timetable using a randomized approach with max attempts
		for {
			timetable := make(map[string]map[string][]models.TimetableEntry)
			subjectsAssigned := make(map[string]map[string]bool)
			periodsLeft := make(map[string]int)
			status0Assignments := make(map[string]map[string]bool)
			facultyAssignments := make(map[string]map[string]string)
			subjectDailyCount := make(map[string]map[string]int)
			labAssigned := make(map[string]bool)

			for _, subject := range subjects {
				periodsLeft[subject.Name] = subject.Period
				if subject.Status == 0 {
					status0Assignments[subject.Name] = make(map[string]bool)
				}
			}

			for _, day := range days {
				timetable[day] = make(map[string][]models.TimetableEntry)
				subjectsAssigned[day] = make(map[string]bool)
				facultyAssignments[day] = make(map[string]string)
				subjectDailyCount[day] = make(map[string]int)
				labAssigned[day] = false
				if skips, ok := skipTimetable[day]; ok {
					for startTime := range skips {
						timetable[day][startTime] = append(timetable[day][startTime])
					}
				}
			}

			rand.Seed(time.Now().UnixNano())

			for _, day := range days {
				for i := 0; i < len(hours); i++ {
					assigned := false
					for attempts := 0; attempts < 1000; attempts++ {
						var filteredSubjects []models.Subject
						for _, subject := range subjects {
							// Adjust daily count constraint based on subject status
							dailyLimit := 1
							if subject.Status == 0 { // Lab subject
								dailyLimit = 2
							}

							if periodsLeft[subject.Name] > 0 &&
								subjectDailyCount[day][subject.Name] < dailyLimit &&
								(!subjectsAssigned[day][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) {

								if subject.Status == 0 && labAssigned[day] {
									continue
								}
								if subject.Status == 1 && subjectsAssigned[day][subject.Name] {
									continue
								}

								var validClassrooms []models.Classroom
								for _, cls := range classrooms {
									if cls.DepartmentID == subject.DepartmentID {
										for _, semester := range semesters {
											if semester.ID == cls.SemesterID {
												validClassrooms = append(validClassrooms, cls)
												break
											}
										}
									}
								}

								if len(validClassrooms) > 0 {
									filteredSubjects = append(filteredSubjects, subject)
								}
							}
						}

						if len(filteredSubjects) == 0 {
							continue
						}

						subjectIndex := rand.Intn(len(filteredSubjects))
						subject := filteredSubjects[subjectIndex]

						hourIndex := i % len(hours)
						startTime := hours[hourIndex].StartTime
						endTime := hours[hourIndex].EndTime

						if _, ok := timetable[day][startTime]; ok {
							if len(timetable[day][startTime]) > 0 {
								continue
							}
						}

						var availableFaculty []models.Faculty
						for _, fac := range faculty {
							for _, fs := range facultySubjects {
								if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID &&
									fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
									fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
									availableFaculty = append(availableFaculty, fac)
									break
								}
							}
						}

						if len(availableFaculty) == 0 {
							continue
						}

						facultyIndex := rand.Intn(len(availableFaculty))
						selectedFaculty := availableFaculty[facultyIndex]

						if assignedClassroom, exists := facultyAssignments[day][selectedFaculty.FacultyName]; exists && assignedClassroom == startTime {
							continue
						}
						var selectedClassroom models.Classroom
						if subject.Status == 0 && len(labVenues) > 0 { // Check lab venues for lab subjects
							// Use the LabVenue for lab subjects
							selectedLabVenue := labVenues[rand.Intn(len(labVenues))]
							selectedClassroom = models.Classroom{
								ID:            selectedLabVenue.ID,
								ClassroomName: selectedLabVenue.LabVenue, // Map LabVenue fields to Classroom-like structure
								DepartmentID:  selectedLabVenue.DepartmentID,
								SemesterID:    selectedLabVenue.SemesterID,
							}
						} else {
							// Regular classroom assignment
							for _, cls := range classrooms {
								if cls.DepartmentID == subject.DepartmentID {
									for _, semester := range semesters {
										if semester.ID == cls.SemesterID {
											selectedClassroom = cls
											break
										}
									}
									break
								}
							}
						}

						entry := models.TimetableEntry{
							DayName:      day,
							StartTime:    startTime,
							EndTime:      endTime,
							SubjectName:  subject.Name,
							FacultyName:  selectedFaculty.FacultyName,
							Classroom:    selectedClassroom.ClassroomName,
							Status:       subject.Status,
							SemesterID:   semesterID,
							DepartmentID: departmentID,
							AcademicYear: academicYearID,
							CourseCode:   subject.CourseCode,
							SectionID:    sectionID,
						}

						if _, ok := timetable[day]; !ok {
							timetable[day] = make(map[string][]models.TimetableEntry)
						}
						if _, ok := timetable[day][startTime]; !ok {
							timetable[day][startTime] = []models.TimetableEntry{}
						}

						if subject.Status == 0 {
							if _, ok := status0Assignments[subject.Name][startTime]; !ok {
								if i < len(hours)-1 {
									nextHourIndex := (hourIndex + 1) % len(hours)
									nextStartTime := hours[nextHourIndex].StartTime
									if IsPeriodAvailable(existingTimetable, day, nextStartTime, "") {
										entry2 := entry
										entry2.StartTime = nextStartTime
										entry2.EndTime = hours[nextHourIndex].EndTime

										timetable[day][startTime] = append(timetable[day][startTime], entry)
										timetable[day][nextStartTime] = append(timetable[day][nextStartTime], entry2)
										periodsLeft[subject.Name] -= 2
										subjectsAssigned[day][subject.Name] = true
										status0Assignments[subject.Name][startTime] = true
										status0Assignments[subject.Name][nextStartTime] = true
										facultyAssignments[day][selectedFaculty.FacultyName] = nextStartTime
										subjectDailyCount[day][subject.Name] += 2
										labAssigned[day] = true
										assigned = true
										break
									}
								}
								continue
							}
						}

						timetable[day][startTime] = append(timetable[day][startTime], entry)
						periodsLeft[subject.Name]--
						subjectDailyCount[day][subject.Name]++
						subjectsAssigned[day][subject.Name] = true
						if subject.Status == 0 {
							status0Assignments[subject.Name][startTime] = true
						}
						facultyAssignments[day][selectedFaculty.FacultyName] = startTime
						if subject.Status == 0 { // Mark as lab subject assigned
							labAssigned[day] = true
						}
						assigned = true
						break
					}
					if !assigned {
						// fmt.Printf("Warning: Could not assign a subject for %s during period %d\n", day, i+1)
					}
				}
			}

			allAssigned := true
			for subjectName, remainingPeriods := range periodsLeft {
				if remainingPeriods > 0 {
					fmt.Printf("Warning: Subject %s has %d periods left unassigned.\n", subjectName, remainingPeriods)
					allAssigned = false
				}
			}

			// Check if all periods are filled
			periodsFilled := true
			for _, day := range days {
				for _, hour := range hours {
					startTime := hour.StartTime
					if len(timetable[day][startTime]) == 0 {
						periodsFilled = false
						break
					}
				}
				if !periodsFilled {
					break
				}
			}

			// Return timetable if all constraints are satisfied
			if allAssigned && periodsFilled {
				return timetable
			}
		}
	} else {
		// For sections other than the first, get the section ID to use as reference
		// Fix: Handle case where sectionsInSameSemester is empty
		var firstSectionID int

		// Check if there are any sections in the same semester
		if len(sectionsInSameSemester) > 0 {
			firstSectionID = sectionsInSameSemester[0].ID
		} else {
			// If no sections found in the same semester, use the existing section ID or generate new timetable
			if existingSectionID > 0 {
				firstSectionID = existingSectionID
			} else {
				// No reference section found, generate a new timetable
				return generateRandomTimetable(workingDays, hours, subjects, faculty, classrooms, facultySubjects, sections, semesters, departmentID, semesterID, academicYearID, sectionID)
			}
		}

		// Fetch the first section's timetable
		firstSectionTimetable, err := FetchSectionTimetable(departmentID, semesterID, academicYearID, firstSectionID)
		if err != nil {
			fmt.Println("Error fetching first section timetable:", err)
			return nil
		}

		labVenues, err := GetLabVenue()
		if err != nil {
			fmt.Println("Error fetching lab venues:", err)
			return nil
		}

		skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
		if err != nil {
			fmt.Println("Error fetching timetable skips:", err)
			return nil
		}

		// Map subject to classrooms and faculty to subjects
		subjectClassrooms := map[string][]models.Classroom{}
		for _, subject := range subjects {
			for _, cls := range classrooms {
				if cls.DepartmentID == subject.DepartmentID {
					subjectClassrooms[subject.Name] = append(subjectClassrooms[subject.Name], cls)
				}
			}
		}

		if len(firstSectionTimetable) == 0 {
			fmt.Println("First section timetable not found, generating new timetable")
			return generateRandomTimetable(workingDays, hours, subjects, faculty, classrooms, facultySubjects, sections, semesters, departmentID, semesterID, academicYearID, sectionID)
		}

		// Create a new timetable based on first section's schedule
		timetable := make(map[string]map[string][]models.TimetableEntry)

		// Initialize timetable structure
		for _, day := range days {
			timetable[day] = make(map[string][]models.TimetableEntry)
		}

		// Track faculty assignments to avoid conflicts
		facultyAssignments := make(map[string]map[string]string)   // day -> faculty -> timeslot
		classroomAssignments := make(map[string]map[string]string) // day -> classroom -> timeslot
		//facultyDailyCount := make(map[string]map[string]int)       // day -> faculty -> count
		labVenueAssignments := make(map[string]map[string]string)

		// Initialize tracking structures
		for _, day := range days {
			facultyAssignments[day] = make(map[string]string)
			classroomAssignments[day] = make(map[string]string)
			//facultyDailyCount[day] = make(map[string]int)
			labVenueAssignments[day] = make(map[string]string)
			if skips, ok := skipTimetable[day]; ok {
				for startTime := range skips {
					timetable[day][startTime] = append(timetable[day][startTime])
				}
			}
		}

		// Iterate through first section's timetable and create entries for current section
		for day, timeSlots := range firstSectionTimetable {
			for startTime, entries := range timeSlots {
				for _, entry := range entries {
					// Find subject ID from name
					var subjectID int
					//var subjectStatus int
					for _, subj := range subjects {
						if subj.Name == entry.SubjectName {
							subjectID = subj.ID
							//subjectStatus = subj.Status
							break
						}
					}

					if skips, ok := skipTimetable[day]; ok {
						if skippedEntries, exists := skips[startTime]; exists {
							// Use the skipTimetable directly
							timetable[day][startTime] = skippedEntries
							continue // Ensure that normal timetable generation is skipped
						}
					}
					var filteredSubjects []models.Subject

					for _, subject := range subjects {

						var validClassrooms []models.Classroom
						for _, cls := range classrooms {
							if cls.DepartmentID == subject.DepartmentID {
								for _, semester := range semesters {
									if semester.ID == cls.SemesterID {
										validClassrooms = append(validClassrooms, cls)
										break
									}
								}
							}
						}
						if len(validClassrooms) > 0 {
							filteredSubjects = append(filteredSubjects, subject)
						}
					}

					if len(filteredSubjects) == 0 {
						continue
					}

					subjectIndex := rand.Intn(len(filteredSubjects))
					subject := filteredSubjects[subjectIndex]
					// Find a suitable faculty for this subject
					var availableFaculty []models.Faculty
					for _, fac := range faculty {
						// Check if faculty is already assigned at this time
						if assignedTime, exists := facultyAssignments[day][fac.FacultyName]; exists && assignedTime == startTime {
							continue
						}

						// Check if faculty has reached daily limit
						// if facultyDailyCount[day][fac.FacultyName] >= 2 {
						// 	continue
						// }

						// Check if faculty can teach this subject
						canTeach := false
						for _, fs := range facultySubjects {
							if fs.FacultyID == fac.ID && fs.SubjectID == subjectID &&
								fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
								fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
								canTeach = true
								break
							}
						}

						if canTeach {
							availableFaculty = append(availableFaculty, fac)
						}
					}

					// If no faculty available, try to find any faculty
					if len(availableFaculty) == 0 {
						fmt.Printf("Warning: No preferred faculty for subject %s in section %d. Trying alternative faculty.\n", entry.SubjectName, sectionID)

						for _, fac := range faculty {
							// Check if faculty is already assigned at this time
							if assignedTime, exists := facultyAssignments[day][fac.FacultyName]; exists && assignedTime == startTime {
								continue
							}

							// Check if faculty has reached daily limit
							// if facultyDailyCount[day][fac.FacultyName] >= 2 {
							// 	continue
							// }

							availableFaculty = append(availableFaculty, fac)
						}
					}

					// If still no faculty available, create placeholder
					var selectedFaculty models.Faculty
					if len(availableFaculty) == 0 {
						fmt.Printf("Warning: No available faculty for subject %s in section %d\n", entry.SubjectName, sectionID)
						selectedFaculty = models.Faculty{
							FacultyName: "To Be Assigned",
						}
					} else {
						// Randomly select a faculty
						rand.Seed(time.Now().UnixNano())
						facultyIndex := rand.Intn(len(availableFaculty))
						selectedFaculty = availableFaculty[facultyIndex]

						// Mark faculty as assigned
						facultyAssignments[day][selectedFaculty.FacultyName] = startTime
						// facultyDailyCount[day][selectedFaculty.FacultyName]++
					}

					// Select appropriate classroom based on subject type
					var selectedClassroom models.Classroom
					if entry.Status == 0 && len(labVenues) > 0 { // Check lab subjects
						var selectedLabVenue models.LabVenue
						availableLabVenues := []models.LabVenue{}

						// Filter out already assigned lab venues for the same day and time
						for _, lab := range labVenues {
							if labVenueAssignments[day][entry.StartTime] != lab.LabVenue {
								availableLabVenues = append(availableLabVenues, lab)
							}
						}

						if len(availableLabVenues) > 0 {
							selectedLabVenue = availableLabVenues[rand.Intn(len(availableLabVenues))]
							labVenueAssignments[day][entry.StartTime] = selectedLabVenue.LabVenue // Mark as assigned
						} else {
							fmt.Println("Warning: No available lab venue for", day, entry.StartTime)
							continue // Skip this period if no lab venue is available
						}

						// Assign the selected lab venue
						selectedClassroom = models.Classroom{
							ID:            selectedLabVenue.ID,
							ClassroomName: selectedLabVenue.LabVenue,
							DepartmentID:  selectedLabVenue.DepartmentID,
							SemesterID:    selectedLabVenue.SemesterID,
						}
					} else {
						// Regular classroom assignment
						for _, cls := range classrooms {
							if cls.DepartmentID == subject.DepartmentID {
								for _, semester := range semesters {
									if semester.ID == cls.SemesterID {
										selectedClassroom = cls
										break
									}
								}
								break
							}
						}
					}

					// Create a new entry for current section
					newEntry := models.TimetableEntry{
						DayName:      day,
						StartTime:    entry.StartTime,
						EndTime:      entry.EndTime,
						SubjectName:  entry.SubjectName,
						FacultyName:  selectedFaculty.FacultyName,
						Classroom:    selectedClassroom.ClassroomName,
						Status:       entry.Status,
						SemesterID:   semesterID,
						DepartmentID: departmentID,
						AcademicYear: academicYearID,
						CourseCode:   entry.CourseCode,
						SectionID:    sectionID,
					}

					if _, ok := timetable[day]; !ok {
						timetable[day] = make(map[string][]models.TimetableEntry)
					}
					if _, ok := timetable[day][startTime]; !ok {
						timetable[day][startTime] = []models.TimetableEntry{}
					}

					timetable[day][startTime] = append(timetable[day][startTime], newEntry)

					// If this is a lab subject that spans two periods, handle the second period as well

				}
			}
		}

		// Verify that all periods have assignments
		for _, day := range days {
			for _, hour := range hours {
				startTime := hour.StartTime
				if _, ok := timetable[day][startTime]; !ok || len(timetable[day][startTime]) == 0 {
					fmt.Printf("Warning: No assignment for %s at %s in section %d\n", day, startTime, sectionID)
				}
			}
		}

		return timetable
	}
}

// New function to fetch a specific section's timetable
func FetchSectionTimetable(departmentID, semesterID, academicYearID, sectionID int) (map[string]map[string][]models.TimetableEntry, error) {
	sectionTimetable := make(map[string]map[string][]models.TimetableEntry)

	rows, err := config.Database.Query(`
		SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom, status, semester_id, department_id, academic_year, course_code, section_id
		FROM timetable
		WHERE department_id = ? AND semester_id = ? AND academic_year = ? AND section_id = ?`,
		departmentID, semesterID, academicYearID, sectionID)

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

		if _, exists := sectionTimetable[dayName]; !exists {
			sectionTimetable[dayName] = make(map[string][]models.TimetableEntry)
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

		sectionTimetable[dayName][startTime] = append(sectionTimetable[dayName][startTime], entry)
	}

	return sectionTimetable, nil
}

// Verify that a period is available in the existing timetable
func IsPeriodAvailable(existingTimetable map[string]map[string][]models.TimetableEntry, day, startTime, facultyName string) bool {
	return true
}

func generateRandomTimetable(
	workingDays []models.WorkingDay, // Valid working days to consider
	hours []models.Hour,
	subjects []models.Subject,
	faculty []models.Faculty,
	classrooms []models.Classroom,
	facultySubjects []models.FacultySubject,
	section []models.Section,
	semesters []models.Semester,
	departmentID int,
	semesterID int,
	academicYearID int,
	sectionID int,
) FacultyBasedTimetable {
	log.Println("Academic Year ID:", academicYearID)

	// Extract working days into a slice of strings
	var days []string
	for _, wd := range workingDays {
		days = append(days, wd.WorkingDate.Format("2006-01-02")) // Use date in "YYYY-MM-DD" format
	}

	labVenues, err := GetLabVenue()
	if err != nil {
		fmt.Println("Error fetching lab venues:", err)
		return nil
	}

	// Fetch timetable skips
	skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		fmt.Println("Error fetching timetable skips:", err)
		return nil
	}

	// Fetch manual timetable
	manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		fmt.Println("Error fetching manual timetable:", err)
		return nil
	}

	maxAttempts := len(subjects) * len(hours)

	// Function to generate a single timetable
	generate := func() FacultyBasedTimetable {
		timetable := make(FacultyBasedTimetable)
		subjectsAssigned := make(map[string]map[string]bool)
		periodsLeft := make(map[string]int)
		facultyDailyCount := make(map[string]map[string]int)
		status0Assignments := make(map[string]map[string]bool)
		labSubjectAssigned := make(map[string]bool)
		facultyAssignments := make(map[string]map[string]int)
		var labSubjects, nonLabSubjects []models.Subject
		for _, subject := range subjects {
			periodsLeft[subject.Name] = subject.Period
			if subject.Status == 0 {
				labSubjects = append(labSubjects, subject)
				status0Assignments[subject.Name] = make(map[string]bool)
			} else {
				nonLabSubjects = append(nonLabSubjects, subject)
			}
		}
		// Iterate over days and apply timetable skips and manual timetable entries.
		for _, day := range days {
			timetable[day] = make(map[string][]models.TimetableEntry)
			subjectsAssigned[day] = make(map[string]bool)
			facultyDailyCount[day] = make(map[string]int)

			// Apply timetable skips
			if skips, ok := skipTimetable[day]; ok {
				for startTime, entries := range skips {
					for _, entry := range entries {
						timetable[day][startTime] = append(timetable[day][startTime], entry)
						subjectsAssigned[day][entry.SubjectName] = true
					}
				}
			}

			// Apply manual timetable entries
			if manualEntries, ok := manualTimetable[day]; ok {
				for startTime, entries := range manualEntries {
					for _, entry := range entries {
						// Ensure you're not overwriting existing timetable entries
						if len(timetable[day][startTime]) == 0 {
							timetable[day][startTime] = append(timetable[day][startTime], entry)
							subjectsAssigned[day][entry.SubjectName] = true
						}
					}
				}
			}
		}

		rand.Seed(time.Now().UnixNano())

		for _, day := range days {
			for i := 0; i < len(hours); i++ {
				for attempts := 0; attempts < maxAttempts; attempts++ {
					var filteredLabSubjects []models.Subject

					// Allocate a lab venue
					for _, subject := range labSubjects {
						if periodsLeft[subject.Name] > 0 &&
							(!subjectsAssigned[day][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) &&
							!labSubjectAssigned[day] {
							filteredLabSubjects = append(filteredLabSubjects, subject)
						}
					}

					if len(filteredLabSubjects) == 0 {
						break
					}

					subjectIndex := rand.Intn(len(filteredLabSubjects))
					subject := filteredLabSubjects[subjectIndex]

					hourIndex := i
					startTime := hours[hourIndex].StartTime
					endTime := hours[hourIndex].EndTime

					if len(timetable[day][startTime]) > 0 {
						continue
					}

					if subject.Status == 0 && i < len(hours)-1 {
						nextStartTime := hours[i+1].StartTime
						nextEndTime := hours[i+1].EndTime

						if IsPeriodAvailable(timetable, day, nextStartTime, "") {
							facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)

							if facultyName == "" {
								fmt.Println("Error: No faculty available for lab subject", subject.Name)
								return nil
							}
							if facultyDailyCount[day][facultyName] >= 2 {
								continue // Skip to the next attempt if assigned twice
							}
							var labVenue models.LabVenue
							for _, venue := range labVenues {
								if venue.SubjectID == subject.ID {
									labVenue = venue
									break
								}
							}

							if labVenue.ID == 0 {
								continue // No lab venue found for this subject
							}
							log.Println("LAB VENUE:", labVenue)

							entry1 := models.TimetableEntry{
								DayName:      day,
								StartTime:    startTime,
								EndTime:      endTime,
								SubjectName:  subject.Name,
								FacultyName:  facultyName,
								LabVenue:     labVenue.LabVenue,
								Status:       subject.Status,
								SemesterID:   semesterID,
								DepartmentID: departmentID,
								AcademicYear: academicYearID,
								CourseCode:   subject.CourseCode,
								SectionID:    sectionID,
							}

							entry2 := models.TimetableEntry{
								DayName:      day,
								StartTime:    nextStartTime,
								EndTime:      nextEndTime,
								SubjectName:  subject.Name,
								FacultyName:  entry1.FacultyName,
								LabVenue:     entry1.LabVenue,
								Status:       subject.Status,
								SemesterID:   entry1.SemesterID,
								DepartmentID: departmentID,
								AcademicYear: academicYearID,
								CourseCode:   subject.CourseCode,
								SectionID:    sectionID,
							}

							timetable[day][startTime] = append(timetable[day][startTime], entry1)
							timetable[day][nextStartTime] = append(timetable[day][nextStartTime], entry2)

							periodsLeft[subject.Name] -= 2
							subjectsAssigned[day][subject.Name] = true
							status0Assignments[subject.Name][startTime] = true
							labSubjectAssigned[day] = true
							// Ensure facultyAssignments[day] is initialized
							if _, exists := facultyAssignments[day]; !exists {
								facultyAssignments[day] = make(map[string]int)
							}

							// Safely increment the count for facultyName
							facultyAssignments[day][facultyName]++

							facultyDailyCount[day][facultyName] += 2
							break
						}
					}
				}
			}
		}

		// Handle non-lab subjects
		for _, day := range days {
			for i := 0; i < len(hours); i++ {
				for attempts := 0; attempts < maxAttempts; attempts++ {
					var filteredNonLabSubjects []models.Subject
					for _, subject := range nonLabSubjects {
						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day][subject.Name] {
							filteredNonLabSubjects = append(filteredNonLabSubjects, subject)
						}
					}

					if len(filteredNonLabSubjects) == 0 {
						break
					}

					subjectIndex := rand.Intn(len(filteredNonLabSubjects))
					subject := filteredNonLabSubjects[subjectIndex]

					hourIndex := i
					startTime := hours[hourIndex].StartTime
					endTime := hours[hourIndex].EndTime

					if len(timetable[day][startTime]) > 0 {
						continue
					}

					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)

					if facultyName == "" {
						fmt.Println("Error: No faculty available for non-lab subject", subject.Name)
						return nil
					}
					if facultyDailyCount[day][facultyName] >= 1 {
						continue // Skip to the next attempt if assigned twice
					}
					classroomName := selectRandomClassroom(classrooms)
					if classroomName == "" {
						fmt.Println("Error: No classroom found for non-lab subject", subject.Name)
						return nil
					}

					entry := models.TimetableEntry{
						DayName:      day,
						StartTime:    startTime,
						EndTime:      endTime,
						SubjectName:  subject.Name,
						FacultyName:  facultyName,
						Classroom:    classroomName,
						Status:       subject.Status,
						SemesterID:   semesterID,
						DepartmentID: departmentID,
						AcademicYear: academicYearID,
						CourseCode:   subject.CourseCode,
						SectionID:    sectionID,
					}

					timetable[day][startTime] = append(timetable[day][startTime], entry)
					periodsLeft[subject.Name]--
					subjectsAssigned[day][subject.Name] = true
					// Ensure facultyAssignments[day] is initialized
					if _, exists := facultyAssignments[day]; !exists {
						facultyAssignments[day] = make(map[string]int)
					}

					// Safely increment the count for facultyName
					facultyAssignments[day][facultyName]++

					facultyDailyCount[day][facultyName]++
					break
				}
			}
		}

		return timetable
	}

	// Keep generating until all periods are filled
	for {
		timetable := generate()
		allPeriodsFilled := true
		for _, day := range days {
			for _, hour := range hours {
				if len(timetable[day][hour.StartTime]) == 0 {
					allPeriodsFilled = false
					break
				}
			}
			if !allPeriodsFilled {
				break
			}
		}
		if allPeriodsFilled {
			return timetable
		}
	}
}

// func generateRandomTimetable(
// 	workingDays []models.WorkingDay,
// 	hours []models.Hour,
// 	subjects []models.Subject,
// 	faculty []models.Faculty,
// 	classrooms []models.Classroom,
// 	facultySubjects []models.FacultySubject,
// 	section []models.Section,
// 	semesters []models.Semester,
// 	departmentID int,
// 	semesterID int,
// 	academicYearID int,
// 	sectionID int,
// ) FacultyBasedTimetable {
// 	log.Println("Academic Year ID:", academicYearID)

// 	var days []string
// 	for _, wd := range workingDays {
// 		days = append(days, wd.WorkingDate.Format("2006-01-02"))
// 	}

// 	labVenues, err := GetLabVenue()
// 	if err != nil {
// 		fmt.Println("Error fetching lab venues:", err)
// 		return nil
// 	}

// 	skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		fmt.Println("Error fetching timetable skips:", err)
// 		return nil
// 	}

// 	manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		fmt.Println("Error fetching manual timetable:", err)
// 		return nil
// 	}

// 	maxAttempts := len(subjects) * len(hours)

// 	generate := func() FacultyBasedTimetable {
// 		timetable := make(FacultyBasedTimetable)
// 		subjectDailyCount := make(map[string]map[string]int)
// 		periodsLeft := make(map[string]int)
// 		labSubjectAssigned := make(map[string]bool)

// 		var labSubjects, nonLabSubjects []models.Subject
// 		for _, subject := range subjects {
// 			periodsLeft[subject.Name] = subject.Period
// 			if subject.Status == 0 {
// 				labSubjects = append(labSubjects, subject)
// 			} else {
// 				nonLabSubjects = append(nonLabSubjects, subject)
// 			}
// 		}

// 		for _, day := range days {
// 			timetable[day] = make(map[string][]models.TimetableEntry)
// 			subjectDailyCount[day] = make(map[string]int)

// 			if skips, ok := skipTimetable[day]; ok {
// 				for startTime, entries := range skips {
// 					for _, entry := range entries {
// 						timetable[day][startTime] = append(timetable[day][startTime], entry)
// 					}
// 				}
// 			}

// 			if manualEntries, ok := manualTimetable[day]; ok {
// 				for startTime, entries := range manualEntries {
// 					for _, entry := range entries {
// 						if len(timetable[day][startTime]) == 0 {
// 							timetable[day][startTime] = append(timetable[day][startTime], entry)
// 						}
// 					}
// 				}
// 			}
// 		}

// 		rand.Seed(time.Now().UnixNano())

// 		// Assign lab subjects
// 		for _, day := range days {
// 			for i := 0; i < len(hours)-1; i++ {
// 				for attempts := 0; attempts < maxAttempts; attempts++ {
// 					var availableLabSubjects []models.Subject
// 					for _, subject := range labSubjects {
// 						if periodsLeft[subject.Name] > 0 && subjectDailyCount[day][subject.Name] < 2 && !labSubjectAssigned[day] {
// 							availableLabSubjects = append(availableLabSubjects, subject)
// 						}
// 					}

// 					if len(availableLabSubjects) == 0 {
// 						break
// 					}

// 					subject := availableLabSubjects[rand.Intn(len(availableLabSubjects))]

// 					startTime := hours[i].StartTime
// 					endTime := hours[i].EndTime
// 					nextStartTime := hours[i+1].StartTime
// 					nextEndTime := hours[i+1].EndTime

// 					if len(timetable[day][startTime]) > 0 || len(timetable[day][nextStartTime]) > 0 {
// 						continue
// 					}

// 					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
// 					if facultyName == "" {
// 						continue
// 					}

// 					var labVenue models.LabVenue
// 					for _, venue := range labVenues {
// 						if venue.SubjectID == subject.ID {
// 							labVenue = venue
// 							break
// 						}
// 					}

// 					if labVenue.ID == 0 {
// 						continue
// 					}

// 					entry1 := models.TimetableEntry{
// 						DayName:      day,
// 						StartTime:    startTime,
// 						EndTime:      endTime,
// 						SubjectName:  subject.Name,
// 						FacultyName:  facultyName,
// 						LabVenue:     labVenue.LabVenue,
// 						Status:       subject.Status,
// 						SemesterID:   semesterID,
// 						DepartmentID: departmentID,
// 						AcademicYear: academicYearID,
// 						CourseCode:   subject.CourseCode,
// 						SectionID:    sectionID,
// 					}

// 					entry2 := entry1
// 					entry2.StartTime = nextStartTime
// 					entry2.EndTime = nextEndTime

// 					timetable[day][startTime] = append(timetable[day][startTime], entry1)
// 					timetable[day][nextStartTime] = append(timetable[day][nextStartTime], entry2)

// 					periodsLeft[subject.Name] -= 2
// 					subjectDailyCount[day][subject.Name] += 2
// 					labSubjectAssigned[day] = true
// 					break
// 				}
// 			}
// 		}

// 		// Assign non-lab subjects
// 		for _, day := range days {
// 			for i := 0; i < len(hours); i++ {
// 				for attempts := 0; attempts < maxAttempts; attempts++ {
// 					var availableNonLabSubjects []models.Subject
// 					for _, subject := range nonLabSubjects {
// 						if periodsLeft[subject.Name] > 0 && subjectDailyCount[day][subject.Name] < 1 {
// 							availableNonLabSubjects = append(availableNonLabSubjects, subject)
// 						}
// 					}

// 					if len(availableNonLabSubjects) == 0 {
// 						break
// 					}

// 					subject := availableNonLabSubjects[rand.Intn(len(availableNonLabSubjects))]

// 					startTime := hours[i].StartTime
// 					endTime := hours[i].EndTime

// 					if len(timetable[day][startTime]) > 0 {
// 						continue
// 					}

// 					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
// 					if facultyName == "" {
// 						continue
// 					}

// 					classroomName := selectRandomClassroom(classrooms)
// 					if classroomName == "" {
// 						continue
// 					}

// 					entry := models.TimetableEntry{
// 						DayName:      day,
// 						StartTime:    startTime,
// 						EndTime:      endTime,
// 						SubjectName:  subject.Name,
// 						FacultyName:  facultyName,
// 						Classroom:    classroomName,
// 						Status:       subject.Status,
// 						SemesterID:   semesterID,
// 						DepartmentID: departmentID,
// 						AcademicYear: academicYearID,
// 						CourseCode:   subject.CourseCode,
// 						SectionID:    sectionID,
// 					}

// 					timetable[day][startTime] = append(timetable[day][startTime], entry)
// 					periodsLeft[subject.Name]--
// 					subjectDailyCount[day][subject.Name]++
// 					break
// 				}
// 			}
// 		}

// 		return timetable
// 	}

// 	for {
// 		timetable := generate()
// 		allPeriodsFilled := true
// 		for _, day := range days {
// 			for _, hour := range hours {
// 				if len(timetable[day][hour.StartTime]) == 0 {
// 					allPeriodsFilled = false
// 					break
// 				}
// 			}
// 			if !allPeriodsFilled {
// 				break
// 			}
// 		}
// 		if allPeriodsFilled {
// 			return timetable
// 		}
// 	}
// }

func selectRandomClassroom(classrooms []models.Classroom) string {
	if len(classrooms) > 0 {
		return classrooms[rand.Intn(len(classrooms))].ClassroomName
	}
	return ""
}

func selectRandomFaculty(facultyList []models.Faculty, subject models.Subject, facultySubjects []models.FacultySubject, departmentID int, semesterID int, academicYearID int, sectionID int) string {
	var availableFaculty []models.Faculty
	for _, fac := range facultyList {
		for _, fs := range facultySubjects {
			if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID &&
				fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
				fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
				availableFaculty = append(availableFaculty, fac)
				break
			}
		}
	}
	if len(availableFaculty) > 0 {
		return availableFaculty[rand.Intn(len(availableFaculty))].FacultyName
	}
	return ""
}
