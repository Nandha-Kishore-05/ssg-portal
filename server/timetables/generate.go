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
		var status, semesterID, departmentID, academicYearID,sectionID int

		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom, &status, &semesterID, &departmentID, &academicYearID, &courseCode,&sectionID); err != nil {
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
			SectionID: sectionID,
		}

		existingTimetable[facultyName][dayName] = append(existingTimetable[facultyName][dayName], entry)
	}

	return existingTimetable, nil
}

func FetchTimetableSkips(departmentID int, semesterID int, academicYearID int) (map[string]map[string]models.TimetableEntry, error) {
	skipEntries := make(map[string]map[string]models.TimetableEntry)

	query := `
        SELECT day_name, start_time, end_time, subject_name, faculty_name, semester_id, department_id ,classroom,status,academic_year,course_code
        FROM timetable_skips 
        WHERE department_id = ? AND semester_id = ? AND  academic_year = ?`

	rows, err := config.Database.Query(query, departmentID, semesterID, &academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var dayName, startTime, endTime, subjectName, facultyName, classroom, courseCode string
		var status, academicYearID int

		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &semesterID, &departmentID, &classroom, &status, &academicYearID, &courseCode); err != nil {
			return nil, err
		}

		if _, exists := skipEntries[dayName]; !exists {
			skipEntries[dayName] = make(map[string]models.TimetableEntry)
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
			CourseCode:   courseCode,
		
		}

		skipEntries[dayName][startTime] = entry
	}

	return skipEntries, nil
}

func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject, semesters []models.Semester, section []models.Section, academicYear []models.AcademicYear, departmentID int, semesterID int, academicYearID int, sectionID int) map[string]map[string][]models.TimetableEntry {

	existingTimetable, err := FetchExistingTimetable()
	if err != nil {
		fmt.Println("Error fetching existing timetable:", err)
		return nil
	}
	skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID)
	if err != nil {
		fmt.Println("Error fetching timetable skips:", err)
		return nil
	}

	if existingTimetable != nil && len(existingTimetable) > 0 {

		for {
			timetable := make(map[string]map[string][]models.TimetableEntry)
			subjectsAssigned := make(map[string]map[string]bool)
			periodsLeft := make(map[string]int)
			status0Assignments := make(map[string]map[string]bool)
			facultyAssignments := make(map[string]map[string]string)
			facultyDailyCount := make(map[string]map[string]int)


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
				facultyDailyCount[day.DayName] = make(map[string]int)
				if skips, ok := skipTimetable[day.DayName]; ok {
					for startTime, entry := range skips {
						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						subjectsAssigned[day.DayName][entry.SubjectName] = true
						periodsLeft[entry.SubjectName]--
						facultyDailyCount[day.DayName][entry.FacultyName]++
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
							if periodsLeft[subject.Name] > 0 && (!subjectsAssigned[day.DayName][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) {
								if subject.Status == 1 && subjectsAssigned[day.DayName][subject.Name] {
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

						if _, ok := timetable[day.DayName][startTime]; ok {
							if len(timetable[day.DayName][startTime]) > 0 {
								continue
							}
						}

						var availableFaculty []models.Faculty
						for _, fac := range faculty {
							for _, fs := range facultySubjects {
								if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID {
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

						if facultyDailyCount[day.DayName][selectedFaculty.FacultyName] >= 2 {
							continue
						}

						if assignedClassroom, exists := facultyAssignments[day.DayName][selectedFaculty.FacultyName]; exists && assignedClassroom == startTime {
							continue
						}

						var selectedClassroom models.Classroom
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
						if !Available(existingTimetable, day.DayName, startTime, availableFaculty[facultyIndex].FacultyName) {
							continue
						}
						entry := models.TimetableEntry{
							DayName:      day.DayName,
							StartTime:    startTime,
							EndTime:      endTime,
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

						if _, ok := timetable[day.DayName]; !ok {
							timetable[day.DayName] = make(map[string][]models.TimetableEntry)
						}
						if _, ok := timetable[day.DayName][startTime]; !ok {
							timetable[day.DayName][startTime] = []models.TimetableEntry{}
						}

						if subject.Status == 0 {
							if _, ok := status0Assignments[subject.Name][startTime]; !ok {
								if i < len(hours)-1 {
									nextHourIndex := (hourIndex + 1) % len(hours)
									nextStartTime := hours[nextHourIndex].StartTime
									if IsPeriodAvailable(existingTimetable, day.DayName, nextStartTime, "") {
										entry2 := entry
										entry2.StartTime = nextStartTime
										entry2.EndTime = hours[nextHourIndex].EndTime

										timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
										timetable[day.DayName][nextStartTime] = append(timetable[day.DayName][nextStartTime], entry2)
										periodsLeft[subject.Name] -= 2
										subjectsAssigned[day.DayName][subject.Name] = true
										status0Assignments[subject.Name][startTime] = true
										status0Assignments[subject.Name][nextStartTime] = true
										facultyAssignments[day.DayName][selectedFaculty.FacultyName] = nextStartTime
										assigned = true
										break
									}
								}
								continue
							}
						}

						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						periodsLeft[subject.Name]--
						subjectsAssigned[day.DayName][subject.Name] = true
						if subject.Status == 0 {
							status0Assignments[subject.Name][startTime] = true
						}
						facultyAssignments[day.DayName][selectedFaculty.FacultyName] = startTime

						assigned = true
						break
					}
					if !assigned {
						fmt.Printf("Warning: Could not assign a subject for %s during period %d\n", day.DayName, i+1)
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
					if len(timetable[day.DayName][startTime]) == 0 {
						periodsFilled = false
						break
					}
				}
				if !periodsFilled {
					break
				}
			}

			// Regenerate if not all periods are filled or if there are conflicts
			if allAssigned && periodsFilled && !CheckTimetableConflicts(timetable, existingTimetable) {
				return timetable
			}
		}
	} else {
		// Call generateRandomTimetable if no existing timetable is found
		return generateRandomTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, section, semesters, departmentID, semesterID, academicYearID, sectionID)
	}
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
	skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID)
	if err != nil {
		fmt.Println("Error fetching timetable skips:", err)
		return nil
	}


	maxAttempts := len(subjects) * len(hours)

	generate := func() FacultyBasedTimetable {
		timetable := make(FacultyBasedTimetable)
		subjectsAssigned := make(map[string]map[string]bool)
		periodsLeft := make(map[string]int)
		status0Assignments := make(map[string]map[string]bool)
		labSubjectAssigned := make(map[string]bool)

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

		for _, day := range days {
			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
			subjectsAssigned[day.DayName] = make(map[string]bool)
			labSubjectAssigned[day.DayName] = false
			if skips, ok := skipTimetable[day.DayName]; ok {
				for startTime, entry := range skips {
					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
					subjectsAssigned[day.DayName][entry.SubjectName] = true
					periodsLeft[entry.SubjectName]--
				}
			}
		}

		rand.Seed(time.Now().UnixNano())


		for _, day := range days {
			for i := 0; i < len(hours); i++ {
				for attempts := 0; attempts < maxAttempts; attempts++ {
					var filteredLabSubjects []models.Subject
					for _, subject := range labSubjects {
						if periodsLeft[subject.Name] > 0 &&
							(!subjectsAssigned[day.DayName][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) &&
							!labSubjectAssigned[day.DayName] {
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

					if len(timetable[day.DayName][startTime]) > 0 {
						continue
					}

					if subject.Status == 0 && i < len(hours)-1 {
						nextStartTime := hours[i+1].StartTime
						nextEndTime := hours[i+1].EndTime

						if IsPeriodAvailable(timetable, day.DayName, nextStartTime, "") {
							facultyName := selectRandomFaculty(faculty, subject, facultySubjects)
							if facultyName == "" {
								fmt.Println("Error: No faculty available for lab subject", subject.Name)
								return nil
							}

							classroomName := selectRandomClassroom(classrooms)
							if classroomName == "" {
								fmt.Println("Error: No classroom found for lab subject", subject.Name)
								return nil
							}

							entry1 := models.TimetableEntry{
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
								SectionID: sectionID,
							}

							entry2 := models.TimetableEntry{
								DayName:      day.DayName,
								StartTime:    nextStartTime,
								EndTime:      nextEndTime,
								SubjectName:  subject.Name,
								FacultyName:  entry1.FacultyName,
								Classroom:    entry1.Classroom,
								Status:       subject.Status,
								SemesterID:   entry1.SemesterID,
								DepartmentID: departmentID,
								AcademicYear: academicYearID,
								CourseCode:   subject.CourseCode,
								SectionID: sectionID,
							}

							timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry1)
							timetable[day.DayName][nextStartTime] = append(timetable[day.DayName][nextStartTime], entry2)

							periodsLeft[subject.Name] -= 2
							subjectsAssigned[day.DayName][subject.Name] = true
							status0Assignments[subject.Name][startTime] = true
							labSubjectAssigned[day.DayName] = true
							break
						}
					}
				}
			}
		}


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

					facultyName := selectRandomFaculty(faculty, subject, facultySubjects)
					if facultyName == "" {
						fmt.Println("Error: No faculty available for non-lab subject", subject.Name)
						return nil
					}


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
						SectionID: sectionID,
					}

					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
					periodsLeft[subject.Name]--
					subjectsAssigned[day.DayName][subject.Name] = true
					break
				}
			}
		}

		return timetable
	}

	for {
		timetable := generate()
		allPeriodsFilled := true
		for _, day := range days {
			for _, hour := range hours {
				startTime := hour.StartTime
				if len(timetable[day.DayName][startTime]) == 0 {
					fmt.Printf("Empty period found for %s at %s\n", day.DayName, startTime)
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


func selectRandomFaculty(facultyList []models.Faculty, subject models.Subject, facultySubjects []models.FacultySubject) string {
	var availableFaculty []models.Faculty
	for _, fac := range facultyList {
		for _, fs := range facultySubjects {
			if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID {
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