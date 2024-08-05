// package timetables

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

// 	"ssg-portal/config"
// 	"ssg-portal/models"

// )

// // Define the FacultyBasedTimetable type
// type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

// // GenerateTimetable generates a timetable based on the provided days, hours, subjects, faculty, classrooms, and faculty-subject mappings.
// func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject) FacultyBasedTimetable {
// 	var timetable FacultyBasedTimetable
// 	var validTimetable bool

// 	// Seed the random number generator to ensure different random selections on each run.
// 	rand.Seed(time.Now().UnixNano())

// 	// Fetch existing timetable from the database.
// 	existingTimetable, err := FetchExistingTimetable()
// 	if err != nil {
// 		fmt.Println("Error fetching existing timetable:", err)
// 		return nil
// 	}

// 	for !validTimetable {
// 		// Initialize the timetable map to store classes for each faculty and day.
// 		timetable = make(FacultyBasedTimetable)
// 		// Map to track the number of periods assigned to each faculty.
// 		facultyHours := make(map[int]int)

// 		// Create maps to link subject IDs to names, faculty IDs to names, and faculty IDs to departments.
// 		subjectIDToName := make(map[int]string)
// 		facultyIDToName := make(map[int]string)
// 		facultyIDToDepartment := make(map[int]int)
// 		subjectToFaculty := make(map[int]int)

// 		// Populate the subjectIDToName map with subject names.
// 		for _, subject := range subjects {
// 			subjectIDToName[subject.ID] = subject.Name
// 		}

// 		// Populate facultyIDToName and facultyIDToDepartment maps with faculty names and departments.
// 		for _, f := range faculty {
// 			facultyIDToName[f.ID] = f.FacultyName
// 			facultyIDToDepartment[f.ID] = f.DepartmentID
// 		}

// 		// Populate subjectToFaculty map with the faculty assigned to each subject.
// 		for _, fs := range facultySubjects {
// 			subjectToFaculty[fs.SubjectID] = fs.FacultyID
// 		}

// 		// Initialize facultyAssignments to track periods assigned to each faculty for each day.
// 		facultyAssignments := make(map[int]map[string]map[string]bool)
// 		// Initialize periodsAssigned to track periods used for each day.
// 		periodsAssigned := make(map[string]map[string]bool)
// 		// Initialize assignedSubjects to track subjects assigned for each day.
// 		assignedSubjects := make(map[string]map[int]bool)

// 		// Set up the facultyAssignments, periodsAssigned, and assignedSubjects maps for each faculty and day.
// 		for _, f := range faculty {
// 			facultyAssignments[f.ID] = make(map[string]map[string]bool)
// 			for _, day := range days {
// 				facultyAssignments[f.ID][day.DayName] = make(map[string]bool)
// 				assignedSubjects[day.DayName] = make(map[int]bool) // Track subjects assigned for the day
// 			}
// 		}

// 		for _, day := range days {
// 			// Initialize used periods for the current day to ensure periods are unique.
// 			periodsAssigned[day.DayName] = make(map[string]bool)

// 			for _, hour := range hours {
// 				// Create a unique key for the current period based on start and end times.
// 				periodKey := fmt.Sprintf("%s-%s", hour.StartTime, hour.EndTime)

// 				// Skip if this period has already been assigned for the day.
// 				if periodsAssigned[day.DayName][periodKey] {
// 					continue
// 				}
// 				// Mark the period as used for the current day.
// 				periodsAssigned[day.DayName][periodKey] = true

// 				var subject models.Subject
// 				for {
// 					// Randomly select a subject.
// 					subjectIndex := rand.Intn(len(subjects))
// 					subject = subjects[subjectIndex]
// 					// Check if the subject has already been assigned for this day.
// 					if assignedSubjects[day.DayName][subject.ID] {
// 						continue
// 					}
// 					// Get the faculty assigned to this subject.
// 					facultyID, facultyAssigned := subjectToFaculty[subject.ID]
// 					// Check if the faculty is assigned and if the period has not been assigned to this faculty yet.
// 					if facultyAssigned {
// 						if !facultyAssignments[facultyID][day.DayName][periodKey] {
// 							// Check if the faculty has fewer than 6 periods assigned.
// 							if facultyHours[facultyID] < 6 {
// 								break
// 							}
// 						} else {
// 							// If the period is already assigned to the selected faculty, try another faculty.
// 							facultyID = findAlternativeFaculty(subject.ID, day.DayName, periodKey, faculty, facultyAssignments, facultyHours, subjectToFaculty)
// 							if facultyID == -1 {
// 								// No alternative faculty found; break out of the loop to try a different subject.
// 								break
// 							}
// 							break
// 						}
// 					}
// 				}

// 				facultyID, exists := subjectToFaculty[subject.ID]
// 				if !exists {
// 					// Skip if no faculty is assigned to this subject.
// 					continue
// 				}

// 				// Initialize the facultyAssignments map for the faculty and day if not already initialized.
// 				if _, ok := facultyAssignments[facultyID][day.DayName]; !ok {
// 					facultyAssignments[facultyID][day.DayName] = make(map[string]bool)
// 				}

// 				// Skip if the period is already assigned to this faculty on this day.
// 				if facultyAssignments[facultyID][day.DayName][periodKey] {
// 					continue
// 				}

// 				for _, classroom := range classrooms {
// 					// Create a new timetable entry.
// 					entry := models.TimetableEntry{
// 						DayName:      day.DayName,
// 						StartTime:    hour.StartTime,
// 						EndTime:      hour.EndTime,
// 						SubjectName:  subjectIDToName[subject.ID],
// 						FacultyName:  facultyIDToName[facultyID],
// 						Classroom:    classroom.ClassroomName,
// 						DepartmentID: facultyIDToDepartment[facultyID],
// 					}

// 					// Initialize the timetable map for the faculty if not already initialized.
// 					if _, ok := timetable[facultyIDToName[facultyID]]; !ok {
// 						timetable[facultyIDToName[facultyID]] = make(map[string][]models.TimetableEntry)
// 					}
// 					// Add the entry to the timetable for the faculty and day.
// 					timetable[facultyIDToName[facultyID]][day.DayName] = append(timetable[facultyIDToName[facultyID]][day.DayName], entry)

// 					// Increment the count of periods assigned to the faculty.
// 					facultyHours[facultyID]++
// 					// Mark the period as assigned to this faculty for the day.
// 					facultyAssignments[facultyID][day.DayName][periodKey] = true
// 					// Mark the subject as assigned for the day.
// 					assignedSubjects[day.DayName][subject.ID] = true

// 					// Break out of the loop after successfully assigning a subject to this period.
// 					break
// 				}
// 			}
// 		}

// 		// Validate the timetable to check for conflicts.
// 		validTimetable = !CheckTimetableConflicts(timetable, existingTimetable)
// 	}

// 	return timetable
// }

// // findAlternativeFaculty finds an alternative faculty member for the given subject, day, and period.
// func findAlternativeFaculty(subjectID int, dayName, periodKey string, faculty []models.Faculty, facultyAssignments map[int]map[string]map[string]bool, facultyHours map[int]int, subjectToFaculty map[int]int) int {
// 	for _, f := range faculty {
// 		if _, ok := facultyAssignments[f.ID][dayName]; !ok {
// 			facultyAssignments[f.ID][dayName] = make(map[string]bool)
// 		}
// 		if !facultyAssignments[f.ID][dayName][periodKey] && facultyHours[f.ID] < 6 {
// 			return f.ID
// 		}
// 	}
// 	return -1
// }

// // CheckTimetableConflicts checks for conflicts in the generated timetable by comparing it with the existing timetable from the database.
// func CheckTimetableConflicts(timetable FacultyBasedTimetable, existingTimetable FacultyBasedTimetable) bool {
// 	for facultyName, days := range timetable {
// 		for dayName, entries := range days {
// 			for _, entry := range entries {
// 				for _, existingEntry := range existingTimetable[facultyName][dayName] {
// 					if entry.StartTime == existingEntry.StartTime && entry.EndTime == existingEntry.EndTime && entry.Classroom == existingEntry.Classroom {
// 						return true // Conflict found.
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return false
// }
// func FetchExistingTimetable() (FacultyBasedTimetable, error) {
// 	var timetable FacultyBasedTimetable

// 	query := "SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom FROM timetable"
// 	rows, err := config.Database.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching existing timetable: %v", err)
// 	}
// 	defer rows.Close()

// 	// Initialize the timetable map
// 	timetable = make(FacultyBasedTimetable)

// 	for rows.Next() {
// 		var dayName, startTime, endTime, subjectName, facultyName, classroom string

// 		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom); err != nil {
// 			return nil, fmt.Errorf("error scanning row: %v", err)
// 		}

// 		if _, ok := timetable[facultyName]; !ok {
// 			timetable[facultyName] = make(map[string][]models.TimetableEntry)
// 		}
// 		entry := models.TimetableEntry{
// 			DayName:     dayName,
// 			StartTime:   startTime,
// 			EndTime:     endTime,
// 			SubjectName: subjectName,
// 			FacultyName: facultyName,
// 			Classroom:   classroom,
// 		}
// 		timetable[facultyName][dayName] = append(timetable[facultyName][dayName], entry)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating over rows: %v", err)
// 	}

//		return timetable, nil
//	}
package timetables

import (
	"fmt"
	"math/rand"
	"time"

	"ssg-portal/config"
	"ssg-portal/models"
)

// Define the FacultyBasedTimetable type
type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

// Define the FacultyAssignment struct
type FacultyAssignment struct {
	FacultyID int
	DayName   string
	PeriodKey string
}

// GenerateTimetable generates a timetable based on the provided days, hours, subjects, faculty, classrooms, and faculty-subject mappings.
func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject) FacultyBasedTimetable {
	var timetable FacultyBasedTimetable
	var validTimetable bool

	// Seed the random number generator to ensure different random selections on each run.
	rand.Seed(time.Now().UnixNano())

	// Fetch existing timetable from the database.
	existingTimetable, err := FetchExistingTimetable()
	if err != nil {
		fmt.Println("Error fetching existing timetable:", err)
		return nil
	}

	for !validTimetable {
		// Initialize the timetable map to store classes for each faculty and day.
		timetable = make(FacultyBasedTimetable)
		// Map to track the number of periods assigned to each faculty.
		facultyHours := make(map[int]int)

		// Create maps to link subject IDs to names, faculty IDs to names, and faculty IDs to departments.
		subjectIDToName := make(map[int]string)
		facultyIDToName := make(map[int]string)
		facultyIDToDepartment := make(map[int]int)
		subjectToFaculty := make(map[int]int)

		// Populate the subjectIDToName map with subject names.
		for _, subject := range subjects {
			subjectIDToName[subject.ID] = subject.Name
		}

		// Populate facultyIDToName and facultyIDToDepartment maps with faculty names and departments.
		for _, f := range faculty {
			facultyIDToName[f.ID] = f.FacultyName
			facultyIDToDepartment[f.ID] = f.DepartmentID
		}

		// Populate subjectToFaculty map with the faculty assigned to each subject.
		for _, fs := range facultySubjects {
			subjectToFaculty[fs.SubjectID] = fs.FacultyID
		}

		// Initialize facultyAssignments to track periods assigned to each faculty for each day.
		facultyAssignments := make(map[int]map[string]map[string]bool)
		// Initialize periodsAssigned to track periods used for each day.
		periodsAssigned := make(map[string]map[string]bool)
		// Initialize assignedSubjects to track subjects assigned for each day.
		assignedSubjects := make(map[string]map[int]bool)

		// Initialize 2D array for faculty assignment tracking.
		facultyAssignmentMatrix := make(map[int]map[string]map[string]bool)
		for _, f := range faculty {
			facultyAssignmentMatrix[f.ID] = make(map[string]map[string]bool)
			for _, day := range days {
				facultyAssignmentMatrix[f.ID][day.DayName] = make(map[string]bool)
			}
		}

		// Set up the facultyAssignments, periodsAssigned, and assignedSubjects maps for each faculty and day.
		for _, f := range faculty {
			facultyAssignments[f.ID] = make(map[string]map[string]bool)
			for _, day := range days {
				facultyAssignments[f.ID][day.DayName] = make(map[string]bool)
				assignedSubjects[day.DayName] = make(map[int]bool) // Track subjects assigned for the day
			}
		}

		for _, day := range days {
			// Initialize used periods for the current day to ensure periods are unique.
			periodsAssigned[day.DayName] = make(map[string]bool)

			for _, hour := range hours {
				// Create a unique key for the current period based on start and end times.
				periodKey := fmt.Sprintf("%s-%s", hour.StartTime, hour.EndTime)

				// Skip if this period has already been assigned for the day.
				if periodsAssigned[day.DayName][periodKey] {
					continue
				}
				// Mark the period as used for the current day.
				periodsAssigned[day.DayName][periodKey] = true

				var subject models.Subject
				for {
					// Randomly select a subject.
					subjectIndex := rand.Intn(len(subjects))
					subject = subjects[subjectIndex]
					// Check if the subject has already been assigned for this day.
					if assignedSubjects[day.DayName][subject.ID] {
						continue
					}
					// Get the faculty assigned to this subject.
					facultyID, facultyAssigned := subjectToFaculty[subject.ID]
					// Check if the faculty is assigned and if the period has not been assigned to this faculty yet.
					if facultyAssigned {
						if !facultyAssignmentMatrix[facultyID][day.DayName][periodKey] {
							// Check if the faculty has fewer than 6 periods assigned.
							if facultyHours[facultyID] < 6 {
								break
							}
						} else {
							// If the period is already assigned to the selected faculty, try another faculty.
							facultyID = findAlternativeFaculty(subject.ID, day.DayName, periodKey, faculty, facultyAssignments, facultyHours, subjectToFaculty)
							if facultyID == -1 {
								// No alternative faculty found; break out of the loop to try a different subject.
								break
							}
							break
						}
					}
				}

				facultyID, exists := subjectToFaculty[subject.ID]
				if !exists {
					// Skip if no faculty is assigned to this subject.
					continue
				}

				// Initialize the facultyAssignments map for the faculty if not already initialized.
				if _, ok := facultyAssignments[facultyID][day.DayName]; !ok {
					facultyAssignments[facultyID][day.DayName] = make(map[string]bool)
				}

				// Skip if the period is already assigned to this faculty on this day.
				if facultyAssignments[facultyID][day.DayName][periodKey] {
					continue
				}

				for _, classroom := range classrooms {
					// Create a new timetable entry.
					entry := models.TimetableEntry{
						DayName:      day.DayName,
						StartTime:    hour.StartTime,
						EndTime:      hour.EndTime,
						SubjectName:  subjectIDToName[subject.ID],
						FacultyName:  facultyIDToName[facultyID],
						Classroom:    classroom.ClassroomName,
						DepartmentID: facultyIDToDepartment[facultyID],
					}

					// Initialize the timetable map for the faculty if not already initialized.
					if _, ok := timetable[facultyIDToName[facultyID]]; !ok {
						timetable[facultyIDToName[facultyID]] = make(map[string][]models.TimetableEntry)
					}
					// Add the entry to the timetable for the faculty and day.
					timetable[facultyIDToName[facultyID]][day.DayName] = append(timetable[facultyIDToName[facultyID]][day.DayName], entry)

					// Increment the count of periods assigned to the faculty.
					facultyHours[facultyID]++
					// Mark the period as assigned to this faculty for the day.
					facultyAssignments[facultyID][day.DayName][periodKey] = true
					// Mark the subject as assigned for the day.
					assignedSubjects[day.DayName][subject.ID] = true
					// Mark the period as assigned in the matrix
					facultyAssignmentMatrix[facultyID][day.DayName][periodKey] = true

					// Break out of the loop after successfully assigning a subject to this period.
					break
				}
			}
		}

		// Validate the timetable to check for conflicts.
		validTimetable = !CheckTimetableConflicts(timetable, existingTimetable)
	}

	return timetable
}

// findAlternativeFaculty finds an alternative faculty member for the given subject, day, and period.
func findAlternativeFaculty(subjectID int, dayName, periodKey string, faculty []models.Faculty, facultyAssignments map[int]map[string]map[string]bool, facultyHours map[int]int, subjectToFaculty map[int]int) int {
	for _, f := range faculty {
		if _, ok := facultyAssignments[f.ID][dayName]; !ok {
			facultyAssignments[f.ID][dayName] = make(map[string]bool)
		}
		if !facultyAssignments[f.ID][dayName][periodKey] && facultyHours[f.ID] < 6 {
			return f.ID
		}
	}
	return -1
}

// CheckTimetableConflicts checks for conflicts in the generated timetable by comparing it with the existing timetable from the database.
func CheckTimetableConflicts(timetable FacultyBasedTimetable, existingTimetable FacultyBasedTimetable) bool {
	// Implement conflict checking logic here.
	for facultyName, days := range timetable {
		for dayName, entries := range days {
			for _, entry := range entries {
				for _, existingEntry := range existingTimetable[facultyName][dayName] {
					if entry.StartTime == existingEntry.StartTime && entry.EndTime == existingEntry.EndTime && entry.Classroom == existingEntry.Classroom {
						return true // Conflict found.
					}
				}
			}
		}
	}
	return false
}

// FetchExistingTimetable fetches the existing timetable from the database.
func FetchExistingTimetable() (FacultyBasedTimetable, error) {
	// Implement fetching logic here.

	var timetable FacultyBasedTimetable

	query := "SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom FROM timetable"
	rows, err := config.Database.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching existing timetable: %v", err)
	}
	defer rows.Close()

	// Initialize the timetable map
	timetable = make(FacultyBasedTimetable)

	for rows.Next() {
		var dayName, startTime, endTime, subjectName, facultyName, classroom string

		if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if _, ok := timetable[facultyName]; !ok {
			timetable[facultyName] = make(map[string][]models.TimetableEntry)
		}
		entry := models.TimetableEntry{
			DayName:     dayName,
			StartTime:   startTime,
			EndTime:     endTime,
			SubjectName: subjectName,
			FacultyName: facultyName,
			Classroom:   classroom,
		}

		timetable[facultyName][dayName] = append(timetable[facultyName][dayName], entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return nil, nil
}
