package timetables

import (
	"fmt"
	"math/rand"
	"time"

	"ssg-portal/config"
	"ssg-portal/models"
)

type FacultyBasedTimetable map[string]map[string][]models.TimetableEntry

// Function to fetch existing timetable data from the database
func FetchExistingTimetable() (map[string]map[string][]models.TimetableEntry, error) {
	var existingTimetable map[string]map[string][]models.TimetableEntry
	rows, err := config.Database.Query("SELECT day_name, start_time, end_time, subject_name, faculty_name, classroom FROM timetable")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existingTimetable = make(map[string]map[string][]models.TimetableEntry)
	for rows.Next() {
		var dayName, startTime, endTime, subjectName, facultyName, classroom string
		err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &classroom)
		if err != nil {
			return nil, err
		}
		if _, exists := existingTimetable[dayName]; !exists {
			existingTimetable[dayName] = make(map[string][]models.TimetableEntry)
		}
		periodKey := fmt.Sprintf("%s-%s", startTime, endTime)
		entry := models.TimetableEntry{
			DayName:     dayName,
			StartTime:   startTime,
			EndTime:     endTime,
			SubjectName: subjectName,
			FacultyName: facultyName,
			Classroom:   classroom,
		}
		existingTimetable[dayName][periodKey] = append(existingTimetable[dayName][periodKey], entry)
	}

	return existingTimetable, nil
}

// GenerateTimetable generates a timetable based on existing entries and constraints
func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject) FacultyBasedTimetable {
	var timetable FacultyBasedTimetable
	var validTimetable bool

	// Seed the random number generator for randomness
	rand.Seed(time.Now().UnixNano())

	// Fetch existing timetable data
	existingTimetable, err := FetchExistingTimetable()
	if err != nil {
		fmt.Println("Error fetching existing timetable:", err)
		return nil // Return nil if there is an error fetching existing timetable
	}

	// If there is existing timetable data, try to resolve conflicts
	if len(existingTimetable) > 0 {
		for !validTimetable {
			timetable = make(FacultyBasedTimetable)
			facultyHours := make(map[int]int)

			// Maps for lookup
			subjectIDToName := make(map[int]string)
			facultyIDToName := make(map[int]string)
			facultyIDToDepartment := make(map[int]int)
			subjectToFaculty := make(map[int]int)

			// Populate maps with subject names and faculty information
			for _, subject := range subjects {
				subjectIDToName[subject.ID] = subject.Name
			}

			for _, f := range faculty {
				facultyIDToName[f.ID] = f.FacultyName
				facultyIDToDepartment[f.ID] = f.DepartmentID
			}

			for _, fs := range facultySubjects {
				subjectToFaculty[fs.SubjectID] = fs.FacultyID
			}

			// Initialize maps to track assignments and periods
			facultyAssignments := make(map[string]map[string]bool)
			periodsAssigned := make(map[string]map[string]bool)
			assignedSubjects := make(map[string]map[int]bool)
			classroomPeriods := make(map[string]map[string]bool)

			for _, day := range days {
				facultyAssignments[day.DayName] = make(map[string]bool)
				assignedSubjects[day.DayName] = make(map[int]bool)
				classroomPeriods[day.DayName] = make(map[string]bool)
				periodsAssigned[day.DayName] = make(map[string]bool)
			}

			// Loop through each day and hour to assign subjects and classrooms
			for _, day := range days {
				for _, hour := range hours {
					periodKey := fmt.Sprintf("%s-%s", hour.StartTime, hour.EndTime)
					if periodsAssigned[day.DayName][periodKey] {
						continue // Skip if the period is already assigned
					}
					periodsAssigned[day.DayName][periodKey] = true

					var subject models.Subject
					for {
						subjectIndex := rand.Intn(len(subjects))
						subject = subjects[subjectIndex]
						if assignedSubjects[day.DayName][subject.ID] {
							continue // Skip if the subject is already assigned for the day
						}
						facultyID, facultyAssigned := subjectToFaculty[subject.ID]
						if facultyAssigned {
							if facultyHours[facultyID] < 6 {
								break // Break loop if faculty hours are less than 6
							}
						}
					}

					facultyID, exists := subjectToFaculty[subject.ID]
					if !exists {
						continue // Skip if no faculty assigned for the subject
					}

					if _, ok := facultyAssignments[day.DayName]; !ok {
						facultyAssignments[day.DayName] = make(map[string]bool)
					}

					if facultyAssignments[day.DayName][periodKey] {
						continue // Skip if the faculty is already assigned for the period
					}

					assigned := false
					for _, classroom := range classrooms {
						if classroomPeriods[day.DayName][periodKey] {
							continue // Skip if the period is already assigned to a classroom
						}

						if !IsPeriodAvailable(existingTimetable, day.DayName, periodKey, facultyIDToName[facultyID]) {
							continue // Skip if the period is not available for the faculty
						}

						entry := models.TimetableEntry{
							DayName:     day.DayName,
							StartTime:   hour.StartTime,
							EndTime:     hour.EndTime,
							SubjectName: subjectIDToName[subject.ID],
							FacultyName: facultyIDToName[facultyID],
							Classroom:   classroom.ClassroomName,
						}

						if _, ok := timetable[facultyIDToName[facultyID]]; !ok {
							timetable[facultyIDToName[facultyID]] = make(map[string][]models.TimetableEntry)
						}
						timetable[facultyIDToName[facultyID]][day.DayName] = append(timetable[facultyIDToName[facultyID]][day.DayName], entry)

						facultyHours[facultyID]++
						facultyAssignments[day.DayName][periodKey] = true
						assignedSubjects[day.DayName][subject.ID] = true
						classroomPeriods[day.DayName][periodKey] = true

						assigned = true
						break // Break loop once a classroom is assigned
					}

					if !assigned {
						continue // Skip to the next period if no classroom was assigned
					}
				}
			}

			// Check if the generated timetable conflicts with existing timetable
			validTimetable = !CheckTimetableConflicts(timetable, existingTimetable)
		}
	} else {
		// If no existing timetable, generate a new timetable
		for !validTimetable {
			timetable = make(FacultyBasedTimetable)
			facultyHours := make(map[int]int)

			// Maps for lookup
			subjectIDToName := make(map[int]string)
			facultyIDToName := make(map[int]string)
			facultyIDToDepartment := make(map[int]int)
			subjectToFaculty := make(map[int]int)

			// Populate maps with subject names and faculty information
			for _, subject := range subjects {
				subjectIDToName[subject.ID] = subject.Name
			}

			for _, f := range faculty {
				facultyIDToName[f.ID] = f.FacultyName
				facultyIDToDepartment[f.ID] = f.DepartmentID
			}

			for _, fs := range facultySubjects {
				subjectToFaculty[fs.SubjectID] = fs.FacultyID
			}

			// Initialize maps to track assignments and periods
			facultyAssignments := make(map[string]map[string]bool)
			periodsAssigned := make(map[string]map[string]bool)
			assignedSubjects := make(map[string]map[int]bool)
			classroomPeriods := make(map[string]map[string]bool)

			for _, day := range days {
				facultyAssignments[day.DayName] = make(map[string]bool)
				assignedSubjects[day.DayName] = make(map[int]bool)
				classroomPeriods[day.DayName] = make(map[string]bool)
				periodsAssigned[day.DayName] = make(map[string]bool)
			}

			// Loop through each day and hour to assign subjects and classrooms
			for _, day := range days {
				for _, hour := range hours {
					periodKey := fmt.Sprintf("%s-%s", hour.StartTime, hour.EndTime)
					if periodsAssigned[day.DayName][periodKey] {
						continue // Skip if the period is already assigned
					}
					periodsAssigned[day.DayName][periodKey] = true

					var subject models.Subject
					for {
						subjectIndex := rand.Intn(len(subjects))
						subject = subjects[subjectIndex]
						if assignedSubjects[day.DayName][subject.ID] {
							continue // Skip if the subject is already assigned for the day
						}
						facultyID, facultyAssigned := subjectToFaculty[subject.ID]
						if facultyAssigned {
							if facultyHours[facultyID] < 6 {
								break // Break loop if faculty hours are less than 6
							}
						}
					}

					facultyID, exists := subjectToFaculty[subject.ID]
					if !exists {
						continue // Skip if no faculty assigned for the subject
					}

					if _, ok := facultyAssignments[day.DayName]; !ok {
						facultyAssignments[day.DayName] = make(map[string]bool)
					}

					if facultyAssignments[day.DayName][periodKey] {
						continue // Skip if the faculty is already assigned for the period
					}

					assigned := false
					for _, classroom := range classrooms {
						if classroomPeriods[day.DayName][periodKey] {
							continue // Skip if the period is already assigned to a classroom
						}

						entry := models.TimetableEntry{
							DayName:     day.DayName,
							StartTime:   hour.StartTime,
							EndTime:     hour.EndTime,
							SubjectName: subjectIDToName[subject.ID],
							FacultyName: facultyIDToName[facultyID],
							Classroom:   classroom.ClassroomName,
						}

						if _, ok := timetable[facultyIDToName[facultyID]]; !ok {
							timetable[facultyIDToName[facultyID]] = make(map[string][]models.TimetableEntry)
						}
						timetable[facultyIDToName[facultyID]][day.DayName] = append(timetable[facultyIDToName[facultyID]][day.DayName], entry)

						facultyHours[facultyID]++
						facultyAssignments[day.DayName][periodKey] = true
						assignedSubjects[day.DayName][subject.ID] = true
						classroomPeriods[day.DayName][periodKey] = true

						assigned = true
						break // Break loop once a classroom is assigned
					}

					if !assigned {
						continue // Skip to the next period if no classroom was assigned
					}
				}
			}

			// Check if the generated timetable is valid and doesn't conflict with existing timetable
			validTimetable = !CheckTimetableConflicts(timetable, existingTimetable)
		}
	}

	return timetable // Return the generated timetable
}

// IsPeriodAvailable checks if the period is available for the faculty based on the existing timetable
func IsPeriodAvailable(existingTimetable map[string]map[string][]models.TimetableEntry, day, period, facultyName string) bool {
	for _, entries := range existingTimetable[day] {
		for _, entry := range entries {
			if entry.FacultyName == facultyName && entry.StartTime == period {
				return false // Conflict found if the same faculty is already assigned to the period
			}
		}
	}
	return true // Return true if no conflict is found
}

// Function to check for timetable conflicts

func CheckTimetableConflicts(timetable FacultyBasedTimetable, existingTimetable map[string]map[string][]models.TimetableEntry) bool {
	for dayName, periods := range existingTimetable {
		for _, entries := range periods {
			for _, entry := range entries {
				if facultyEntries, exists := timetable[entry.FacultyName]; exists {
					if dayEntries, exists := facultyEntries[dayName]; exists {
						for _, newEntry := range dayEntries {
							if newEntry.StartTime == entry.StartTime && newEntry.EndTime == entry.EndTime {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}
