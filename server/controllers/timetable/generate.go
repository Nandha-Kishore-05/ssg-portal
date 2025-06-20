package timetable

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"ssg-portal/config"
	"ssg-portal/controllers/lab"
	"ssg-portal/models"
	"time"
)

type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

func PrecomputeFacultyAvailability(
	existingTimetable map[string]map[string][]models.TimetableEntry,
	facultyList []string) map[string]map[string]map[string]bool {

	// Structure: facultyName -> dayName -> startTime -> isUnavailable
	unavailabilityCache := make(map[string]map[string]map[string]bool)

	for _, facultyName := range facultyList {
		unavailabilityCache[facultyName] = make(map[string]map[string]bool)

		if facultyEntries, exists := existingTimetable[facultyName]; exists {
			for dayName, dayEntries := range facultyEntries {
				if _, ok := unavailabilityCache[facultyName][dayName]; !ok {
					unavailabilityCache[facultyName][dayName] = make(map[string]bool)
				}

				for _, entry := range dayEntries {
					unavailabilityCache[facultyName][dayName][entry.StartTime] = true
				}
			}
		}
	}

	return unavailabilityCache
}

func IsFacultyAvailableFast(
	unavailabilityCache map[string]map[string]map[string]bool,
	facultyList []string,
	dayName string,
	startTime string) []string {

	var unavailableFaculties []string

	for _, facultyName := range facultyList {
		if dayMap, exists := unavailabilityCache[facultyName]; exists {
			if timeMap, exists := dayMap[dayName]; exists {
				if timeMap[startTime] {
					unavailableFaculties = append(unavailableFaculties, facultyName)
				}
			}
		}
	}

	return unavailableFaculties
}

func AreAllSubjectFacultyAvailableAcrossSectionsFast(
	unavailabilityCache map[string]map[string]map[string]bool,
	facultySubjectMap map[int]map[int][]models.Faculty,
	subjectID int,
	dayName string,
	startTime string,
	sectionMap map[int]bool) bool {

	sections := make([]int, 0, len(sectionMap))
	for sectionID := range sectionMap {
		sections = append(sections, sectionID)
	}

	// Get all unique faculty names for this subject across all sections
	facultyNames := make(map[string]bool)
	for _, sectionID := range sections {
		if faculty, exists := facultySubjectMap[subjectID]; exists {
			if sectionFaculty, exists := faculty[sectionID]; exists {
				for _, f := range sectionFaculty {
					facultyNames[f.FacultyName] = true
				}
			}
		}
	}

	// Convert to slice for checking
	facultyList := make([]string, 0, len(facultyNames))
	for name := range facultyNames {
		facultyList = append(facultyList, name)
	}

	// Quick check using the precomputed cache
	unavailableFaculties := IsFacultyAvailableFast(unavailabilityCache, facultyList, dayName, startTime)
	return len(unavailableFaculties) == 0
}

func BuildSectionMap(departmentID, semesterID, academicYearID int) (map[int]bool, error) {
	sectionMap := make(map[int]bool)

	rows, err := config.Database.Query(`
        SELECT DISTINCT section_id 
        FROM faculty_subjects 
        WHERE department_id = ? AND semester_id = ? AND academic_year_id = ?`,
		departmentID, semesterID, academicYearID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sectionID int
		if err := rows.Scan(&sectionID); err != nil {
			return nil, err
		}
		sectionMap[sectionID] = true
	}

	return sectionMap, nil
}

// Function to build the faculty-subject map in memory (called during initialization)
func BuildFacultySubjectMap(departmentID, semesterID, academicYearID int) (map[int]map[int][]models.Faculty, error) {
	// Map structure: subject_id -> section_id -> []Faculty
	facultySubjectMap := make(map[int]map[int][]models.Faculty)

	rows, err := config.Database.Query(`
        SELECT fs.subject_id, fs.section_id, f.id, f.name 
        FROM faculty f
        JOIN faculty_subjects fs ON f.id = fs.faculty_id
        WHERE fs.department_id = ? 
        AND fs.semester_id = ? 
        AND fs.academic_year_id = ?`,
		departmentID, semesterID, academicYearID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subjectID, sectionID, facultyID int
		var facultyName string

		if err := rows.Scan(&subjectID, &sectionID, &facultyID, &facultyName); err != nil {
			return nil, err
		}

		// Initialize inner map if not exists
		if facultySubjectMap[subjectID] == nil {
			facultySubjectMap[subjectID] = make(map[int][]models.Faculty)
		}

		// Add faculty to the appropriate section
		faculty := models.Faculty{
			ID:          facultyID,
			FacultyName: facultyName,
		}

		facultySubjectMap[subjectID][sectionID] = append(facultySubjectMap[subjectID][sectionID], faculty)
	}

	return facultySubjectMap, nil
}

func IsLabVenueAvailable(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, labVenueName string) bool {
	// Check all faculty entries in existing timetable for the given day and time
	for _, facultySchedule := range existingTimetable {
		if daySchedule, exists := facultySchedule[dayName]; exists {
			for _, entry := range daySchedule {
				if entry.StartTime == startTime && entry.Classroom == labVenueName {
					return false // Lab venue is already occupied
				}
			}
		}
	}
	return true
}

func AreAllLabVenuesAvailableForAllSections(existingTimetable map[string]map[string][]models.TimetableEntry,
	labVenuesBySubject map[int][]models.LabVenue, subjectID int, dayName, startTime1, startTime2 string,
	sectionMap map[int]bool) bool {

	if venues, exists := labVenuesBySubject[subjectID]; exists {

		sections := make([]int, 0, len(sectionMap))
		for sectionID := range sectionMap {
			sections = append(sections, sectionID)
		}
		//log.Println("nikitha",sections)
		// Check if we have enough lab venues for all sections
		if len(venues) < len(sections) {
			return false // Not enough lab venues for all sections
		}

		// Try to find a valid assignment of venues to sections
		return canAssignVenuesToAllSections(existingTimetable, venues, dayName, startTime1, startTime2, sectionMap)
	}
	return false // No venues found for this subject
}

func canAssignVenuesToAllSections(existingTimetable map[string]map[string][]models.TimetableEntry,
	venues []models.LabVenue, dayName, startTime1, startTime2 string, sectionMap map[int]bool) bool {

	// For simplicity, we'll use a greedy approach
	// In practice, you might want to use a more sophisticated matching algorithm
	sections := make([]int, 0, len(sectionMap))
	for sectionID := range sectionMap {
		sections = append(sections, sectionID)
	}
	availableVenues := make([]models.LabVenue, 0)

	// Find all venues that are completely free for both time slots
	for _, venue := range venues {
		if IsLabVenueAvailable(existingTimetable, dayName, startTime1, venue.LabVenue) &&
			IsLabVenueAvailable(existingTimetable, dayName, startTime2, venue.LabVenue) {
			availableVenues = append(availableVenues, venue)
		}
	}

	result := len(availableVenues) >= len(sections)

	// Check if we have enough available venues for all sections
	return result
}

func GenerateTimetable(
	days []models.Day,
	//workingDays []models.WorkingDay,
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
	// var days []string
	// for _, wd := range workingDays {
	// 	days = append(days, wd.WorkingDate.Format("2006-01-02"))
	// }

	// Check if this is a multi-section department for the same year and semester
	var sectionsInSameSemester []models.Section
	for _, section := range sections {
		if section.DepartmentID == departmentID && section.SemesterID == semesterID && section.AcademicYear == academicYearID {
			sectionsInSameSemester = append(sectionsInSameSemester, section)
		}
	}
	for _, section := range sections {
		log.Printf("Section: ID=%d | DepartmentID=%d (%T), SemesterID=%d (%T), AcademicYear=%v (%T)",
			section.ID,
			section.DepartmentID, section.DepartmentID,
			section.SemesterID, section.SemesterID,
			section.AcademicYear, section.AcademicYear,
		)

		log.Printf("Expected: DepartmentID=%d (%T), SemesterID=%d (%T), AcademicYearID=%d (%T)",
			departmentID, departmentID,
			semesterID, semesterID,
			academicYearID, academicYearID,
		)
	}

	log.Println("sectionsInSameSemester", sectionsInSameSemester)
	existingTimetable, err := FetchExistingTimetable(academicYearID)
	if err != nil {
		fmt.Println("Error fetching existing timetable:", err)
		return nil
	}

	if len(existingTimetable) == 0 {
		return generateRandomTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, sections, semesters, departmentID, semesterID, academicYearID, sectionID)
	}

	// Fetch any one section ID
	var existingSectionID int
	sectionQuery := `
	SELECT section_id FROM timetable 
	WHERE department_id = ? AND semester_id = ? AND academic_year = ?
	LIMIT 1
`
	err = config.Database.QueryRow(sectionQuery, departmentID, semesterID, academicYearID).Scan(&existingSectionID)

	if err == sql.ErrNoRows {
		log.Println("No section found for given criteria. Setting section ID to 0.")
		existingSectionID = 0
	} else if err != nil {
		log.Println("Error fetching section:", err)
	}

	if existingSectionID == 0 {
		log.Println("kishore")

		// Pre-build all data structures (same as before)
		existingTimetable, err := FetchExistingTimetable(academicYearID)
		if err != nil {
			fmt.Println("Error fetching existing timetable:", err)
			return nil
		}

		skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
		if err != nil {
			fmt.Println("Error fetching timetable skips:", err)
			return nil
		}

		allFacultyNames := make([]string, 0, len(faculty))
		for _, f := range faculty {
			allFacultyNames = append(allFacultyNames, f.FacultyName)
		}

		facultyAvailabilityCache := PrecomputeFacultyAvailability(existingTimetable, allFacultyNames)

		sectionMap, err := BuildSectionMap(departmentID, semesterID, academicYearID)
		if err != nil {
			log.Println("Error building section map:", err)
			return nil
		}

		facultySubjectMap, err := BuildFacultySubjectMap(departmentID, semesterID, academicYearID)
		if err != nil {
			log.Println("Error building faculty-subject map:", err)
			return nil
		}

		labVenues, err := lab.GetLabVenue()
		if err != nil {
			fmt.Println("Error fetching lab venues:", err)
			return nil
		}

		classroomsByDept := make(map[int][]models.Classroom)
		for _, cls := range classrooms {
			classroomsByDept[cls.DepartmentID] = append(classroomsByDept[cls.DepartmentID], cls)
		}

		facultyForSubjectSection := make(map[int]map[int][]models.Faculty)
		for _, fs := range facultySubjects {
			if fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
				fs.AcademicYear == academicYearID && fs.SectionID == sectionID {

				if facultyForSubjectSection[fs.SubjectID] == nil {
					facultyForSubjectSection[fs.SubjectID] = make(map[int][]models.Faculty)
				}

				for _, f := range faculty {
					if f.ID == fs.FacultyID {
						facultyForSubjectSection[fs.SubjectID][fs.SectionID] = append(
							facultyForSubjectSection[fs.SubjectID][fs.SectionID], f)
						break
					}
				}
			}
		}

		// Define the exact allowed lab slot pairs
		allowedLabPairs := [][2]string{
			{"08:45:00", "09:35:00"}, // 1st & 2nd periods
			{"10:40:00", "11:30:00"}, // 3rd & 4th periods
			{"13:30:00", "14:20:00"}, // 5th & 6th periods
		}

		// Pre-find all valid consecutive slot pairs for labs
		labSlots := make(map[string][][2]string)
		for _, day := range days {
			for _, pair := range allowedLabPairs {
				slot1 := pair[0]
				slot2 := pair[1]

				// Check if both slots are free
				free1 := IsLabPeriodAvailable(existingTimetable, day.DayName, slot1, "")
				free2 := IsLabPeriodAvailable(existingTimetable, day.DayName, slot2, "")

				if free1 && free2 {
					labSlots[day.DayName] = append(labSlots[day.DayName], [2]string{slot1, slot2})
				}
			}
		}

		labVenuesBySubject := make(map[int][]models.LabVenue)
		for _, labVenue := range labVenues {
			labVenuesBySubject[labVenue.SubjectID] = append(labVenuesBySubject[labVenue.SubjectID], labVenue)
		}

		// Limited retry with fast lab assignment
		for {
			timetable := make(map[string]map[string][]models.TimetableEntry)
			subjectsAssigned := make(map[string]map[string]bool)
			periodsLeft := make(map[string]int)
			status0Assignments := make(map[string]map[string]bool)
			facultyAssignments := make(map[string]map[string]string)
			subjectDailyCount := make(map[string]map[string]int)
			labAssignedday := make(map[string]int)
			usedLabSlots := make(map[string]map[string]bool) // Track used lab slots

			// Initialize
			for _, subject := range subjects {
				periodsLeft[subject.Name] = subject.Period
				if subject.Status == 0 {
					status0Assignments[subject.Name] = make(map[string]bool)
				}
			}

			for _, day := range days {
				timetable[day.DayName] = make(map[string][]models.TimetableEntry)
				subjectsAssigned[day.DayName] = make(map[string]bool)
				facultyAssignments[day.DayName] = make(map[string]string)
				subjectDailyCount[day.DayName] = make(map[string]int)
				labAssignedday[day.DayName] = 0
				usedLabSlots[day.DayName] = make(map[string]bool)

				// Apply skips
				if skips, ok := skipTimetable[day.DayName]; ok {
					for startTime, entries := range skips {
						for _, entry := range entries {
							timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
							subjectsAssigned[day.DayName][entry.SubjectName] = true
						}
					}
				}
			}

			rand.Seed(time.Now().UnixNano())

			// FAST LAB ASSIGNMENT: Process all lab subjects first
			labSuccess := true
			for _, subject := range subjects {
				if subject.Status != 0 || periodsLeft[subject.Name] <= 0 {
					continue
				}

				sessionsNeeded := periodsLeft[subject.Name] / 2
				labTryCount := 0
				for session := 0; session < sessionsNeeded; session++ {
					placed := false

					// Try each day in random order
					dayOrder := rand.Perm(len(days))
					for _, dayIdx := range dayOrder {
						if labTryCount >= 10 { // 🚫 Stop trying after 10 failed attempts
							break
						}
						day := days[dayIdx]

						if labAssignedday[day.DayName] >= 2 {
							continue
						}

						// Try available lab slots for this day
						if slots, exists := labSlots[day.DayName]; exists {
							rand.Shuffle(len(slots), func(i, j int) { slots[i], slots[j] = slots[j], slots[i] })

							for _, slot := range slots {
								if labTryCount >= 10 { // 🚫 Stop trying after 10 failed attempts
									break
								}
								slot1, slot2 := slot[0], slot[1]

								// Check if slots are already used
								if usedLabSlots[day.DayName][slot1] || usedLabSlots[day.DayName][slot2] {
								
									continue
								}

								// Quick faculty check
								if !AreAllSubjectFacultyAvailableAcrossSectionsFast(
									facultyAvailabilityCache, facultySubjectMap,
									subject.ID, day.DayName, slot1, sectionMap) {
									continue
								}

								// NEW: Check lab venue availability for both slots
								if !AreAllLabVenuesAvailableForAllSections(existingTimetable, labVenuesBySubject,
									subject.ID, day.DayName, slot1, slot2, sectionMap) {
										labTryCount++
									continue // Skip this timing if not all sections can get lab venues
								}

								// Get faculty and venue
								var availableFaculty []models.Faculty
								if subjectFaculty, exists := facultyForSubjectSection[subject.ID]; exists {
									if sectionFaculty, exists := subjectFaculty[sectionID]; exists {
										availableFaculty = sectionFaculty
									}
								}

								if len(availableFaculty) == 0 {
									continue
								}

								selectedFaculty := availableFaculty[rand.Intn(len(availableFaculty))]

								var selectedClassroom models.Classroom
								if venues, exists := labVenuesBySubject[subject.ID]; exists && len(venues) > 0 {
									// Find the first available lab venue for both slots
									var availableVenue models.LabVenue
									found := false
									for _, venue := range venues {
										if IsLabVenueAvailable(existingTimetable, day.DayName, slot1, venue.LabVenue) &&
											IsLabVenueAvailable(existingTimetable, day.DayName, slot2, venue.LabVenue) {
											availableVenue = venue
											found = true
											break
										}
									}
									if !found {
										continue // No available venue found
									}

									selectedClassroom = models.Classroom{
										ID:            availableVenue.ID,
										ClassroomName: availableVenue.LabVenue,
										DepartmentID:  availableVenue.DepartmentID,
										SemesterID:    availableVenue.SemesterID,
									}
								}

								// Find correct hour indices for slot1 and slot2
								var slot1Index, slot2Index int = -1, -1
								for idx, hour := range hours {
									if hour.StartTime == slot1 {
										slot1Index = idx
									}
									if hour.StartTime == slot2 {
										slot2Index = idx
									}
								}

								// Create entries with correct EndTime mapping
								slots := []struct {
									startTime string
									endTime   string
								}{
									{slot1, hours[slot1Index].EndTime},
									{slot2, hours[slot2Index].EndTime},
								}

								for _, slotInfo := range slots {
									entry := models.TimetableEntry{
										DayName:      day.DayName,
										StartTime:    slotInfo.startTime,
										EndTime:      slotInfo.endTime,
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
									timetable[day.DayName][slotInfo.startTime] = append(timetable[day.DayName][slotInfo.startTime], entry)
									usedLabSlots[day.DayName][slotInfo.startTime] = true
								}

								periodsLeft[subject.Name] -= 2
								subjectsAssigned[day.DayName][subject.Name] = true
								status0Assignments[subject.Name][slot1] = true
								status0Assignments[subject.Name][slot2] = true
								facultyAssignments[day.DayName][selectedFaculty.FacultyName] = slot2
								subjectDailyCount[day.DayName][subject.Name] += 2
								labAssignedday[day.DayName]++
									labTryCount = 0
								placed = true
								break
							}
						}

						if placed {
							break
						}
					}

					if !placed {
						labSuccess = false
						break
					}
				}

				if !labSuccess {
					break
				}
			}

			if !labSuccess {
				continue // Retry
			}

			// REGULAR SUBJECT ASSIGNMENT (same logic as before but simplified)
			for _, day := range days {
				for i := 0; i < len(hours); i++ {
					startTime := hours[i].StartTime

					if _, exists := timetable[day.DayName][startTime]; exists && len(timetable[day.DayName][startTime]) > 0 {
						continue
					}

					// Find eligible regular subjects
					var eligibleSubjects []models.Subject
					for _, subject := range subjects {
						if subject.Status == 0 || periodsLeft[subject.Name] <= 0 ||
							subjectsAssigned[day.DayName][subject.Name] {
							continue
						}

						if !AreAllSubjectFacultyAvailableAcrossSectionsFast(
							facultyAvailabilityCache, facultySubjectMap,
							subject.ID, day.DayName, startTime, sectionMap) {
							continue
						}

						if len(classroomsByDept[subject.DepartmentID]) > 0 {
							eligibleSubjects = append(eligibleSubjects, subject)
						}
					}

					if len(eligibleSubjects) == 0 {
						continue
					}

					subject := eligibleSubjects[rand.Intn(len(eligibleSubjects))]

					var availableFaculty []models.Faculty
					if subjectFaculty, exists := facultyForSubjectSection[subject.ID]; exists {
						if sectionFaculty, exists := subjectFaculty[sectionID]; exists {
							availableFaculty = sectionFaculty
						}
					}

					if len(availableFaculty) == 0 {
						continue
					}

					selectedFaculty := availableFaculty[rand.Intn(len(availableFaculty))]

					if assignedTime, exists := facultyAssignments[day.DayName][selectedFaculty.FacultyName]; exists && assignedTime == startTime {
						continue
					}

					var selectedClassroom models.Classroom
					if len(classroomsByDept[subject.DepartmentID]) > 0 {
						selectedClassroom = classroomsByDept[subject.DepartmentID][rand.Intn(len(classroomsByDept[subject.DepartmentID]))]
					}

					entry := models.TimetableEntry{
						DayName:      day.DayName,
						StartTime:    startTime,
						EndTime:      hours[i].EndTime,
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

					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
					periodsLeft[subject.Name]--
					subjectDailyCount[day.DayName][subject.Name]++
					subjectsAssigned[day.DayName][subject.Name] = true
					facultyAssignments[day.DayName][selectedFaculty.FacultyName] = startTime
				}
			}

			// Check completion
			allAssigned := true
			for _, remaining := range periodsLeft {
				if remaining > 0 {
					allAssigned = false
					break
				}
			}

			if allAssigned {
				return timetable
			}
		}

		//log.Println("Warning: Could not assign all periods within retry limit")
		return nil

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
				return generateRandomTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, sections, semesters, departmentID, semesterID, academicYearID, sectionID)
			}
		}

		// Fetch the first section's timetable
		firstSectionTimetable, err := FetchSectionTimetable(departmentID, semesterID, academicYearID, firstSectionID)
		if err != nil {
			fmt.Println("Error fetching first section timetable:", err)
			return nil
		}

		labVenues, err := lab.GetLabVenue()
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
			return generateRandomTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, sections, semesters, departmentID, semesterID, academicYearID, sectionID)
		}

		// Create a new timetable based on first section's schedule
		timetable := make(map[string]map[string][]models.TimetableEntry)

		// Initialize timetable structure
		for _, day := range days {
			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
		}

		// Track faculty assignments to avoid conflicts
		facultyAssignments := make(map[string]map[string]string)   // day -> faculty -> timeslot
		classroomAssignments := make(map[string]map[string]string) // day -> classroom -> timeslot
		//facultyDailyCount := make(map[string]map[string]int)       // day -> faculty -> count
		labVenueAssignments := make(map[string]map[string]string)

		// Initialize tracking structures
		for _, day := range days {
			facultyAssignments[day.DayName] = make(map[string]string)
			classroomAssignments[day.DayName] = make(map[string]string)
			//facultyDailyCount[day] = make(map[string]int)
			labVenueAssignments[day.DayName] = make(map[string]string)
			if skips, ok := skipTimetable[day.DayName]; ok {
				for startTime := range skips {
					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime])
				}
			}
		}

		// Store failed lab allocations for swapping
		failedLabAllocations := []struct {
			Day         string
			StartTime   string
			SubjectID   int
			SubjectName string
			Entry       struct {
				DayName     string
				StartTime   string
				EndTime     string
				SubjectName string
				Status      int
				CourseCode  string
			}
		}{}

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
					//var labOccupancyTracker = make(map[string]map[string]map[int]int)

					// Replace your lab venue allocation section with this:

					if entry.Status == 0 && len(labVenues) > 0 { // Check lab subjects
						var selectedLabVenue models.LabVenue
						var foundLab bool = false

						// Get current lab occupancy from database for this academic year
						labOccupancy, err := lab.GetLabOccupancyFromDB(day, startTime, academicYearID)
						if err != nil {
							fmt.Printf("Error getting lab occupancy: %v\n", err)
							continue
						}

						// Filter lab venues that match the subject ID
						matchingLabVenues := []models.LabVenue{}
						for _, lab := range labVenues {
							if lab.SubjectID == subjectID { // Use subjectID instead of entry.SubjectID
								matchingLabVenues = append(matchingLabVenues, lab)
							}
						}

						if len(matchingLabVenues) == 0 {
							fmt.Printf("Warning: No lab venues found for subject ID %d\n", subjectID)
							continue
						}

						// Check each matching lab venue for availability based on max_sections
						for _, lab := range matchingLabVenues {
							// Get current occupancy from database
							currentOccupancy := labOccupancy[lab.LabVenue] // Use lab name from database

							// Check if this lab has capacity (current occupancy < max_sections)
							if currentOccupancy < lab.MaxSections {
								selectedLabVenue = lab
								foundLab = true
								log.Printf("Allocated lab: %s (ID: %d) for subject %d at %s %s. Current occupancy: %d/%d (Academic Year: %d)",
									lab.LabVenue, lab.ID, subjectID, day, startTime,
									currentOccupancy+1, lab.MaxSections, academicYearID)
								break
							} else {
								log.Printf("Lab %s is at max capacity: %d/%d at %s %s (Academic Year: %d)",
									lab.LabVenue, currentOccupancy, lab.MaxSections, day, startTime, academicYearID)
							}
						}

						if !foundLab {
							fmt.Printf("Warning: All lab venues for subject %d are at max capacity at %s %s (Academic Year: %d). Adding to failed allocations for swapping.\n",
								subjectID, day, startTime, academicYearID)

							// Store failed allocation for potential swapping
							failedLabAllocations = append(failedLabAllocations, struct {
								Day         string
								StartTime   string
								SubjectID   int
								SubjectName string
								Entry       struct {
									DayName     string
									StartTime   string
									EndTime     string
									SubjectName string
									Status      int
									CourseCode  string
								}
							}{
								Day:         day,
								StartTime:   startTime,
								SubjectID:   subjectID,
								SubjectName: entry.SubjectName,
								Entry: struct {
									DayName     string
									StartTime   string
									EndTime     string
									SubjectName string
									Status      int
									CourseCode  string
								}{
									DayName:     entry.DayName,
									StartTime:   entry.StartTime,
									EndTime:     entry.EndTime,
									SubjectName: entry.SubjectName,
									Status:      entry.Status,
									CourseCode:  entry.CourseCode,
								},
							})
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

		// Process failed lab allocations through swapping mechanism
		// Process failed lab allocations through swapping mechanism
		for _, failedAlloc := range failedLabAllocations {
			fmt.Printf("Attempting to swap lab subject %s from %s %s\n", failedAlloc.SubjectName, failedAlloc.Day, failedAlloc.StartTime)

			// Find all consecutive periods for the failed lab subject
			failedConsecutivePeriods := findConsecutiveLabPeriods(timetable, firstSectionTimetable, failedAlloc.Day, failedAlloc.StartTime, failedAlloc.SubjectName, sectionID, hours)
			fmt.Printf("Failed subject %s has %d consecutive periods: %v\n", failedAlloc.SubjectName, len(failedConsecutivePeriods), failedConsecutivePeriods)

			// Find existing lab subjects in timetable that can be swapped
			swapFound := false
			for swapDay, timeSlots := range timetable {
				for swapStartTime, entries := range timeSlots {
					for _, existingEntry := range entries {
						// Check if this is a lab subject from same dept, semester, academic year, section
						if existingEntry.Status == 0 && // Lab subject
							existingEntry.DepartmentID == departmentID &&
							existingEntry.SemesterID == semesterID &&
							existingEntry.AcademicYear == academicYearID &&
							existingEntry.SectionID == sectionID &&
							(swapDay != failedAlloc.Day || swapStartTime != failedAlloc.StartTime) { // Different time slot

							// Find all consecutive periods for the existing lab subject
							swapConsecutivePeriods := findConsecutiveLabPeriods(timetable, firstSectionTimetable, swapDay, swapStartTime, existingEntry.SubjectName, sectionID, hours)
							fmt.Printf("Existing subject %s has %d consecutive periods: %v\n", existingEntry.SubjectName, len(swapConsecutivePeriods), swapConsecutivePeriods)

							// Only proceed if both subjects have the same number of consecutive periods
							if len(swapConsecutivePeriods) != len(failedConsecutivePeriods) {
								fmt.Printf("Cannot swap: Different number of consecutive periods. Swap subject has %d periods, Failed subject has %d periods\n",
									len(swapConsecutivePeriods), len(failedConsecutivePeriods))
								continue
							}

							// Find subject ID of existing entry
							var existingSubjectID int
							for _, subj := range subjects {
								if subj.Name == existingEntry.SubjectName {
									existingSubjectID = subj.ID
									break
								}
							}

							// Check lab venue availability for all consecutive periods
							allVenuesAvailable := true

							// Check if failed subject's lab venue is available at ALL existing entry's consecutive time slots
							for _, swapPeriod := range swapConsecutivePeriods {
								failedLabOccupancy, err := lab.GetLabOccupancyFromDB(swapDay, swapPeriod, academicYearID)
								if err != nil {
									allVenuesAvailable = false
									break
								}

								// Get failed subject's lab venues
								failedSubjectLabVenues := []models.LabVenue{}
								for _, labVen := range labVenues {
									if labVen.SubjectID == failedAlloc.SubjectID {
										failedSubjectLabVenues = append(failedSubjectLabVenues, labVen)
									}
								}

								var failedSubjectCanFitAtSwapTime bool = false
								for _, failedLabVen := range failedSubjectLabVenues {
									if failedLabOccupancy[failedLabVen.LabVenue] < failedLabVen.MaxSections {
										failedSubjectCanFitAtSwapTime = true
										break
									}
								}

								if !failedSubjectCanFitAtSwapTime {
									allVenuesAvailable = false
									break
								}
							}

							if !allVenuesAvailable {
								continue
							}

							// Check if existing subject's lab venue is available at ALL failed subject's consecutive time slots
							for _, failedPeriod := range failedConsecutivePeriods {
								existingLabOccupancy, err := lab.GetLabOccupancyFromDB(failedAlloc.Day, failedPeriod, academicYearID)
								if err != nil {
									allVenuesAvailable = false
									break
								}

								existingSubjectLabVenues := []models.LabVenue{}
								for _, labVen := range labVenues {
									if labVen.SubjectID == existingSubjectID {
										existingSubjectLabVenues = append(existingSubjectLabVenues, labVen)
									}
								}

								var existingSubjectCanFitAtFailedTime bool = false
								for _, existingLabVen := range existingSubjectLabVenues {
									if existingLabOccupancy[existingLabVen.LabVenue] < existingLabVen.MaxSections {
										existingSubjectCanFitAtFailedTime = true
										break
									}
								}

								if !existingSubjectCanFitAtFailedTime {
									allVenuesAvailable = false
									break
								}
							}

							if !allVenuesAvailable {
								continue
							}

							// Check faculty availability for ALL consecutive periods
							var failedSubjectFacultyAvailable bool = true
							var existingSubjectFacultyAvailable bool = true

							// Check if failed subject's faculty is available at ALL swap time periods
							for _, swapPeriod := range swapConsecutivePeriods {
								facultyFound := false
								for _, fac := range faculty {
									if assignedTime, exists := facultyAssignments[swapDay][fac.FacultyName]; !exists || assignedTime != swapPeriod {
										// Check if faculty can teach failed subject
										for _, fs := range facultySubjects {
											if fs.FacultyID == fac.ID && fs.SubjectID == failedAlloc.SubjectID &&
												fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
												fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
												facultyFound = true
												break
											}
										}
										if facultyFound {
											break
										}
									}
								}
								if !facultyFound {
									failedSubjectFacultyAvailable = false
									break
								}
							}

							// Check if existing subject's faculty is available at ALL failed time periods
							for _, failedPeriod := range failedConsecutivePeriods {
								facultyFound := false
								for _, fac := range faculty {
									if assignedTime, exists := facultyAssignments[failedAlloc.Day][fac.FacultyName]; !exists || assignedTime != failedPeriod {
										// Check if faculty can teach existing subject
										for _, fs := range facultySubjects {
											if fs.FacultyID == fac.ID && fs.SubjectID == existingSubjectID &&
												fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
												fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
												facultyFound = true
												break
											}
										}
										if facultyFound {
											break
										}
									}
								}
								if !facultyFound {
									existingSubjectFacultyAvailable = false
									break
								}
							}

							// If all conditions met, perform the swap
							if failedSubjectFacultyAvailable && existingSubjectFacultyAvailable {

								fmt.Printf("Swapping %d consecutive periods: %s (%v) <-> %s (%v)\n",
									len(swapConsecutivePeriods), failedAlloc.SubjectName, failedConsecutivePeriods,
									existingEntry.SubjectName, swapConsecutivePeriods)

								// Get lab venues for both subjects (get the first available one for each)
								var failedSubjectLabVenue models.LabVenue
								failedLabOccupancy, _ := lab.GetLabOccupancyFromDB(swapDay, swapConsecutivePeriods[0], academicYearID)
								for _, labVen := range labVenues {
									if labVen.SubjectID == failedAlloc.SubjectID {
										if failedLabOccupancy[labVen.LabVenue] < labVen.MaxSections {
											failedSubjectLabVenue = labVen
											break
										}
									}
								}

								var existingSubjectLabVenue models.LabVenue
								existingLabOccupancy, _ := lab.GetLabOccupancyFromDB(failedAlloc.Day, failedConsecutivePeriods[0], academicYearID)
								for _, labVen := range labVenues {
									if labVen.SubjectID == existingSubjectID {
										if existingLabOccupancy[labVen.LabVenue] < labVen.MaxSections {
											existingSubjectLabVenue = labVen
											break
										}
									}
								}

								// Get faculty for both subjects (get the first available one for each)
								var failedSubjectFaculty models.Faculty
								for _, fac := range faculty {
									if assignedTime, exists := facultyAssignments[swapDay][fac.FacultyName]; !exists || assignedTime != swapConsecutivePeriods[0] {
										for _, fs := range facultySubjects {
											if fs.FacultyID == fac.ID && fs.SubjectID == failedAlloc.SubjectID &&
												fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
												fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
												failedSubjectFaculty = fac
												break
											}
										}
										if failedSubjectFaculty.ID != 0 {
											break
										}
									}
								}

								var existingSubjectFaculty models.Faculty
								for _, fac := range faculty {
									if assignedTime, exists := facultyAssignments[failedAlloc.Day][fac.FacultyName]; !exists || assignedTime != failedConsecutivePeriods[0] {
										for _, fs := range facultySubjects {
											if fs.FacultyID == fac.ID && fs.SubjectID == existingSubjectID &&
												fs.DepartmentID == departmentID && fs.SemesterID == semesterID &&
												fs.AcademicYear == academicYearID && fs.SectionID == sectionID {
												existingSubjectFaculty = fac
												break
											}
										}
										if existingSubjectFaculty.ID != 0 {
											break
										}
									}
								}

								// STEP 1: Remove ALL consecutive periods of existing subject from swap time slots
								for _, swapPeriod := range swapConsecutivePeriods {
									if entries, exists := timetable[swapDay][swapPeriod]; exists {
										var newEntries []models.TimetableEntry
										for _, entry := range entries {
											if !(entry.SubjectName == existingEntry.SubjectName &&
												entry.SectionID == sectionID &&
												entry.Status == 0) {
												newEntries = append(newEntries, entry)
											}
										}
										timetable[swapDay][swapPeriod] = newEntries
									}
									// Clear faculty assignment
									delete(facultyAssignments[swapDay], existingEntry.FacultyName)
								}

								// STEP 2: Remove ALL consecutive periods of failed subject from failed time slots
								for _, failedPeriod := range failedConsecutivePeriods {
									if entries, exists := timetable[failedAlloc.Day][failedPeriod]; exists {
										var newEntries []models.TimetableEntry
										for _, entry := range entries {
											if !(entry.SubjectName == failedAlloc.SubjectName &&
												entry.SectionID == sectionID &&
												entry.Status == 0) {
												newEntries = append(newEntries, entry)
											}
										}
										timetable[failedAlloc.Day][failedPeriod] = newEntries
									}
								}

								// STEP 3: Add failed subject to ALL swap time periods
								for _, swapPeriod := range swapConsecutivePeriods {
									// Find the correct end time for the swap time slot
									var swapEndTime string
									for _, hour := range hours {
										if hour.StartTime == swapPeriod {
											swapEndTime = hour.EndTime
											break
										}
									}

									// Create new entry for failed subject at swap time
									newFailedEntry := models.TimetableEntry{
										DayName:      swapDay,
										StartTime:    swapPeriod,
										EndTime:      swapEndTime,
										SubjectName:  failedAlloc.SubjectName,
										FacultyName:  failedSubjectFaculty.FacultyName,
										Classroom:    failedSubjectLabVenue.LabVenue,
										Status:       failedAlloc.Entry.Status,
										SemesterID:   semesterID,
										DepartmentID: departmentID,
										AcademicYear: academicYearID,
										CourseCode:   failedAlloc.Entry.CourseCode,
										SectionID:    sectionID,
									}

									if _, ok := timetable[swapDay]; !ok {
										timetable[swapDay] = make(map[string][]models.TimetableEntry)
									}
									if _, ok := timetable[swapDay][swapPeriod]; !ok {
										timetable[swapDay][swapPeriod] = []models.TimetableEntry{}
									}
									timetable[swapDay][swapPeriod] = append(timetable[swapDay][swapPeriod], newFailedEntry)

									// Update faculty assignments for swap periods
									facultyAssignments[swapDay][failedSubjectFaculty.FacultyName] = swapPeriod
								}

								// STEP 4: Add existing subject to ALL failed time periods
								for _, failedPeriod := range failedConsecutivePeriods {
									// Find the correct end time for the failed time slot
									var failedEndTime string
									for _, hour := range hours {
										if hour.StartTime == failedPeriod {
											failedEndTime = hour.EndTime
											break
										}
									}

									// Create entry for existing subject at failed time
									newExistingEntry := models.TimetableEntry{
										DayName:      failedAlloc.Day,
										StartTime:    failedPeriod,
										EndTime:      failedEndTime,
										SubjectName:  existingEntry.SubjectName,
										FacultyName:  existingSubjectFaculty.FacultyName,
										Classroom:    existingSubjectLabVenue.LabVenue,
										Status:       existingEntry.Status,
										SemesterID:   semesterID,
										DepartmentID: departmentID,
										AcademicYear: academicYearID,
										CourseCode:   existingEntry.CourseCode,
										SectionID:    sectionID,
									}

									if _, ok := timetable[failedAlloc.Day]; !ok {
										timetable[failedAlloc.Day] = make(map[string][]models.TimetableEntry)
									}
									if _, ok := timetable[failedAlloc.Day][failedPeriod]; !ok {
										timetable[failedAlloc.Day][failedPeriod] = []models.TimetableEntry{}
									}
									timetable[failedAlloc.Day][failedPeriod] = append(timetable[failedAlloc.Day][failedPeriod], newExistingEntry)

									// Update faculty assignments for failed periods
									facultyAssignments[failedAlloc.Day][existingSubjectFaculty.FacultyName] = failedPeriod
								}

								swapFound = true
								fmt.Printf("Successfully swapped lab subjects with %d consecutive periods: %s <-> %s\n",
									len(swapConsecutivePeriods), failedAlloc.SubjectName, existingEntry.SubjectName)
								break
							}
						}
					}
					if swapFound {
						break
					}
				}
				if swapFound {
					break
				}
			}

			if !swapFound {
				fmt.Printf("Could not find suitable swap for lab subject %s at %s %s\n", failedAlloc.SubjectName, failedAlloc.Day, failedAlloc.StartTime)
			}
		}

		// Verify that all periods have assignments
		for _, day := range days {
			for _, hour := range hours {
				startTime := hour.StartTime
				if _, ok := timetable[day.DayName][startTime]; !ok || len(timetable[day.DayName][startTime]) == 0 {
					//fmt.Printf("Warning: No assignment for %s at %s in section %d\n", day, startTime, sectionID)
				}
			}
		}

		return timetable
	}
}
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

func IsLabPeriodAvailable(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
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

func findConsecutiveLabPeriods(timetable map[string]map[string][]models.TimetableEntry, firstSectionTimetable map[string]map[string][]models.TimetableEntry, day, startTime, subjectName string, sectionID int, hours []models.Hour) []string {
	consecutivePeriods := []string{startTime}

	// Find the index of current start time
	var currentIndex int = -1
	for i, hour := range hours {
		if hour.StartTime == startTime {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return consecutivePeriods
	}

	// Check forward consecutive periods
	for i := currentIndex + 1; i < len(hours); i++ {
		nextStartTime := hours[i].StartTime
		found := false

		// Check in current timetable
		if entries, exists := timetable[day][nextStartTime]; exists {
			for _, entry := range entries {
				if entry.SubjectName == subjectName && entry.SectionID == sectionID && entry.Status == 0 {
					consecutivePeriods = append(consecutivePeriods, nextStartTime)
					found = true
					break
				}
			}
		}

		// If not found in current timetable, check in first section timetable
		if !found {
			if firstSectionEntries, exists := firstSectionTimetable[day][nextStartTime]; exists {
				for _, entry := range firstSectionEntries {
					if entry.SubjectName == subjectName && entry.Status == 0 {
						consecutivePeriods = append(consecutivePeriods, nextStartTime)
						found = true
						break
					}
				}
			}
		}

		// If no consecutive period found, break
		if !found {
			break
		}
	}

	// Check backward consecutive periods (in case we started from middle of a lab session)
	for i := currentIndex - 1; i >= 0; i-- {
		prevStartTime := hours[i].StartTime
		found := false

		// Check in current timetable
		if entries, exists := timetable[day][prevStartTime]; exists {
			for _, entry := range entries {
				if entry.SubjectName == subjectName && entry.SectionID == sectionID && entry.Status == 0 {
					// Insert at beginning to maintain chronological order
					consecutivePeriods = append([]string{prevStartTime}, consecutivePeriods...)
					found = true
					break
				}
			}
		}

		// If not found in current timetable, check in first section timetable
		if !found {
			if firstSectionEntries, exists := firstSectionTimetable[day][prevStartTime]; exists {
				for _, entry := range firstSectionEntries {
					if entry.SubjectName == subjectName && entry.Status == 0 {
						// Insert at beginning to maintain chronological order
						consecutivePeriods = append([]string{prevStartTime}, consecutivePeriods...)
						found = true
						break
					}
				}
			}
		}

		// If no consecutive period found, break
		if !found {
			break
		}
	}

	return consecutivePeriods
}
