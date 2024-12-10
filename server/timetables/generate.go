// package timetables

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"ssg-portal/config"
// 	"ssg-portal/models"
// 	"time"
// )

// type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

// func FetchExistingTimetable() (map[string]map[string][]models.TimetableEntry, error) {
// 	existingTimetable := make(map[string]map[string][]models.TimetableEntry)

// 	rows, err := config.Database.Query(`
//         SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom, status, semester_id, department_id, academic_year, course_code ,section_id
//         FROM timetable`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var dayName, startTime, endTime, subjectName, facultyName, classroom string
// 		var courseCode []byte
// 		var status, semesterID, departmentID, academicYearID, sectionID int

// 		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom, &status, &semesterID, &departmentID, &academicYearID, &courseCode, &sectionID); err != nil {
// 			return nil, err
// 		}

// 		courseCodeStr := string(courseCode)

// 		if _, exists := existingTimetable[facultyName]; !exists {
// 			existingTimetable[facultyName] = make(map[string][]models.TimetableEntry)
// 		}

// 		entry := models.TimetableEntry{
// 			DayName:      dayName,
// 			StartTime:    startTime,
// 			EndTime:      endTime,
// 			SubjectName:  subjectName,
// 			FacultyName:  facultyName,
// 			Classroom:    classroom,
// 			Status:       status,
// 			SemesterID:   semesterID,
// 			DepartmentID: departmentID,
// 			AcademicYear: academicYearID,
// 			CourseCode:   courseCodeStr,
// 			SectionID:    sectionID,
// 		}

// 		existingTimetable[facultyName][dayName] = append(existingTimetable[facultyName][dayName], entry)
// 	}

// 	return existingTimetable, nil
// }
// func FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID int) (map[string]map[string][]models.TimetableEntry, error) {
// 	skipEntries := make(map[string]map[string][]models.TimetableEntry)

// 	query := `
// 	SELECT day_name, start_time, end_time, subject_name, faculty_name, semester_id, department_id, classroom, status, academic_year, course_code, section_id
// 	FROM timetable_skips
// 	WHERE department_id = ? AND semester_id = ? AND academic_year = ? AND section_id = ?`

// 	rows, err := config.Database.Query(query, departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var dayName, startTime, endTime, subjectName, facultyName, classroom, courseCode string
// 		var semesterID, departmentID, status, academicYear, sectionID int

// 		// Scan values from the row
// 		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &semesterID, &departmentID, &classroom, &status, &academicYear, &courseCode, &sectionID); err != nil {
// 			return nil, err
// 		}

// 		entry := models.TimetableEntry{
// 			DayName:      dayName,
// 			StartTime:    startTime,
// 			EndTime:      endTime,
// 			SubjectName:  subjectName,
// 			FacultyName:  facultyName,
// 			Classroom:    classroom,
// 			Status:       status,
// 			SemesterID:   semesterID,
// 			DepartmentID: departmentID,
// 			AcademicYear: academicYear,
// 			CourseCode:   courseCode,
// 			SectionID:    sectionID,
// 		}
// 		if skipEntries[dayName] == nil {
// 			skipEntries[dayName] = make(map[string][]models.TimetableEntry)
// 		}
// 		// Append the entry to the corresponding day and time
// 		skipEntries[dayName][startTime] = append(skipEntries[dayName][startTime], entry)
// 	}

// 	// Check for errors after iteration
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return skipEntries, nil
// }

// func FetchManualTimetable(departmentID int, semesterID int, academicYearID int, sectionID int) (map[string]map[string][]models.TimetableEntry, error) {

// 	manualTimetable := make(map[string]map[string][]models.TimetableEntry)

// 	query := `
//         SELECT day_name, start_time, end_time, classroom, semester_id, department_id,
//                subject_name, faculty_name, status, academic_year, course_code, section_id
//         FROM manual_timetable
//         WHERE department_id = ? AND semester_id = ? AND academic_year = ? AND section_id = ?`

// 	rows, err := config.Database.Query(query, departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch manual timetable: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var entry models.TimetableEntry
// 		if err := rows.Scan(&entry.DayName, &entry.StartTime, &entry.EndTime, &entry.Classroom,
// 			&entry.SemesterID, &entry.DepartmentID, &entry.SubjectName, &entry.FacultyName,
// 			&entry.Status, &entry.AcademicYear, &entry.CourseCode, &entry.SectionID); err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %w", err)
// 		}

// 		// Initialize the nested map for the day if it doesn't exist
// 		if _, exists := manualTimetable[entry.DayName]; !exists {
// 			manualTimetable[entry.DayName] = make(map[string][]models.TimetableEntry)
// 		}

// 		// Append the entry to the correct day and start time
// 		manualTimetable[entry.DayName][entry.StartTime] = append(manualTimetable[entry.DayName][entry.StartTime], entry)
// 	}

// 	// Check for any errors encountered during iteration
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error during rows iteration: %w", err)
// 	}

// 	return manualTimetable, nil
// }

// func IsPeriodAvailable(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
// 	if _, ok := existingTimetable[facultyName]; ok {
// 		if entries, ok := existingTimetable[facultyName][dayName]; ok {
// 			for _, entry := range entries {
// 				if entry.StartTime == startTime {
// 					return false
// 				}
// 			}
// 		}
// 	}
// 	return true
// }
// func Available(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
// 	for _, entries := range existingTimetable[facultyName][dayName] {
// 		if entries.StartTime == startTime {
// 			return false
// 		}
// 	}
// 	return true
// }

// func CheckTimetableConflicts(generatedTimetable FacultyBasedTimetable, existingTimetable map[string]map[string][]models.TimetableEntry) bool {
// 	for facultyName, days := range generatedTimetable {
// 		for dayName, entries := range days {
// 			if existingEntries, ok := existingTimetable[facultyName][dayName]; ok {
// 				for _, entry := range entries {
// 					for _, existingEntry := range existingEntries {
// 						// Check for conflicts in StartTime, Classroom, and Subject Name
// 						if entry.StartTime == existingEntry.StartTime &&
// 							entry.EndTime == existingEntry.EndTime &&
// 							entry.Classroom == existingEntry.Classroom &&
// 							entry.SubjectName == existingEntry.SubjectName {
// 							return true // Conflict found
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return false // No conflicts found
// }

// func GenerateTimetable(
// 	days []models.Day,
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

// 	// Fetch the existing timetable
// 	existingTimetable, err := FetchExistingTimetable()
// 	if err != nil {
// 		fmt.Println("Error fetching existing timetable:", err)
// 		return nil
// 	}

// 	// If no existing timetable is found, generate a random timetable
// 	if len(existingTimetable) == 0 {
// 		return generateRandomTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, section, semesters, departmentID, semesterID, academicYearID, sectionID)
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

// 	// Pre-compute mappings to reduce repeated lookups
// 	subjectClassrooms := map[string][]models.Classroom{}
// 	for _, subject := range subjects {
// 		for _, cls := range classrooms {
// 			if cls.DepartmentID == subject.DepartmentID {
// 				subjectClassrooms[subject.Name] = append(subjectClassrooms[subject.Name], cls)
// 			}
// 		}
// 	}

// 	// Map faculty subjects to reduce repeated filtering
// 	facultySubjectMap := map[int]map[int]bool{}
// 	for _, fs := range facultySubjects {
// 		if facultySubjectMap[fs.FacultyID] == nil {
// 			facultySubjectMap[fs.FacultyID] = map[int]bool{}
// 		}
// 		facultySubjectMap[fs.FacultyID][fs.SubjectID] = true
// 	}

// 	// Map for tracking lab subject assignments (status == 0)
// 	status0Assignments := make(map[string]map[string]bool) // Day -> {time -> true}
// 	labAssigned := make(map[string]bool)                   // Day -> is lab assigned

// 	for {
// 		timetable := make(map[string]map[string][]models.TimetableEntry)
// 		subjectsAssigned := make(map[string]map[string]bool)
// 		periodsLeft := make(map[string]int)

// 		facultyDailyCount := make(map[string]map[string]int)

// 		// Initialize periods left for each subject
// 		for _, subject := range subjects {
// 			periodsLeft[subject.Name] = subject.Period
// 		}

// 		// Initialize facultyAssignments for each day
// 		facultyAssignments := make(map[string]map[string]string)

// 		// Incorporate manual timetable into the generated timetable
// 		for _, day := range days {
// 			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
// 			subjectsAssigned[day.DayName] = make(map[string]bool)

// 			facultyDailyCount[day.DayName] = make(map[string]int)
// 			labAssigned[day.DayName] = false

// 			// Initialize facultyAssignments[day.DayName]
// 			facultyAssignments[day.DayName] = make(map[string]string)

// 			// If skips exist for the day, add them
// 			if skips, ok := skipTimetable[day.DayName]; ok {
// 				for startTime, entries := range skips { // `entries` is a slice of `models.TimetableEntry`
// 					for _, entry := range entries { // Iterate over individual `models.TimetableEntry` structs
// 						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
// 						subjectsAssigned[day.DayName][entry.SubjectName] = true
// 					}
// 				}
// 			}

// 			// Add manual timetable entries for the current day
// 			if manualEntries, ok := manualTimetable[day.DayName]; ok {
// 				for startTime, entries := range manualEntries {
// 					for _, entry := range entries {
// 						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
// 						subjectsAssigned[day.DayName][entry.SubjectName] = true
// 						facultyDailyCount[day.DayName][entry.FacultyName]++
// 						if entry.Status == 0 {
// 							labAssigned[day.DayName] = true
// 						}
// 					}
// 				}
// 			}
// 		}

// 		// Automatic generation for remaining periods...
// 		rand.Seed(time.Now().UnixNano())

// 		for _, day := range days {
// 			for _, hour := range hours {
// 				startTime := hour.StartTime

// 				if len(timetable[day.DayName][startTime]) > 0 {
// 					continue
// 				}

// 				assigned := false
// 				for attempts := 0; attempts < 1000; attempts++ {

// 					var filteredSubjects []models.Subject
// 					for _, subject := range subjects {
// 						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
// 							if subject.Status == 0 && labAssigned[day.DayName] {
// 								continue
// 							}
// 							filteredSubjects = append(filteredSubjects, subject)
// 						}
// 					}

// 					if len(filteredSubjects) == 0 {
// 						break
// 					}

// 					subject := filteredSubjects[rand.Intn(len(filteredSubjects))]

// 					var availableFaculty []models.Faculty
// 					for _, fac := range faculty {
// 						if facultySubjectMap[fac.ID][subject.ID] && facultyDailyCount[day.DayName][fac.FacultyName] < 2 {
// 							availableFaculty = append(availableFaculty, fac)
// 						}
// 					}

// 					if len(availableFaculty) == 0 {
// 						continue
// 					}

// 					selectedFaculty := availableFaculty[rand.Intn(len(availableFaculty))]

// 					// Check if the faculty has already been assigned to this start time
// 					if assignedClassroom, exists := facultyAssignments[day.DayName][selectedFaculty.FacultyName]; exists && assignedClassroom == startTime {
// 						continue // Faculty is already assigned to this time slot
// 					}

// 					// Select a classroom
// 					selectedClassroom := subjectClassrooms[subject.Name][0]
// 					if !Available(existingTimetable, day.DayName, startTime, selectedFaculty.FacultyName) {
// 						continue
// 					}

// 					// Create the timetable entry
// 					entry := models.TimetableEntry{
// 						DayName:      day.DayName,
// 						StartTime:    startTime,
// 						EndTime:      hour.EndTime,
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

// 					// Assign and track period
// 					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
// 					periodsLeft[subject.Name]--
// 					facultyDailyCount[day.DayName][selectedFaculty.FacultyName]++
// 					subjectsAssigned[day.DayName][subject.Name] = true
// 					facultyAssignments[day.DayName][selectedFaculty.FacultyName] = startTime

// 					// Track lab assignment if the subject is a lab (status 0)
// 					if subject.Status == 0 {
// 						status0Assignments[day.DayName][startTime] = true
// 						labAssigned[day.DayName] = true
// 					}

// 					assigned = true
// 					break
// 				}

// 				if !assigned {
// 					fmt.Printf("Warning: Could not assign a subject for %s at %s\n", day.DayName, startTime)
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

// 		if allAssigned && !CheckTimetableConflicts(timetable, existingTimetable) {
// 			return timetable
// 		}
// 	}
// }
// func generateRandomTimetable(
// 	days []models.Day,
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

// 	// Fetch timetable skips
// 	skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		fmt.Println("Error fetching timetable skips:", err)
// 		return nil
// 	}

// 	// Fetch manual timetable
// 	manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
// 	if err != nil {
// 		fmt.Println("Error fetching manual timetable:", err)
// 		return nil
// 	}

// 	maxAttempts := len(subjects) * len(hours)

// 	// Function to generate a single timetable
// 	generate := func() FacultyBasedTimetable {
// 		timetable := make(FacultyBasedTimetable)
// 		subjectsAssigned := make(map[string]map[string]bool)
// 		periodsLeft := make(map[string]int)
// 		labSubjectAssigned := make(map[string]bool)
// 		facultyDailyCount := make(map[string]map[string]int)

// 		var labSubjects, nonLabSubjects []models.Subject
// 		for _, subject := range subjects {
// 			periodsLeft[subject.Name] = subject.Period
// 			if subject.Status == 0 {
// 				labSubjects = append(labSubjects, subject)
// 			} else {
// 				nonLabSubjects = append(nonLabSubjects, subject)
// 			}
// 		}

// 		// Pre-fill timetable with manual timetable data
// 		for _, day := range days {
// 			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
// 			subjectsAssigned[day.DayName] = make(map[string]bool)
// 			labSubjectAssigned[day.DayName] = false
// 			facultyDailyCount[day.DayName] = make(map[string]int)

// 			if skips, ok := skipTimetable[day.DayName]; ok {
// 				for startTime, entries := range skips { // `entries` is a slice of `models.TimetableEntry`
// 					for _, entry := range entries { // Iterate over individual `models.TimetableEntry` structs
// 						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
// 						subjectsAssigned[day.DayName][entry.SubjectName] = true
// 					}
// 				}
// 			}
// 			// Add manual timetable entries
// 			if manualEntries, ok := manualTimetable[day.DayName]; ok {
// 				for startTime, entries := range manualEntries {
// 					for _, entry := range entries {
// 						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
// 						subjectsAssigned[day.DayName][entry.SubjectName] = true
// 					}
// 				}
// 			}
// 		}

// 		rand.Seed(time.Now().UnixNano())

// 		// Generate timetable for lab subjects first
// 		for _, day := range days {
// 			for i := 0; i < len(hours); i++ {
// 				startTime := hours[i].StartTime
// 				if len(timetable[day.DayName][startTime]) > 0 {
// 					continue
// 				}

// 				for attempts := 0; attempts < maxAttempts; attempts++ {
// 					// Filter eligible lab subjects
// 					var filteredLabSubjects []models.Subject
// 					for _, subject := range labSubjects {
// 						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] && !labSubjectAssigned[day.DayName] {
// 							filteredLabSubjects = append(filteredLabSubjects, subject)
// 						}
// 					}

// 					if len(filteredLabSubjects) == 0 {
// 						break
// 					}

// 					subject := filteredLabSubjects[rand.Intn(len(filteredLabSubjects))]

// 					// Assign lab subjects in two consecutive periods
// 					if i < len(hours)-1 {
// 						nextStartTime := hours[i+1].StartTime

// 						// Randomly select faculty and classroom
// 						facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
// 						if facultyName == "" || facultyDailyCount[day.DayName][facultyName] >= 2 {
// 							continue
// 						}

// 						classroomName := selectRandomClassroom(classrooms)
// 						if classroomName == "" {
// 							continue
// 						}

// 						// Create timetable entries
// 						entry1 := models.TimetableEntry{
// 							DayName:      day.DayName,
// 							StartTime:    startTime,
// 							EndTime:      hours[i].EndTime,
// 							SubjectName:  subject.Name,
// 							FacultyName:  facultyName,
// 							Classroom:    classroomName,
// 							Status:       subject.Status,
// 							SemesterID:   semesterID,
// 							DepartmentID: departmentID,
// 							AcademicYear: academicYearID,
// 							CourseCode:   subject.CourseCode,
// 							SectionID:    sectionID,
// 						}
// 						entry2 := entry1
// 						entry2.StartTime = nextStartTime
// 						entry2.EndTime = hours[i+1].EndTime

// 						// Add entries
// 						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry1)
// 						timetable[day.DayName][nextStartTime] = append(timetable[day.DayName][nextStartTime], entry2)

// 						// Update assignment tracking
// 						periodsLeft[subject.Name] -= 2
// 						subjectsAssigned[day.DayName][subject.Name] = true
// 						labSubjectAssigned[day.DayName] = true
// 						facultyDailyCount[day.DayName][facultyName] += 2
// 						break
// 					}
// 				}
// 			}
// 		}

// 		// Generate timetable for non-lab subjects
// 		for _, day := range days {
// 			for i := 0; i < len(hours); i++ {
// 				startTime := hours[i].StartTime
// 				if len(timetable[day.DayName][startTime]) > 0 {
// 					continue
// 				}

// 				for attempts := 0; attempts < maxAttempts; attempts++ {
// 					// Filter eligible non-lab subjects
// 					var filteredNonLabSubjects []models.Subject
// 					for _, subject := range nonLabSubjects {
// 						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
// 							filteredNonLabSubjects = append(filteredNonLabSubjects, subject)
// 						}
// 					}

// 					if len(filteredNonLabSubjects) == 0 {
// 						break
// 					}

// 					subject := filteredNonLabSubjects[rand.Intn(len(filteredNonLabSubjects))]

// 					// Randomly select faculty and classroom
// 					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
// 					if facultyName == "" || facultyDailyCount[day.DayName][facultyName] >= 1 {
// 						continue
// 					}

// 					classroomName := selectRandomClassroom(classrooms)
// 					if classroomName == "" {
// 						continue
// 					}

// 					// Add timetable entry
// 					entry := models.TimetableEntry{
// 						DayName:      day.DayName,
// 						StartTime:    startTime,
// 						EndTime:      hours[i].EndTime,
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

// 					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
// 					periodsLeft[subject.Name]--
// 					subjectsAssigned[day.DayName][subject.Name] = true
// 					facultyDailyCount[day.DayName][facultyName]++
// 					break
// 				}
// 			}
// 		}

// 		return timetable
// 	}

// 	// Keep generating until all periods are filled
// 	for {
// 		timetable := generate()
// 		allPeriodsFilled := true
// 		for _, day := range days {
// 			for _, hour := range hours {
// 				if len(timetable[day.DayName][hour.StartTime]) == 0 {
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

//	func selectRandomClassroom(classrooms []models.Classroom) string {
//		if len(classrooms) > 0 {
//			return classrooms[rand.Intn(len(classrooms))].ClassroomName
//		}
//		return ""
//	}
//
//	func selectRandomFaculty(facultyList []models.Faculty, subject models.Subject, facultySubjects []models.FacultySubject, departmentID int, semesterID int, academicYearID int, sectionID int) string {
//		var availableFaculty []models.Faculty
//		for _, fac := range facultyList {
//			for _, fs := range facultySubjects {
//				// Check if the faculty ID matches, and also check the additional criteria
//				if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID &&
//					fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
//					fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
//					availableFaculty = append(availableFaculty, fac)
//					break
//				}
//			}
//		}
//		if len(availableFaculty) > 0 {
//			return availableFaculty[rand.Intn(len(availableFaculty))].FacultyName
//		}
//		return ""
//	}
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

		// Append the entry to the correct day and start time
		manualTimetable[entry.DayName][entry.StartTime] = append(manualTimetable[entry.DayName][entry.StartTime], entry)
	}

	// Check for any errors encountered during iteration
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
func Available(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
	for _, entries := range existingTimetable[facultyName][dayName] {
		if entries.StartTime == startTime {
			return false
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
						// Check for conflicts in StartTime, Classroom, and Subject Name
						if entry.StartTime == existingEntry.StartTime &&
							entry.EndTime == existingEntry.EndTime &&
							entry.Classroom == existingEntry.Classroom &&
							entry.SubjectName == existingEntry.SubjectName {
							return true // Conflict found
						}
					}
				}
			}
		}
	}
	return false // No conflicts found
}
func GenerateTimetable(
	days []models.Day,
	hours []models.Hour,
	subjects []models.Subject,
	faculty []models.Faculty,
	classrooms []models.Classroom,
	facultySubjects []models.FacultySubject,
	semesters []models.Semester,
	section []models.Section,
	academicYear []models.AcademicYear,
	departmentID int,
	semesterID int,
	academicYearID int,
	sectionID int) map[string]map[string][]models.TimetableEntry {

	// Fetch the existing timetable
	existingTimetable, err := FetchExistingTimetable()
	if err != nil {
		fmt.Println("Error fetching existing timetable:", err)
		return nil
	}

	// If no existing timetable is found, generate a random timetable
	if len(existingTimetable) == 0 {
		return generateRandomTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, section, semesters, departmentID, semesterID, academicYearID, sectionID)
	}

	// Fetch skips (blocked periods for the timetable)
	skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		fmt.Println("Error fetching timetable skips:", err)
		return nil
	}

	// Fetch manual timetable entries
	manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		fmt.Println("Error fetching manual timetable:", err)
		return nil
	}

	// Pre-compute mappings to reduce repeated lookups
	subjectClassrooms := map[string][]models.Classroom{}
	for _, subject := range subjects {
		for _, cls := range classrooms {
			if cls.DepartmentID == subject.DepartmentID {
				subjectClassrooms[subject.Name] = append(subjectClassrooms[subject.Name], cls)
			}
		}
	}

	// Map faculty subjects to reduce repeated filtering
	facultySubjectMap := map[int]map[int]bool{}
	for _, fs := range facultySubjects {
		if facultySubjectMap[fs.FacultyID] == nil {
			facultySubjectMap[fs.FacultyID] = map[int]bool{}
		}
		facultySubjectMap[fs.FacultyID][fs.SubjectID] = true
	}

	// Map for tracking lab subject assignments (status == 0)
	status0Assignments := make(map[string]map[string]bool) // Day -> {time -> true}
	labAssigned := make(map[string]bool)                   // Day -> is lab assigned

	for {
		timetable := make(map[string]map[string][]models.TimetableEntry)
		subjectsAssigned := make(map[string]map[string]bool)
		periodsLeft := make(map[string]int)

		facultyDailyCount := make(map[string]map[string]int)

		// Initialize periods left for each subject
		for _, subject := range subjects {
			periodsLeft[subject.Name] = subject.Period
		}

		// Initialize facultyAssignments for each day
		facultyAssignments := make(map[string]map[string]string)

		// Incorporate manual timetable into the generated timetable
		for _, day := range days {
			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
			subjectsAssigned[day.DayName] = make(map[string]bool)

			facultyDailyCount[day.DayName] = make(map[string]int)
			labAssigned[day.DayName] = false

			// Initialize facultyAssignments[day.DayName]
			facultyAssignments[day.DayName] = make(map[string]string)

			// If skips exist for the day, add them
			if skips, ok := skipTimetable[day.DayName]; ok {
				for startTime, entries := range skips { // `entries` is a slice of `models.TimetableEntry`
					for _, entry := range entries { // Iterate over individual `models.TimetableEntry` structs
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						subjectsAssigned[day.DayName][entry.SubjectName] = true
					}
				}
			}

			// Add manual timetable entries for the current day
			if manualEntries, ok := manualTimetable[day.DayName]; ok {
				for startTime, entries := range manualEntries {
					for _, entry := range entries {
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						subjectsAssigned[day.DayName][entry.SubjectName] = true
						facultyDailyCount[day.DayName][entry.FacultyName]++
						if entry.Status == 0 {
							labAssigned[day.DayName] = true
						}
					}
				}
			}
		}

		// Automatic generation for remaining periods...
		rand.Seed(time.Now().UnixNano())

		for _, day := range days {
			for _, hour := range hours {
				startTime := hour.StartTime

				if len(timetable[day.DayName][startTime]) > 0 {
					continue
				}

				assigned := false
				for attempts := 0; attempts < 1000; attempts++ {

					var filteredSubjects []models.Subject
					for _, subject := range subjects {
						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
							if subject.Status == 0 && labAssigned[day.DayName] {
								continue
							}
							filteredSubjects = append(filteredSubjects, subject)
						}
					}

					if len(filteredSubjects) == 0 {
						break
					}

					subject := filteredSubjects[rand.Intn(len(filteredSubjects))]

					var availableFaculty []models.Faculty
					for _, fac := range faculty {
						if facultySubjectMap[fac.ID][subject.ID] && facultyDailyCount[day.DayName][fac.FacultyName] < 2 {
							availableFaculty = append(availableFaculty, fac)
						}
					}

					if len(availableFaculty) == 0 {
						continue
					}

					selectedFaculty := availableFaculty[rand.Intn(len(availableFaculty))]

					// Check if the faculty has already been assigned to this start time
					if assignedClassroom, exists := facultyAssignments[day.DayName][selectedFaculty.FacultyName]; exists && assignedClassroom == startTime {
						continue // Faculty is already assigned to this time slot
					}

					// Select a classroom
					selectedClassroom := subjectClassrooms[subject.Name][0]
					if !Available(existingTimetable, day.DayName, startTime, selectedFaculty.FacultyName) {
						continue
					}

					// Create the timetable entry
					entry := models.TimetableEntry{
						DayName:      day.DayName,
						StartTime:    startTime,
						EndTime:      hour.EndTime,
						SubjectName:  subject.Name,
						FacultyName:  selectedFaculty.FacultyName,
						Classroom:    selectedClassroom.ClassroomName,
						Status:       subject.Status,
						SemesterID:   selectedClassroom.SemesterID,
						DepartmentID: departmentID,
						AcademicYear: academicYearID,
						CourseCode:   subject.CourseCode,
						SectionID:    sectionID,
					}

					// Assign and track period
					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
					periodsLeft[subject.Name]--
					facultyDailyCount[day.DayName][selectedFaculty.FacultyName]++
					subjectsAssigned[day.DayName][subject.Name] = true
					facultyAssignments[day.DayName][selectedFaculty.FacultyName] = startTime

					// Track lab assignment if the subject is a lab (status 0)
					if subject.Status == 0 {
						status0Assignments[day.DayName][startTime] = true
						labAssigned[day.DayName] = true
					}

					assigned = true
					break
				}

				if !assigned {
					fmt.Printf("Warning: Could not assign a subject for %s at %s\n", day.DayName, startTime)
				}
			}
		}

		// Ensure all periods are filled
		allAssigned := true
		for subjectName, remainingPeriods := range periodsLeft {
			if remainingPeriods > 0 {
				fmt.Printf("Warning: Subject %s has %d periods left unassigned.\n", subjectName, remainingPeriods)
				allAssigned = false
			}
		}

		if allAssigned && !CheckTimetableConflicts(timetable, existingTimetable) {
			return timetable
		}
	}
}
func generateRandomTimetable(
	days []models.Day,
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
		labSubjectAssigned := make(map[string]bool)
		facultyDailyCount := make(map[string]map[string]int)

		var labSubjects, nonLabSubjects []models.Subject
		for _, subject := range subjects {
			periodsLeft[subject.Name] = subject.Period
			if subject.Status == 0 {
				labSubjects = append(labSubjects, subject)
			} else {
				nonLabSubjects = append(nonLabSubjects, subject)
			}
		}

		// Pre-fill timetable with manual timetable data
		for _, day := range days {
			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
			subjectsAssigned[day.DayName] = make(map[string]bool)
			labSubjectAssigned[day.DayName] = false
			facultyDailyCount[day.DayName] = make(map[string]int)

			if skips, ok := skipTimetable[day.DayName]; ok {
				for startTime, entries := range skips { // `entries` is a slice of `models.TimetableEntry`
					for _, entry := range entries { // Iterate over individual `models.TimetableEntry` structs
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						subjectsAssigned[day.DayName][entry.SubjectName] = true
					}
				}
			}
			// Add manual timetable entries
			if manualEntries, ok := manualTimetable[day.DayName]; ok {
				for startTime, entries := range manualEntries {
					for _, entry := range entries {
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						subjectsAssigned[day.DayName][entry.SubjectName] = true
					}
				}
			}
		}

		rand.Seed(time.Now().UnixNano())

		// Generate timetable for lab subjects first
		for _, day := range days {
			for i := 0; i < len(hours); i++ {
				startTime := hours[i].StartTime
				if len(timetable[day.DayName][startTime]) > 0 {
					continue
				}

				for attempts := 0; attempts < maxAttempts; attempts++ {
					// Filter eligible lab subjects
					var filteredLabSubjects []models.Subject
					for _, subject := range labSubjects {
						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] && !labSubjectAssigned[day.DayName] {
							filteredLabSubjects = append(filteredLabSubjects, subject)
						}
					}

					if len(filteredLabSubjects) == 0 {
						break
					}

					subject := filteredLabSubjects[rand.Intn(len(filteredLabSubjects))]

					// Assign lab subjects in two consecutive periods
					if i < len(hours)-1 {
						nextStartTime := hours[i+1].StartTime

						// Randomly select faculty and classroom
						facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
						if facultyName == "" || facultyDailyCount[day.DayName][facultyName] >= 2 {
							continue
						}

						classroomName := selectRandomClassroom(classrooms)
						if classroomName == "" {
							continue
						}

						// Create timetable entries
						entry1 := models.TimetableEntry{
							DayName:      day.DayName,
							StartTime:    startTime,
							EndTime:      hours[i].EndTime,
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
						entry2 := entry1
						entry2.StartTime = nextStartTime
						entry2.EndTime = hours[i+1].EndTime

						// Add entries
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry1)
						timetable[day.DayName][nextStartTime] = append(timetable[day.DayName][nextStartTime], entry2)

						// Update assignment tracking
						periodsLeft[subject.Name] -= 2
						subjectsAssigned[day.DayName][subject.Name] = true
						labSubjectAssigned[day.DayName] = true
						facultyDailyCount[day.DayName][facultyName] += 2
						break
					}
				}
			}
		}

		// Generate timetable for non-lab subjects
		for _, day := range days {
			for i := 0; i < len(hours); i++ {
				startTime := hours[i].StartTime
				if len(timetable[day.DayName][startTime]) > 0 {
					continue
				}

				for attempts := 0; attempts < maxAttempts; attempts++ {
					// Filter eligible non-lab subjects
					var filteredNonLabSubjects []models.Subject
					for _, subject := range nonLabSubjects {
						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
							filteredNonLabSubjects = append(filteredNonLabSubjects, subject)
						}
					}

					if len(filteredNonLabSubjects) == 0 {
						break
					}

					subject := filteredNonLabSubjects[rand.Intn(len(filteredNonLabSubjects))]

					// Randomly select faculty and classroom
					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
					if facultyName == "" || facultyDailyCount[day.DayName][facultyName] >= 1 {
						continue
					}

					classroomName := selectRandomClassroom(classrooms)
					if classroomName == "" {
						continue
					}

					// Add timetable entry
					entry := models.TimetableEntry{
						DayName:      day.DayName,
						StartTime:    startTime,
						EndTime:      hours[i].EndTime,
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

					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
					periodsLeft[subject.Name]--
					subjectsAssigned[day.DayName][subject.Name] = true
					facultyDailyCount[day.DayName][facultyName]++
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
				if len(timetable[day.DayName][hour.StartTime]) == 0 {
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
			// Check if the faculty ID matches, and also check the additional criteria
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
