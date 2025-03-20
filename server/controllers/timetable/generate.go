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

func AreAllSubjectFacultyAvailableAcrossSections(
	existingTimetable map[string]map[string][]models.TimetableEntry,
	subjectID int,
	dayName string,
	startTime string,
	departmentID int,
	semesterID int,
	academicYearID int,
	sectionMap map[int]bool,
	facultySubjectMap map[int]map[int][]models.Faculty) bool {

	log.Println("Checking faculty availability across sections")

	// Get all sections from the in-memory map instead of database
	sections := make([]int, 0)
	for sectionID := range sectionMap {
		// Filter by department, semester, and academic year (assuming this logic is handled elsewhere)
		sections = append(sections, sectionID)
	}

	// For each section, check if the faculty assigned to this subject is available
	for _, sectionID := range sections {
		log.Println("Checking section:", sectionID)

		// Get the faculty assigned to this subject from the hashmap
		faculty := facultySubjectMap[subjectID]
		if faculty == nil {
			log.Println("No faculty found for subject:", subjectID)
			continue
		}

		// Convert faculty map into a list of names
		var facultyNames []string
		for _, facultyList := range faculty {
			for _, f := range facultyList {
				facultyNames = append(facultyNames, f.FacultyName)
			}
		}

		log.Println("Total faculties:", len(facultyNames))

		// Check availability of all faculties at the given day & time
		unavailableFaculties := IsFacultyAvailable(existingTimetable, facultyNames, dayName, startTime)

		if len(unavailableFaculties) > 0 {
			log.Println("Unavailable faculties:", unavailableFaculties)
			return false // If at least one faculty is unavailable, reject the slot
		}
	}

	return true // All faculties across all sections are available for this subject
}

func IsFacultyAvailable(
	existingTimetable map[string]map[string][]models.TimetableEntry,
	facultyList []string,
	dayName string,
	startTime string) []string { // Returns a slice of unavailable faculty names

	log.Println("Checking faculty availability")

	var unavailableFaculties []string

	for _, facultyName := range facultyList {
		log.Println("Checking faculty:", facultyName)

		if facultyEntries, exists := existingTimetable[facultyName]; exists {
			if dayEntries, exists := facultyEntries[dayName]; exists {
				for _, entry := range dayEntries {
					if entry.StartTime == startTime {
						log.Println("Faculty", facultyName, "is unavailable at", startTime, "on", dayName)
						unavailableFaculties = append(unavailableFaculties, facultyName)
						break // No need to check further for this faculty
					}
				}
			}
		}
	}

	return unavailableFaculties // Return the list of unavailable faculties
}

// Function to build the section map in memory (called during initialization)
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

	existingTimetable, err := FetchExistingTimetable()
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

	// Check if sectionsInSameSemester has elements before accessing index 0
	if existingSectionID == 0 {
		log.Println("kishore")
		// Fetch the existing timetable and handle errors
		existingTimetable, err := FetchExistingTimetable()
		if err != nil {
			fmt.Println("Error fetching existing timetable:", err)
			return nil
		}

		// Precompute faculty availability once
		allFacultyNames := make([]string, 0, len(faculty))
		for _, f := range faculty {
			allFacultyNames = append(allFacultyNames, f.FacultyName)
		}
		log.Println("fac", len(allFacultyNames))
		facultyAvailabilityCache := PrecomputeFacultyAvailability(existingTimetable, allFacultyNames)

		// Build maps once
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

		// Fetch other required data
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

		// Precompute classroom assignments by department
		classroomsByDept := make(map[int][]models.Classroom)
		for _, cls := range classrooms {
			classroomsByDept[cls.DepartmentID] = append(classroomsByDept[cls.DepartmentID], cls)
		}

		// Precompute faculty for subjects and sections
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

		// Try to generate a valid timetable with retry logic
		for {
			timetable := make(map[string]map[string][]models.TimetableEntry)
			subjectsAssigned := make(map[string]map[string]bool)
			periodsLeft := make(map[string]int)
			status0Assignments := make(map[string]map[string]bool)
			facultyAssignments := make(map[string]map[string]string)
			subjectDailyCount := make(map[string]map[string]int)
			labAssigned := make(map[string]bool)

			// Initialize maps
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
				labAssigned[day.DayName] = false
				if skips, ok := skipTimetable[day.DayName]; ok {
					for startTime := range skips {
						timetable[day.DayName][startTime] = []models.TimetableEntry{}
					}
				}
			}

			rand.Seed(time.Now().UnixNano())

			// Process days - this is the main optimization target
			//	success := true
			for _, day := range days {
				for i := 0; i < len(hours); i++ {
					hourIndex := i % len(hours)
					startTime := hours[hourIndex].StartTime

					// Skip if already assigned or in skip list
					if _, exists := timetable[day.DayName][startTime]; exists && len(timetable[day.DayName][startTime]) > 0 {
						continue
					}

					// Filter eligible subjects - precompute once per slot
					var eligibleSubjects []models.Subject
					for _, subject := range subjects {
						// Check basic eligibility
						dailyLimit := 2
						if subject.Status != 0 {
							dailyLimit = 1
						}

						if periodsLeft[subject.Name] <= 0 ||
							subjectDailyCount[day.DayName][subject.Name] >= dailyLimit ||
							(subjectsAssigned[day.DayName][subject.Name] &&
								(subject.Status == 1 || (subject.Status == 0 && len(status0Assignments[subject.Name]) > 0))) {
							continue
						}

						// Lab-specific checks
						if subject.Status == 0 && labAssigned[day.DayName] {
							continue
						}

						// Check faculty availability using cache
						if !AreAllSubjectFacultyAvailableAcrossSectionsFast(
							facultyAvailabilityCache,
							facultySubjectMap,
							subject.ID,
							day.DayName,
							startTime,
							sectionMap) {
							continue
						}

						// Verify classroom availability
						if len(classroomsByDept[subject.DepartmentID]) > 0 {
							eligibleSubjects = append(eligibleSubjects, subject)
						}
					}

					// No eligible subjects for this slot
					if len(eligibleSubjects) == 0 {
						continue
					}

					// Pick a random subject
					subject := eligibleSubjects[rand.Intn(len(eligibleSubjects))]
					endTime := hours[hourIndex].EndTime

					// Get available faculty for this subject and section
					var availableFaculty []models.Faculty
					if subjectFaculty, exists := facultyForSubjectSection[subject.ID]; exists {
						if sectionFaculty, exists := subjectFaculty[sectionID]; exists {
							availableFaculty = sectionFaculty
						}
					}

					if len(availableFaculty) == 0 {
						continue
					}

					// Pick a random faculty
					selectedFaculty := availableFaculty[rand.Intn(len(availableFaculty))]

					// Check faculty schedule conflict
					if assignedTime, exists := facultyAssignments[day.DayName][selectedFaculty.FacultyName]; exists && assignedTime == startTime {
						continue
					}
   
					// Select classroom (lab or regular)
					var selectedClassroom models.Classroom
					if subject.Status == 0 && len(labVenues) > 0 {
						selectedLabVenue := labVenues[rand.Intn(len(labVenues))]
						selectedClassroom = models.Classroom{
							ID:            selectedLabVenue.ID,
							ClassroomName: selectedLabVenue.LabVenue,
							DepartmentID:  selectedLabVenue.DepartmentID,
							SemesterID:    selectedLabVenue.SemesterID,
						}
					} else if len(classroomsByDept[subject.DepartmentID]) > 0 {
						selectedClassroom = classroomsByDept[subject.DepartmentID][rand.Intn(len(classroomsByDept[subject.DepartmentID]))]
					}

					// Create the timetable entry
					entry := models.TimetableEntry{
						DayName:      day.DayName,
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

					// Handle lab subjects (which need two consecutive slots)
					if subject.Status == 0 {
						if i < len(hours)-1 {
							nextHourIndex := (hourIndex + 1) % len(hours)
							nextStartTime := hours[nextHourIndex].StartTime

							// Check availability of next slot
							if IsLabPeriodAvailable(existingTimetable, day.DayName, nextStartTime, "") {
								entry2 := entry
								entry2.StartTime = nextStartTime
								entry2.EndTime = hours[nextHourIndex].EndTime

								// Assign both slots
								timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
								timetable[day.DayName][nextStartTime] = append(timetable[day.DayName][nextStartTime], entry2)

								// Update tracking information
								periodsLeft[subject.Name] -= 2
								subjectsAssigned[day.DayName][subject.Name] = true
								status0Assignments[subject.Name][startTime] = true
								status0Assignments[subject.Name][nextStartTime] = true
								facultyAssignments[day.DayName][selectedFaculty.FacultyName] = nextStartTime
								subjectDailyCount[day.DayName][subject.Name] += 2
								labAssigned[day.DayName] = true
							} else {
								continue // Can't assign lab here
							}
						} else {
							continue // Not enough slots left for lab
						}
					} else {
						// Regular subject assignment
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						periodsLeft[subject.Name]--
						subjectDailyCount[day.DayName][subject.Name]++
						subjectsAssigned[day.DayName][subject.Name] = true
						facultyAssignments[day.DayName][selectedFaculty.FacultyName] = startTime
					}
				}
			}

			// Check if all subjects were fully assigned
			allAssigned := true
			for subjectName, remaining := range periodsLeft {
				if remaining > 0 {
					allAssigned = false
					log.Printf("Subject %s has %d periods left unassigned\n", subjectName, remaining)
					break
				}
			}

			if allAssigned {
				return timetable
			}
		}

		//log.Println("Failed to generate a valid timetable after multiple attempts")
		//return nil
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
