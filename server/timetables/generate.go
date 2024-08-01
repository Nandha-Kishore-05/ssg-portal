package timetables

import (
	"fmt"
	"math/rand"
	"time"

	"ssg-portal/models"
)

func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject) map[string][]models.Class {
	timetable := make(map[string][]models.Class)
	facultyHours := make(map[int]int)

	rand.Seed(time.Now().UnixNano())

	subjectIDToName := make(map[int]string)
	for _, subject := range subjects {
		subjectIDToName[subject.ID] = subject.Name
	}

	facultyIDToName := make(map[int]string)
	facultyIDToDepartment := make(map[int]int)
	for _, f := range faculty {
		facultyIDToName[f.ID] = f.FacultyName
		facultyIDToDepartment[f.ID] = f.DepartmentID
	}

	subjectToFaculty := make(map[int]int)
	for _, fs := range facultySubjects {
		subjectToFaculty[fs.SubjectID] = fs.FacultyID
	}

	facultyAssignments := make(map[int]map[string]map[string]bool)

	for _, day := range days {
		usedPeriods := make(map[string]struct{})
		assignedSubjects := make(map[string]struct{})

		for _, hour := range hours {

			periodKey := fmt.Sprintf("%s-%s", hour.StartTime, hour.EndTime)
			if _, exists := usedPeriods[periodKey]; exists {
				continue
			}
			usedPeriods[periodKey] = struct{}{}

			var subject models.Subject
			for {
				subjectIndex := rand.Intn(len(subjects))
				subject = subjects[subjectIndex]
				if _, assigned := assignedSubjects[subject.Name]; !assigned {
					break
				}
			}
			assignedSubjects[subject.Name] = struct{}{}

			facultyID, exists := subjectToFaculty[subject.ID]
			if !exists {
				continue
			}

			if _, ok := facultyAssignments[facultyID]; !ok {
				facultyAssignments[facultyID] = make(map[string]map[string]bool)
			}
			if _, ok := facultyAssignments[facultyID][day.DayName]; !ok {
				facultyAssignments[facultyID][day.DayName] = make(map[string]bool)
			}

			if facultyHours[facultyID] >= 36 {
				continue
			}

			conflict := false
			for assignedPeriod := range facultyAssignments[facultyID][day.DayName] {
				if assignedPeriod == periodKey {
					conflict = true
					break
				}
			}
			if conflict {
				continue
			}

			for _, classroom := range classrooms {
				class := models.Class{
					Day:         day.DayName,
					StartTime:   hour.StartTime,
					EndTime:     hour.EndTime,
					SubjectName: subject.Name,
					FacultyName: facultyIDToName[facultyID],
					Classroom:   classroom.ClassroomName,
				}
				timetable[day.DayName] = append(timetable[day.DayName], class)
				facultyHours[facultyID]++
				facultyAssignments[facultyID][day.DayName][periodKey] = true
				break
			}
		}
	}

	return timetable
}
