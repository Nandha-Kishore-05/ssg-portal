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

func IsPeriodAvailable(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
	if _, ok := existingTimetable[facultyName]; ok {
		if entries, ok := existingTimetable[facultyName][dayName]; ok {
			for _, entry := range entries {
				if entry.StartTime == startTime {
					return false
				}
			}
		}
	}
	return true
}


func CheckPeriodAvailability(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
	// Check if the faculty already has a period scheduled for the given day and time
	if entries, exists := existingTimetable[facultyName][dayName]; exists {
		for _, entry := range entries {
			if entry.StartTime == startTime {
				return false
			}
		}
	}
	return true
}


func CheckTimetableConflicts(generatedTimetable FacultyBasedTimetable, existingTimetable map[string]map[string][]models.TimetableEntry) bool {
	for facultyName, days := range generatedTimetable {
		for dayName, entries := range days {
			if existingEntries, ok := existingTimetable[facultyName][dayName]; ok {
				for _, entry := range entries {
					for _, existingEntry := range existingEntries {
						if entry.StartTime == existingEntry.StartTime &&
							entry.EndTime == existingEntry.EndTime &&
							entry.Classroom == existingEntry.Classroom &&
							entry.SubjectName == existingEntry.SubjectName {
							return true
						}
					}
				}
			}
		}
	}
	return false
}





// func GenerateTimetable(
// 	workingDays []models.WorkingDay,
// 	hours []models.Hour,
// 	subjects []models.Subject,
// 	faculty []models.Faculty,
// 	classrooms []models.Classroom,
// 	facultySubjects []models.FacultySubject,
// 	semesters []models.Semester,
// 	section []models.Section,
// 	academicYear []models.AcademicYear,
// 	departmentID, semesterID, academicYearID, sectionID int) map[string]map[string][]models.TimetableEntry {

// 	// Initialize data structures
// 	var days []string
// 	for _, wd := range workingDays {
// 		days = append(days, wd.WorkingDate.Format("2006-01-02"))
// 	}

// 	// Fetch the existing timetable and handle errors
// 	existingTimetable, err := FetchExistingTimetable()
// 	if err != nil {
// 		fmt.Println("Error fetching existing timetable:", err)
// 		return nil
// 	}
// 	if len(existingTimetable) == 0 {
// 		return generateRandomTimetable(workingDays, hours, subjects, faculty, classrooms, facultySubjects, section, semesters, departmentID, semesterID, academicYearID, sectionID)
// 	}

// 	// Fetch lab venues, skips, and manual timetable entries
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

// 	// manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
// 	// if err != nil {
// 	// 	fmt.Println("Error fetching manual timetable:", err)
// 	// 	return nil
// 	// }

// 	// Map subject to classrooms and faculty to subjects
// 	subjectClassrooms := map[string][]models.Classroom{}
// 	for _, subject := range subjects {
// 		for _, cls := range classrooms {
// 			if cls.DepartmentID == subject.DepartmentID {
// 				subjectClassrooms[subject.Name] = append(subjectClassrooms[subject.Name], cls)
// 			}
// 		}
// 	}

// 	facultySubjectMap := map[int]map[int]bool{}
// 	for _, fs := range facultySubjects {
// 		if facultySubjectMap[fs.FacultyID] == nil {
// 			facultySubjectMap[fs.FacultyID] = map[int]bool{}
// 		}
// 		facultySubjectMap[fs.FacultyID][fs.SubjectID] = true
// 	}

// 	// Generate timetable using a randomized approach with max attempts


// 	for {
// 		timetable := make(map[string]map[string][]models.TimetableEntry)
// 		subjectsAssigned := make(map[string]map[string]bool)
// 		periodsLeft := make(map[string]int)
// 		status0Assignments := make(map[string]map[string]bool)
// 		facultyAssignments := make(map[string]map[string]string)
// 		facultyDailyCount := make(map[string]map[string]int) 
// 		labAssigned := make(map[string]bool)

// 		for _, subject := range subjects {
// 			periodsLeft[subject.Name] = subject.Period
// 			if subject.Status == 0 {
// 				status0Assignments[subject.Name] = make(map[string]bool)
// 			}
// 		}

// 		for _, day := range days {
// 			timetable[day] = make(map[string][]models.TimetableEntry)
// 			subjectsAssigned[day] = make(map[string]bool)
// 			facultyAssignments[day] = make(map[string]string)
// 			facultyDailyCount[day] = make(map[string]int) 
// 			labAssigned[day] = false
// 			if skips, ok := skipTimetable[day]; ok {
// 				for startTime, _ := range skips {
// 					timetable[day][startTime] = append(timetable[day][startTime])
// 					// subjectsAssigned[day][entry] = true
// 					// periodsLeft[entry.SubjectName]--
// 					// if entry.Status == 0 { // Assuming Status 1 indicates a lab subject
// 					// 	labAssigned[day] = true // Mark lab as assigned
// 					// }
// 				}
// 			}
// 		}

// 		rand.Seed(time.Now().UnixNano())

// 		for _, day := range days {
// 			for i := 0; i < len(hours); i++ {
// 				assigned := false
// 				for attempts := 0; attempts < 1000; attempts++ {
// 					var filteredSubjects []models.Subject
// 					for _, subject := range subjects {
// 						if periodsLeft[subject.Name] > 0 && (!subjectsAssigned[day][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) {
	
// 							if subject.Status == 0 && labAssigned[day] {
// 								continue
// 							}
// 							if subject.Status == 1 && subjectsAssigned[day][subject.Name] {
// 								continue
// 							}
// 							var validClassrooms []models.Classroom
// 							for _, cls := range classrooms {
// 								if cls.DepartmentID == subject.DepartmentID {
// 									for _, semester := range semesters {
// 										if semester.ID == cls.SemesterID {
// 											validClassrooms = append(validClassrooms, cls)
// 											break
// 										}
// 									}
// 								}
// 							}

// 							if len(validClassrooms) > 0 {
// 								filteredSubjects = append(filteredSubjects, subject)
// 							}
// 						}
// 					}

// 					if len(filteredSubjects) == 0 {
// 						continue
// 					}

// 					subjectIndex := rand.Intn(len(filteredSubjects))
// 					subject := filteredSubjects[subjectIndex]

// 					hourIndex := i % len(hours)
// 					startTime := hours[hourIndex].StartTime
// 					endTime := hours[hourIndex].EndTime

// 					if _, ok := timetable[day][startTime]; ok {
// 						if len(timetable[day][startTime]) > 0 {
// 							continue
// 						}
// 					}

	
// 					// }
// 					var availableFaculty []models.Faculty
// 					for _, fac := range faculty {
// 						for _, fs := range facultySubjects {

// 							if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID &&
// 								fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
// 								fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
// 								availableFaculty = append(availableFaculty, fac)
// 								break
// 							}
// 						}
// 					}
	
		

// 					if len(availableFaculty) == 0 {
// 						continue
// 					}

// 					facultyIndex := rand.Intn(len(availableFaculty))
// 					selectedFaculty := availableFaculty[facultyIndex]

// 					if facultyDailyCount[day][selectedFaculty.FacultyName] >= 2 {
// 						continue
// 					}

// 					if assignedClassroom, exists := facultyAssignments[day][selectedFaculty.FacultyName]; exists && assignedClassroom == startTime {
// 						continue
// 					}

// 					var selectedClassroom models.Classroom
// 					for _, cls := range classrooms {
// 						if cls.DepartmentID == subject.DepartmentID {
// 							for _, semester := range semesters {
// 								if semester.ID == cls.SemesterID {
// 									selectedClassroom = cls
// 									break
// 								}
// 							}
// 							break
// 						}
// 					}
				
// 					if !Available(existingTimetable, day, startTime, availableFaculty[facultyIndex].FacultyName) {
// 						continue
// 					}
// 					entry := models.TimetableEntry{
// 						DayName:      day,
// 						StartTime:    startTime,
// 						EndTime:      endTime,
// 						SubjectName:  subject.Name,
// 						FacultyName:  selectedFaculty.FacultyName,
// 						Classroom:    selectedClassroom.ClassroomName,
// 						Status:       subject.Status,
// 						SemesterID:   selectedClassroom.SemesterID,
// 						DepartmentID: departmentID,
// 						AcademicYear: academicYearID,
// 						CourseCode:   subject.CourseCode,
// 						SectionID:    sectionID,
// 					}

// 					if _, ok := timetable[day]; !ok {
// 						timetable[day] = make(map[string][]models.TimetableEntry)
// 					}
// 					if _, ok := timetable[day][startTime]; !ok {
// 						timetable[day][startTime] = []models.TimetableEntry{}
// 					}

// 					if subject.Status == 0 {
// 						if _, ok := status0Assignments[subject.Name][startTime]; !ok {
// 							if i < len(hours)-1 {
// 								nextHourIndex := (hourIndex + 1) % len(hours)
// 								nextStartTime := hours[nextHourIndex].StartTime
// 								if IsPeriodAvailable(existingTimetable, day, nextStartTime, "") {
// 									entry2 := entry
// 									entry2.StartTime = nextStartTime
// 									entry2.EndTime = hours[nextHourIndex].EndTime

// 									timetable[day][startTime] = append(timetable[day][startTime], entry)
// 									timetable[day][nextStartTime] = append(timetable[day][nextStartTime], entry2)
// 									periodsLeft[subject.Name] -= 2
// 									subjectsAssigned[day][subject.Name] = true
// 									status0Assignments[subject.Name][startTime] = true
// 									status0Assignments[subject.Name][nextStartTime] = true
// 									facultyAssignments[day][selectedFaculty.FacultyName] = nextStartTime
// 									facultyDailyCount[day][selectedFaculty.FacultyName]+= 2
// 									labAssigned[day] = true
// 									assigned = true
// 									break
// 								}
// 							}
// 							continue
// 						}
// 					}

				
// 					timetable[day][startTime] = append(timetable[day][startTime], entry)
// 					periodsLeft[subject.Name]--
// 					facultyDailyCount[day][selectedFaculty.FacultyName]++
// 					subjectsAssigned[day][subject.Name] = true
// 					if subject.Status == 0 {
// 						status0Assignments[subject.Name][startTime] = true
// 					}
// 					facultyAssignments[day][selectedFaculty.FacultyName] = startTime
// 					facultyDailyCount[day][selectedFaculty.FacultyName]++ 
// 					if subject.Status == 0 { // Mark as lab subject assigned
// 						labAssigned[day] = true
// 					}
// 					assigned = true
// 					break
// 				}
// 				if !assigned {
// 					fmt.Printf("Warning: Could not assign a subject for %s during period %d\n", day, i+1)
// 				}
// 			}
// 		}

// 		allAssigned := true
// 		for subjectName, remainingPeriods := range periodsLeft {
// 			if remainingPeriods > 0 {
// 				fmt.Printf("Warning: Subject %s has %d periods left unassigned.\n", subjectName, remainingPeriods)
// 				allAssigned = false
// 			}
// 		}

// 		// Check if all periods are filled
// 		periodsFilled := true
// 		for _, day := range days {
// 			for _, hour := range hours {
// 				startTime := hour.StartTime
// 				if len(timetable[day][startTime]) == 0 {
// 					fmt.Printf("Warning: Subject %s has %d periods left unassigned.\n")
// 					periodsFilled = false
// 					break
// 				}
// 			}
// 			if !periodsFilled {
// 				break
// 			}
// 		}

// 		// Regenerate if not all periods are filled or if there are conflicts
// 		if allAssigned && periodsFilled && !TimetableConflicts(timetable, existingTimetable) {
// 			return timetable
// 		}
// 	}

// }


func GenerateTimetable(
	workingDays []models.WorkingDay,
	hours []models.Hour,
	subjects []models.Subject,
	faculty []models.Faculty,
	classrooms []models.Classroom,
	facultySubjects []models.FacultySubject,
	semesters []models.Semester,
	section []models.Section,
	academicYear []models.AcademicYear,
	departmentID, semesterID, academicYearID, sectionID int) map[string]map[string][]models.TimetableEntry {

	// Initialize data structures
	var days []string
	for _, wd := range workingDays {
		days = append(days, wd.WorkingDate.Format("2006-01-02"))
	}

	// Fetch the existing timetable and handle errors
	existingTimetable, err := FetchExistingTimetable()
	if err != nil {
		fmt.Println("Error fetching existing timetable:", err)
		return nil
	}
	if len(existingTimetable) == 0 {
		return generateRandomTimetable(workingDays, hours, subjects, faculty, classrooms, facultySubjects, section, semesters, departmentID, semesterID, academicYearID, sectionID)
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

	// manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
	// if err != nil {
	// 	fmt.Println("Error fetching manual timetable:", err)
	// 	return nil
	// }

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

		// Initialize a map to track assigned subjects per section
	

	// Generate timetable using a randomized approach with max attempts
	for {
		timetable := make(map[string]map[string][]models.TimetableEntry)
		subjectsAssigned := make(map[string]map[string]bool)
		periodsLeft := make(map[string]int)
		status0Assignments := make(map[string]map[string]bool)
		facultyAssignments := make(map[string]map[string]string)
		facultyDailyCount := make(map[string]map[string]int)
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
			facultyDailyCount[day] = make(map[string]int)
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
						if periodsLeft[subject.Name] > 0 && (!subjectsAssigned[day][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) {

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

					if facultyDailyCount[day][selectedFaculty.FacultyName] >= 2 {
						continue
					}

					if assignedClassroom, exists := facultyAssignments[day][selectedFaculty.FacultyName]; exists && assignedClassroom == startTime {
						continue
					}
					var selectedClassroom models.Classroom
					if subject.Status == 0 && len(labVenues) > 0 { // Check lab venues for lab subjects
						// Use the LabVenue for lab subjects
						selectedLabVenue := labVenues[rand.Intn(len(labVenues))]
						selectedClassroom = models.Classroom{
							ID:           selectedLabVenue.ID,
							ClassroomName: selectedLabVenue.LabVenue, // Map LabVenue fields to Classroom-like structure
							DepartmentID: selectedLabVenue.DepartmentID,
							SemesterID:   selectedLabVenue.SemesterID,
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
					

					if !Available(existingTimetable, day, startTime, availableFaculty[facultyIndex].FacultyName) {
						continue
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
									facultyDailyCount[day][selectedFaculty.FacultyName] += 2
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
					facultyDailyCount[day][selectedFaculty.FacultyName]++
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
					fmt.Printf("Warning: Could not assign a subject for %s during period %d\n", day, i+1)
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

		// Regenerate if not all periods are filled or if there are conflicts
		if allAssigned && periodsFilled && !TimetableConflicts(timetable, existingTimetable) {
			return timetable
		}
	}
}



func Available(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
for _, entries := range existingTimetable[facultyName][dayName] {
	if entries.StartTime == startTime {
		return false
	}
}
return true
}

func TimetableConflicts(generatedTimetable FacultyBasedTimetable, existingTimetable map[string]map[string][]models.TimetableEntry) bool {
for facultyName, days := range generatedTimetable {
	for dayName, entries := range days {
		if existingEntries, ok := existingTimetable[facultyName][dayName]; ok {
			for _, entry := range entries {
				for _, existingEntry := range existingEntries {
					if entry.StartTime == existingEntry.StartTime && entry.Classroom == existingEntry.Classroom {
						return true
					}
				}
			}
		}
	}
}
return false
}



// func Available(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
// 	for _, entries := range existingTimetable[facultyName][dayName] {
// 		if entries.StartTime == startTime {
// 			return false
// 		}
// 	}
// 	return true
// }
// func CheckTimetableConflicts(generatedTimetable map[string]map[string][]models.TimetableEntry, existingTimetable map[string]map[string][]models.TimetableEntry) bool {
//     for day, periods := range generatedTimetable {
//         for startTime, entries := range periods {
//             // Check conflicts within the generated timetable
//             if len(entries) > 1 {
//                 return true // Conflict: Multiple entries for the same period
//             }

//             for _, entry := range entries {
//                 // Check conflicts with the existing timetable
//                 if existingEntries, ok := existingTimetable[day][startTime]; ok {
//                     for _, existingEntry := range existingEntries {
//                         if entry.FacultyName == existingEntry.FacultyName ||
//                             entry.Classroom == existingEntry.Classroom {
//                             return true // Conflict: Faculty or classroom overlap
//                         }
//                     }
//                 }
//             }
//         }
//     }
//     return false
// }


// func GenerateTimetable(
// 	workingDays []models.WorkingDay,
// 	hours []models.Hour,
// 	subjects []models.Subject,
// 	faculty []models.Faculty,
// 	classrooms []models.Classroom,
// 	facultySubjects []models.FacultySubject,
// 	semesters []models.Semester,
// 	section []models.Section,
// 	academicYear []models.AcademicYear,
// 	departmentID int,
// 	semesterID int,
// 	academicYearID int,
// 	sectionID int) map[string]map[string][]models.TimetableEntry {

// 	var days []string
// 	for _, wd := range workingDays {
// 		days = append(days, wd.WorkingDate.Format("2006-01-02")) // Use date in "YYYY-MM-DD" format
// 	} // Fetch the existing timetable

// 	existingTimetable, err := FetchExistingTimetable()
// 	if err != nil {
// 		fmt.Println("Error fetching existing timetable:", err)
// 		return nil
// 	}
// 	// If no existing timetable is found, generate a random timetable
// 	if len(existingTimetable) == 0 {
// 		return generateRandomTimetable(workingDays, hours, subjects, faculty, classrooms, facultySubjects, section, semesters, departmentID, semesterID, academicYearID, sectionID)
// 	}

// 	labVenues, err := GetLabVenue()
// 	if err != nil {
// 		fmt.Println("Error fetching lab venues:", err)
// 		return nil
// 	}
// 	// Fetch skips (blocked periods for the timetable)
// 	skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		fmt.Println("Error fetching timetable skips:", err)
// 		return nil
// 	}
// 	// Fetch manual timetable entries
// 	manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		fmt.Println("Error fetching manual timetable:", err)
// 		return nil
// 	}
// 	maxAttempts := len(subjects) * len(hours)

// 	subjectClassrooms := map[string][]models.Classroom{}
// 	for _, subject := range subjects {
// 		for _, cls := range classrooms {
// 			if cls.DepartmentID == subject.DepartmentID {
// 				subjectClassrooms[subject.Name] = append(subjectClassrooms[subject.Name], cls)
// 			}
// 		}
// 	}

// 	facultySubjectMap := map[int]map[int]bool{}
// 	for _, fs := range facultySubjects {
// 		if facultySubjectMap[fs.FacultyID] == nil {
// 			facultySubjectMap[fs.FacultyID] = map[int]bool{}
// 		}
// 		facultySubjectMap[fs.FacultyID][fs.SubjectID] = true
// 	}
// 	// Day -> is lab assigned
// 	for {
// 		timetable := make(map[string]map[string][]models.TimetableEntry)
// 		subjectsAssigned := make(map[string]map[string]bool)
// 		periodsLeft := make(map[string]int)
// 		facultyDailyCount := make(map[string]map[string]int)
// 		labSubjectAssigned := make(map[string]bool)
// 		status0Assignments := make(map[string]map[string]bool) // Day -> {time -> true}
// 		labAssigned := make(map[string]bool)
// 		facultyAssignments := make(map[string]map[string]int)
// 		// Initialize periods left for each subject
// 		var labSubjects, nonLabSubjects []models.Subject
// 		for _, subject := range subjects {
// 			periodsLeft[subject.Name] = subject.Period
// 			if subject.Status == 0 {
// 				labSubjects = append(labSubjects, subject)
// 				status0Assignments[subject.Name] = make(map[string]bool)
// 			} else {
// 				nonLabSubjects = append(nonLabSubjects, subject)
// 			}
// 		}
// 		// Initialize facultyAssignments for each day

// 		// Incorporate manual timetable into the generated timetable
// 		for _, day := range days {
// 			timetable[day] = make(map[string][]models.TimetableEntry)
// 			subjectsAssigned[day] = make(map[string]bool)
// 			facultyDailyCount[day] = make(map[string]int)

// 			facultyAssignments[day] = make(map[string]int)
// 			// If skips exist for the day, add them
// 			if skips, ok := skipTimetable[day]; ok {
// 				for startTime, entries := range skips { // `entries` is a slice of `models.TimetableEntry`
// 					for _, entry := range entries { // Iterate over individual `models.TimetableEntry` structs
// 						timetable[day][startTime] = append(timetable[day][startTime], entry)
// 						subjectsAssigned[day][entry.SubjectName] = true
// 					}
// 				}
// 			}
// 			// Add manual timetable entries for the current day
// 			if manualEntries, ok := manualTimetable[day]; ok {
// 				for startTime, entries := range manualEntries {
// 					for _, entry := range entries {
// 						timetable[day][startTime] = append(timetable[day][startTime], entry)
// 						subjectsAssigned[day][entry.SubjectName] = true
// 						facultyDailyCount[day][entry.FacultyName]++
// 						if entry.Status == 0 {
// 							labAssigned[day] = true
// 						}
// 					}
// 				}
// 			}
// 		}
// 		// Automatic generation for remaining periods...
// 		rand.Seed(time.Now().UnixNano())

// 		for _, day := range days {
// 			for i := 0; i < len(hours); i++ {
// 				for attempts := 0; attempts < maxAttempts; attempts++ {
// 					var filteredLabSubjects []models.Subject
// 					for _, subject := range labSubjects {
// 						if periodsLeft[subject.Name] > 0 &&
// 							(!subjectsAssigned[day][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) &&
// 							!labSubjectAssigned[day] {
// 							filteredLabSubjects = append(filteredLabSubjects, subject)
// 						}
// 					}

// 					if len(filteredLabSubjects) == 0 {
// 						break
// 					}

// 					subjectIndex := rand.Intn(len(filteredLabSubjects))
// 					subject := filteredLabSubjects[subjectIndex]

// 					hourIndex := i
// 					startTime := hours[hourIndex].StartTime
// 					endTime := hours[hourIndex].EndTime

// 					if len(timetable[day][startTime]) > 0 {
// 						continue
// 					}

// 					if subject.Status == 0 && i < len(hours)-1 {
// 						nextStartTime := hours[i+1].StartTime
// 						nextEndTime := hours[i+1].EndTime

// 						if IsPeriodAvailable(timetable, day, nextStartTime, "") {
// 							facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)

// 							if facultyName == "" {
// 								fmt.Println("Error: No faculty available for lab subject", subject.Name)
// 								return nil
// 							}
// 							if facultyDailyCount[day][facultyName] >= 2 {
// 								continue // Skip to the next attempt if assigned twice
// 							}
// 							var labVenue models.LabVenue
// 							for _, venue := range labVenues {
// 								if venue.SubjectID == subject.ID {
// 									labVenue = venue
// 									break
// 								}
// 							}

// 							if labVenue.ID == 0 {
// 								continue // No lab venue found for this subject
// 							}
// 							log.Println("LAB VENUE:", labVenue)

// 							entry1 := models.TimetableEntry{
// 								DayName:      day,
// 								StartTime:    startTime,
// 								EndTime:      endTime,
// 								SubjectName:  subject.Name,
// 								FacultyName:  facultyName,
// 								LabVenue:     labVenue.LabVenue,
// 								Status:       subject.Status,
// 								SemesterID:   semesterID,
// 								DepartmentID: departmentID,
// 								AcademicYear: academicYearID,
// 								CourseCode:   subject.CourseCode,
// 								SectionID:    sectionID,
// 							}

// 							entry2 := models.TimetableEntry{
// 								DayName:      day,
// 								StartTime:    nextStartTime,
// 								EndTime:      nextEndTime,
// 								SubjectName:  subject.Name,
// 								FacultyName:  entry1.FacultyName,
// 								LabVenue:     entry1.LabVenue,
// 								Status:       subject.Status,
// 								SemesterID:   entry1.SemesterID,
// 								DepartmentID: departmentID,
// 								AcademicYear: academicYearID,
// 								CourseCode:   subject.CourseCode,
// 								SectionID:    sectionID,
// 							}

// 							timetable[day][startTime] = append(timetable[day][startTime], entry1)
// 							timetable[day][nextStartTime] = append(timetable[day][nextStartTime], entry2)

// 							periodsLeft[subject.Name] -= 2
// 							subjectsAssigned[day][subject.Name] = true
// 							status0Assignments[subject.Name][startTime] = true
// 							labSubjectAssigned[day] = true
// 							// Ensure facultyAssignments[day] is initialized
// 							if _, exists := facultyAssignments[day]; !exists {
// 								facultyAssignments[day] = make(map[string]int)
// 							}

// 							// Safely increment the count for facultyName
// 							facultyAssignments[day][facultyName]++
// 							facultyDailyCount[day][facultyName] += 2
// 							break
// 						}
// 					}
// 				}
// 			}
// 		}

// 		for _, day := range days {
// 			for i := 0; i < len(hours); i++ {
// 				for attempts := 0; attempts < maxAttempts; attempts++ {
// 					var filteredNonLabSubjects []models.Subject
// 					for _, subject := range nonLabSubjects {
// 						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day][subject.Name] {
// 							filteredNonLabSubjects = append(filteredNonLabSubjects, subject)
// 						}
// 					}

// 					if len(filteredNonLabSubjects) == 0 {
// 						break
// 					}

// 					subjectIndex := rand.Intn(len(filteredNonLabSubjects))
// 					subject := filteredNonLabSubjects[subjectIndex]

// 					hourIndex := i
// 					startTime := hours[hourIndex].StartTime
// 					endTime := hours[hourIndex].EndTime

// 					if len(timetable[day][startTime]) > 0 {
// 						continue
// 					}
// 					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)

// 					if facultyName == "" {
// 						fmt.Println("Error: No faculty available for non-lab subject", subject.Name)
// 						return nil
// 					}
// 					if facultyDailyCount[day][facultyName] >= 1 {
// 						continue // Skip to the next attempt if assigned twice
// 					}
// 					classroomName := selectRandomClassroom(classrooms)
// 					if classroomName == "" {
// 						fmt.Println("Error: No classroom found for non-lab subject", subject.Name)
// 						return nil
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
// 					subjectsAssigned[day][subject.Name] = true
// 					facultyAssignments[day][facultyName]++
// 					facultyDailyCount[day][facultyName]++
// 					break
// 				}
// 			}
// 		}

// 		// Ensure all periods are filled
// 		allAssigned := true
// 		for subjectName, remainingPeriods := range periodsLeft {
// 			if remainingPeriods > 0 {
// 				fmt.Printf("Warning: Subject %s has %d periods left unassigned.\n", subjectName, remainingPeriods)
// 				allAssigned = false
// 			}
// 		}
// 		if allAssigned && CheckTimetableConflicts(timetable, existingTimetable) {
// 			return timetable
// 		}
// 	}
// }

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
