// package timetables

// import (
// 	"fmt"
// 	"math/rand"
// 	"ssg-portal/config"
// 	"ssg-portal/models"
// 	"time"
// )

// // FacultyBasedTimetable represents the structure of the timetable organized by faculty and day
// type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

// // FetchExistingTimetable fetches existing timetable data from the database
// func FetchExistingTimetable() (map[string]map[string][]models.TimetableEntry, error) {
// 	// Initialize the map to store existing timetable data
// 	existingTimetable := make(map[string]map[string][]models.TimetableEntry)

// 	// Query to fetch timetable data from the database
// 	rows, err := config.Database.Query("SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom, semester_id FROM timetable")
// 	if err != nil {
// 		// Return error if query fails
// 		return nil, err
// 	}
// 	defer rows.Close() // Ensure the rows are closed after processing

// 	// Iterate through the query results
// 	for rows.Next() {
// 		var dayName, startTime, endTime, subjectName, facultyName, classroom string
// 		var semesterID int

// 		// Scan the row data into variables
// 		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom, &semesterID); err != nil {
// 			return nil, err
// 		}

// 		// Initialize the map for a new faculty if it does not exist
// 		if _, exists := existingTimetable[facultyName]; !exists {
// 			existingTimetable[facultyName] = make(map[string][]models.TimetableEntry)
// 		}

// 		// Create a timetable entry from the row data
// 		entry := models.TimetableEntry{
// 			DayName:     dayName,
// 			StartTime:   startTime,
// 			EndTime:     endTime,
// 			SubjectName: subjectName,
// 			FacultyName: facultyName,
// 			Classroom:   classroom,
// 			SemesterID:  semesterID,
// 		}

// 		// Append the entry to the existing timetable
// 		existingTimetable[facultyName][dayName] = append(existingTimetable[facultyName][dayName], entry)
// 	}

// 	// Return the existing timetable data
// 	return existingTimetable, nil
// }

// // GenerateTimetable generates a timetable with the given constraints
// func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject, semesters []models.Semester) map[string]map[string][]models.TimetableEntry {
// 	for {
// 		timetable := make(FacultyBasedTimetable)
// 		subjectsAssigned := make(map[string]map[string]bool) // Track assigned subjects for each day
// 		periodsLeft := make(map[string]int)

// 		// Initialize maps
// 		for _, subject := range subjects {
// 			periodsLeft[subject.Name] = subject.Period
// 		}
// 		for _, day := range days {
// 			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
// 			subjectsAssigned[day.DayName] = make(map[string]bool)
// 		}

// 		rand.Seed(time.Now().UnixNano()) // Seed the random number generator

// 		// Assign subjects to periods
// 		for _, day := range days {
// 			for i := 0; i < 6; i++ { // 6 periods per day
// 				assigned := false
// 				for attempts := 0; attempts < 1000; attempts++ {
// 					// Filter subjects by department and semester
// 					var filteredSubjects []models.Subject
// 					for _, subject := range subjects {
// 						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
// 							// Check if there is an available classroom matching the subject's department and semester
// 							var validClassrooms []models.Classroom
// 							for _, cls := range classrooms {
// 								if cls.DepartmentID == subject.DepartmentID {
// 									// If the classroom has a valid semesterID
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

// 					// Randomly select a filtered subject
// 					subjectIndex := rand.Intn(len(filteredSubjects))
// 					subject := filteredSubjects[subjectIndex]

// 					hourIndex := i % len(hours)
// 					startTime := hours[hourIndex].StartTime
// 					endTime := hours[hourIndex].EndTime

// 					// Find available faculty for the selected subject
// 					var availableFaculty []models.Faculty
// 					for _, fac := range faculty {
// 						for _, fs := range facultySubjects {
// 							if fs.FacultyID == fac.ID && fs.SubjectID == subject.ID {
// 								availableFaculty = append(availableFaculty, fac)
// 								break
// 							}
// 						}
// 					}

// 					if len(availableFaculty) == 0 {
// 						continue
// 					}

// 					// Randomly select an available faculty
// 					facultyIndex := rand.Intn(len(availableFaculty))

// 					// Select a valid classroom
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

// 					entry := models.TimetableEntry{
// 						DayName:     day.DayName,
// 						StartTime:   startTime,
// 						EndTime:     endTime,
// 						SubjectName: subject.Name,
// 						FacultyName: availableFaculty[facultyIndex].FacultyName,
// 						Classroom:   selectedClassroom.ClassroomName,
// 						Status:      subject.Status,
// 						SemesterID:  selectedClassroom.SemesterID,
// 					}

// 					if _, ok := timetable[day.DayName]; !ok {
// 						timetable[day.DayName] = make(map[string][]models.TimetableEntry)
// 					}

// 					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
// 					periodsLeft[subject.Name]--
// 					subjectsAssigned[day.DayName][subject.Name] = true

// 					assigned = true
// 					break
// 				}
// 				if !assigned {
// 					fmt.Printf("Warning: Could not assign a subject for %s during period %d\n", day.DayName, i+1)
// 					// Additional handling for unassigned periods or fallback
// 				}
// 			}
// 		}

// 		// Check if all periods are assigned
// 		allAssigned := true
// 		for subjectName, remainingPeriods := range periodsLeft {
// 			if remainingPeriods > 0 {
// 				fmt.Printf("Warning: Subject %s has %d periods left unassigned.\n", subjectName, remainingPeriods)
// 				allAssigned = false
// 			}
// 		}

// 		if allAssigned {
// 			return timetable
// 		}

// 		fmt.Println("Regenerating timetable due to unassigned periods...")
// 	}
// }

// // transformTimetable transforms an existing timetable to the desired format

//	func transformTimetable(existingTimetable map[string]map[string][]models.TimetableEntry) map[string]map[string][]models.TimetableEntry {
//		// Example transformation function; adjust as needed for your use case
//		return existingTimetable
//	}
package timetables

import (
	"fmt"
	"math/rand"
	"ssg-portal/config"
	"ssg-portal/models"
	"time"
)

// FacultyBasedTimetable represents the structure of the timetable organized by faculty and day
type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

// FetchExistingTimetable fetches existing timetable data from the database
func FetchExistingTimetable() (map[string]map[string][]models.TimetableEntry, error) {
	existingTimetable := make(map[string]map[string][]models.TimetableEntry)

	rows, err := config.Database.Query("SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom, semester_id FROM timetable")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dayName, startTime, endTime, subjectName, facultyName, classroom string
		var semesterID int

		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom, &semesterID); err != nil {
			return nil, err
		}

		if _, exists := existingTimetable[facultyName]; !exists {
			existingTimetable[facultyName] = make(map[string][]models.TimetableEntry)
		}

		entry := models.TimetableEntry{
			DayName:     dayName,
			StartTime:   startTime,
			EndTime:     endTime,
			SubjectName: subjectName,
			FacultyName: facultyName,
			Classroom:   classroom,
			SemesterID:  semesterID,
		}

		existingTimetable[facultyName][dayName] = append(existingTimetable[facultyName][dayName], entry)
	}

	return existingTimetable, nil
}

// GenerateTimetable generates a timetable based on existing entries and constraints
func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject, semesters []models.Semester) map[string]map[string][]models.TimetableEntry {
	existingTimetable, err := FetchExistingTimetable()
	if err != nil {
		fmt.Println("Error fetching existing timetable:", err)
		return nil
	}

	if existingTimetable != nil && len(existingTimetable) > 0 {
		for {
			timetable := make(FacultyBasedTimetable)
			subjectsAssigned := make(map[string]map[string]bool)
			periodsLeft := make(map[string]int)

			for _, subject := range subjects {
				periodsLeft[subject.Name] = subject.Period
			}
			for _, day := range days {
				timetable[day.DayName] = make(map[string][]models.TimetableEntry)
				subjectsAssigned[day.DayName] = make(map[string]bool)
			}

			rand.Seed(time.Now().UnixNano())

			for _, day := range days {
				for i := 0; i < 6; i++ {
					assigned := false
					for attempts := 0; attempts < 1000; attempts++ {
						var filteredSubjects []models.Subject
						for _, subject := range subjects {
							if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
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

						if !IsPeriodAvailable(existingTimetable, day.DayName, startTime, availableFaculty[facultyIndex].FacultyName) {
							continue
						}

						entry := models.TimetableEntry{
							DayName:     day.DayName,
							StartTime:   startTime,
							EndTime:     endTime,
							SubjectName: subject.Name,
							FacultyName: availableFaculty[facultyIndex].FacultyName,
							Classroom:   selectedClassroom.ClassroomName,
							Status:      subject.Status,
							SemesterID:  selectedClassroom.SemesterID,
						}

						if _, ok := timetable[day.DayName]; !ok {
							timetable[day.DayName] = make(map[string][]models.TimetableEntry)
						}

						timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
						periodsLeft[subject.Name]--
						subjectsAssigned[day.DayName][subject.Name] = true

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

			if allAssigned && !CheckTimetableConflicts(timetable, existingTimetable) {
				return timetable
			}

			fmt.Println("Regenerating timetable due to unassigned periods or conflicts...")
		}
	} else {
		return generateRandomTimetable(days, hours, subjects, faculty, classrooms, facultySubjects, semesters)
	}
}

// generateRandomTimetable generates a random timetable
func generateRandomTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject, semesters []models.Semester) map[string]map[string][]models.TimetableEntry {
	for {
		timetable := make(FacultyBasedTimetable)
		subjectsAssigned := make(map[string]map[string]bool)
		periodsLeft := make(map[string]int)

		for _, subject := range subjects {
			periodsLeft[subject.Name] = subject.Period
		}
		for _, day := range days {
			timetable[day.DayName] = make(map[string][]models.TimetableEntry)
			subjectsAssigned[day.DayName] = make(map[string]bool)
		}

		rand.Seed(time.Now().UnixNano())

		for _, day := range days {
			for i := 0; i < 6; i++ {
				assigned := false
				for attempts := 0; attempts < 1000; attempts++ {
					var filteredSubjects []models.Subject
					for _, subject := range subjects {
						if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] {
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

					entry := models.TimetableEntry{
						DayName:     day.DayName,
						StartTime:   startTime,
						EndTime:     endTime,
						SubjectName: subject.Name,
						FacultyName: availableFaculty[facultyIndex].FacultyName,
						Classroom:   selectedClassroom.ClassroomName,
						Status:      subject.Status,
						SemesterID:  selectedClassroom.SemesterID,
					}

					if _, ok := timetable[day.DayName]; !ok {
						timetable[day.DayName] = make(map[string][]models.TimetableEntry)
					}

					timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
					periodsLeft[subject.Name]--
					subjectsAssigned[day.DayName][subject.Name] = true

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

		if allAssigned {
			return timetable
		}

		fmt.Println("Regenerating random timetable due to unassigned periods...")
	}
}

// IsPeriodAvailable checks if a period is available for the faculty
func IsPeriodAvailable(existingTimetable map[string]map[string][]models.TimetableEntry, dayName, startTime, facultyName string) bool {
	for _, entries := range existingTimetable[facultyName][dayName] {
		if entries.StartTime == startTime {
			return false
		}
	}
	return true
}

// CheckTimetableConflicts checks if the generated timetable has conflicts with the existing one
func CheckTimetableConflicts(newTimetable, existingTimetable map[string]map[string][]models.TimetableEntry) bool {
	for faculty, days := range newTimetable {
		for day, entries := range days {
			for _, entry := range entries {
				if _, exists := existingTimetable[faculty][day]; exists {
					for _, existingEntry := range existingTimetable[faculty][day] {
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
