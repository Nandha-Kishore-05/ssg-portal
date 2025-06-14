package timetable

import (
	"fmt"
	"log"
	"math/rand"

	"ssg-portal/controllers/lab"
	"ssg-portal/models"
	"time"
)
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

	// // Extract working days into a slice of strings
	// var days []string
	// for _, wd := range workingDays {
	// 	days = append(days, wd.WorkingDate.Format("2006-01-02")) // Use date in "YYYY-MM-DD" format
	// }

	labVenues, err := lab.GetLabVenue()
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
		//facultyDailyCount := make(map[string]map[string]int)
		status0Assignments := make(map[string]map[string]bool)
		labSubjectAssigned := make(map[string]int)
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

		// Define the exact allowed lab slot pairs
		allowedLabPairs := [][2]string{
			{"08:45:00", "09:35:00"}, // 1st & 2nd periods
			{"10:40:00", "11:30:00"}, // 3rd & 4th periods  
			{"13:30:00", "14:20:00"}, // 5th & 6th periods
		}

		// Pre-find all valid consecutive slot pairs for labs
		labSlots := make(map[string][][2]string)
		
		// Group lab venues by subject for easier lookup
		labVenuesBySubject := make(map[int][]models.LabVenue)
		for _, labVenue := range labVenues {
			labVenuesBySubject[labVenue.SubjectID] = append(labVenuesBySubject[labVenue.SubjectID], labVenue)
		}

		// Iterate over days and apply timetable skips and manual timetable entries.
		for _, day := range days {
			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
			subjectsAssigned[day.DayName] = make(map[string]bool)
			//facultyDailyCount[day.DayName] = make(map[string]int)

			// Apply timetable skips
			if skips, ok := skipTimetable[day.DayName]; ok {
				for startTime, entries := range skips {
					for _, entry := range entries {
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						subjectsAssigned[day.DayName][entry.SubjectName] = true
					}
				}
			}

			// Apply manual timetable entries
			if manualEntries, ok := manualTimetable[day.DayName]; ok {
				for startTime, entries := range manualEntries {
					for _, entry := range entries {
						// Ensure you're not overwriting existing timetable entries
						if len(timetable[day.DayName][startTime]) == 0 {
							timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
							subjectsAssigned[day.DayName][entry.SubjectName] = true
						}
					}
				}
			}

			// Find available lab slots for this day
			for _, pair := range allowedLabPairs {
				slot1 := pair[0]
				slot2 := pair[1]
				
				// Check if both slots are free
				free1 := len(timetable[day.DayName][slot1]) == 0
				free2 := len(timetable[day.DayName][slot2]) == 0
				
				if free1 && free2 {
					labSlots[day.DayName] = append(labSlots[day.DayName], [2]string{slot1, slot2})
				}
			}
		}

		rand.Seed(time.Now().UnixNano())

		// Handle lab subjects first using the predefined slot pairs
		for _, day := range days {
			availableLabSlots := labSlots[day.DayName]
			
			// Shuffle the available lab slots randomly
			rand.Shuffle(len(availableLabSlots), func(i, j int) {
				availableLabSlots[i], availableLabSlots[j] = availableLabSlots[j], availableLabSlots[i]
			})

			for _, slotPair := range availableLabSlots {
				for attempts := 0; attempts < maxAttempts; attempts++ {
					var filteredLabSubjects []models.Subject

					// Find lab subjects that need to be scheduled
					for _, subject := range labSubjects {
						if periodsLeft[subject.Name] > 0 &&
							(!subjectsAssigned[day.DayName][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) &&
							labSubjectAssigned[day.DayName] < 2 {
							filteredLabSubjects = append(filteredLabSubjects, subject)
						}
					}

					if len(filteredLabSubjects) == 0 {
						break
					}

					subjectIndex := rand.Intn(len(filteredLabSubjects))
					subject := filteredLabSubjects[subjectIndex]

					startTime := slotPair[0]
					nextStartTime := slotPair[1]

					// Find the corresponding end times
					var endTime, nextEndTime string
					for _, hour := range hours {
						if hour.StartTime == startTime {
							endTime = hour.EndTime
						}
						if hour.StartTime == nextStartTime {
							nextEndTime = hour.EndTime
						}
					}

					// Check if slots are still available
					if len(timetable[day.DayName][startTime]) > 0 || len(timetable[day.DayName][nextStartTime]) > 0 {
						continue
					}

					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)

					if facultyName == "" {
						fmt.Println("Error: No faculty available for lab subject", subject.Name)
						continue
					}
					
					// if facultyDailyCount[day.DayName][facultyName] >= 2 {
					// 	continue // Skip to the next attempt if assigned twice
					// }

					// Find lab venue for this subject
					var labVenue models.LabVenue
					if venues, exists := labVenuesBySubject[subject.ID]; exists && len(venues) > 0 {
						labVenue = venues[0] // Take the first available venue
					}

					if labVenue.ID == 0 {
						continue // No lab venue found for this subject
					}

					entry1 := models.TimetableEntry{
						DayName:      day.DayName,
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
						DayName:      day.DayName,
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

					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry1)
					timetable[day.DayName][nextStartTime] = append(timetable[day.DayName][nextStartTime], entry2)

					periodsLeft[subject.Name] -= 2
					subjectsAssigned[day.DayName][subject.Name] = true
					status0Assignments[subject.Name][startTime] = true
					labSubjectAssigned[day.DayName]++
					
					// Ensure facultyAssignments[day] is initialized
					if _, exists := facultyAssignments[day.DayName]; !exists {
						facultyAssignments[day.DayName] = make(map[string]int)
					}

					// Safely increment the count for facultyName
					facultyAssignments[day.DayName][facultyName]++
					// facultyDailyCount[day.DayName][facultyName] += 2
					break
				}
			}
		}

		// Handle non-lab subjects
		for _, day := range days {
			for i := 0; i < len(hours); i++ {
				for attempts := 0; attempts < maxAttempts; attempts++ {
					var filteredNonLabSubjects []models.Subject
					for _, subject := range nonLabSubjects {
						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
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

					if len(timetable[day.DayName][startTime]) > 0 {
						continue
					}

					facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)

					if facultyName == "" {
						fmt.Println("Error: No faculty available for non-lab subject", subject.Name)
						return nil
					}
					
					// if facultyDailyCount[day.DayName][facultyName] >= 1 {
					// 	continue // Skip to the next attempt if assigned twice
					// }
					
					classroomName := selectRandomClassroom(classrooms)
					if classroomName == "" {
						fmt.Println("Error: No classroom found for non-lab subject", subject.Name)
						return nil
					}

					entry := models.TimetableEntry{
						DayName:      day.DayName,
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

					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
					periodsLeft[subject.Name]--
					subjectsAssigned[day.DayName][subject.Name] = true
					
					// Ensure facultyAssignments[day] is initialized
					if _, exists := facultyAssignments[day.DayName]; !exists {
						facultyAssignments[day.DayName] = make(map[string]int)
					}

					// Safely increment the count for facultyName
					facultyAssignments[day.DayName][facultyName]++
					//facultyDailyCount[day.DayName][facultyName]++
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