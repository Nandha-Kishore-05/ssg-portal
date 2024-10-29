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

func FetchTimetableSkips(departmentID int, semesterID int, academicYearID int, sectionID int) (map[string]map[string]models.TimetableEntry, error) {
    skipEntries := make(map[string]map[string]models.TimetableEntry)

    query := `
        SELECT day_name, start_time, end_time, subject_name, faculty_name, semester_id, department_id ,classroom,status,academic_year,course_code,section_id
        FROM timetable_skips 
        WHERE department_id = ? AND semester_id = ? AND  academic_year = ? AND section_id = ?`

    rows, err := config.Database.Query(query, departmentID, semesterID, &academicYearID, &sectionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        var dayName, startTime, endTime, subjectName, facultyName, classroom, courseCode string
        var status, academicYearID,sectionID int

        if err := rows.Scan(&dayName, &startTime, &endTime, &subjectName, &facultyName, &semesterID, &departmentID, &classroom, &status, &academicYearID, &courseCode,&sectionID); err != nil {
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
			SectionID: sectionID,
        }

        skipEntries[dayName][startTime] = entry
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
func GenerateTimetable(days []models.Day, hours []models.Hour, subjects []models.Subject, faculty []models.Faculty, classrooms []models.Classroom, facultySubjects []models.FacultySubject, semesters []models.Semester, section []models.Section, academicYear []models.AcademicYear, departmentID int, semesterID int, academicYearID int, sectionID int) map[string]map[string][]models.TimetableEntry {

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

    for {
        timetable := make(map[string]map[string][]models.TimetableEntry)
        subjectsAssigned := make(map[string]map[string]bool)
        periodsLeft := make(map[string]int)
        status0Assignments := make(map[string]map[string]bool)
        facultyAssignments := make(map[string]map[string]string)
        facultyDailyCount := make(map[string]map[string]int)
        labAssigned := make(map[string]bool)

        // Initialize periods left for each subject
        for _, subject := range subjects {
            periodsLeft[subject.Name] = subject.Period
            if subject.Status == 0 {
                status0Assignments[subject.Name] = make(map[string]bool)
            }
        }

        // Incorporate manual timetable into the generated timetable
        for _, day := range days {
            timetable[day.DayName] = make(map[string][]models.TimetableEntry)
            subjectsAssigned[day.DayName] = make(map[string]bool)
            facultyAssignments[day.DayName] = make(map[string]string)
            facultyDailyCount[day.DayName] = make(map[string]int)
            labAssigned[day.DayName] = false

            // If skips exist for the day, add them
            if skips, ok := skipTimetable[day.DayName]; ok {
                for startTime, entry := range skips {
                    timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
                    subjectsAssigned[day.DayName][entry.SubjectName] = true
                    // periodsLeft[entry.SubjectName]--
                    if entry.Status == 0 {
                        labAssigned[day.DayName] = true
                    }
                }
            }

            // Add manual timetable entries for the current day
            if manualEntries, ok := manualTimetable[day.DayName]; ok {
                for startTime, entries := range manualEntries {
                    for _, entry := range entries {
                        timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
                        subjectsAssigned[day.DayName][entry.SubjectName] = true
                        // periodsLeft[entry.SubjectName]--
                        facultyDailyCount[day.DayName][entry.FacultyName]++
                        if entry.Status == 0 {
                            labAssigned[day.DayName] = true
                        }
                    }
                }
            }
        }

        // Continue with the automatic generation for remaining periods...
        rand.Seed(time.Now().UnixNano())

        for _, day := range days {
            for i := 0; i < len(hours); i++ {
                assigned := false
                for attempts := 0; attempts < 1000; attempts++ {
                    var filteredSubjects []models.Subject
                    for _, subject := range subjects {
                        if periodsLeft[subject.Name] > 0 && (!subjectsAssigned[day.DayName][subject.Name] || (subject.Status == 0 && len(status0Assignments[subject.Name]) == 0)) {

                            if subject.Status == 0 && labAssigned[day.DayName] {
                                continue
                            }
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

                    if !Available(existingTimetable, day.DayName, startTime, selectedFaculty.FacultyName) {
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

                    timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
                    periodsLeft[subject.Name]--
                    facultyDailyCount[day.DayName][selectedFaculty.FacultyName]++
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

        if allAssigned && periodsFilled && !CheckTimetableConflicts(timetable, existingTimetable) {
            return timetable
        }
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
    // Add db as a parameter for fetching manual timetable
) FacultyBasedTimetable {
    log.Println("Academic Year ID:", academicYearID)

    // Fetch manual timetable
    manualTimetable, err := FetchManualTimetable(departmentID, semesterID, academicYearID, sectionID)
    if err != nil {
        fmt.Println("Error fetching manual timetable:", err)
        return nil
    }

    skipTimetable, err := FetchTimetableSkips(departmentID, semesterID, academicYearID, sectionID)
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
        facultyAssignments := make(map[string]map[string]int)
        facultyDailyCount := make(map[string]map[string]int)

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

        // Pre-fill timetable with manual timetable data
        for _, day := range days {
            timetable[day.DayName] = make(map[string][]models.TimetableEntry)
            subjectsAssigned[day.DayName] = make(map[string]bool)
            labSubjectAssigned[day.DayName] = false
            facultyAssignments[day.DayName] = make(map[string]int)
            facultyDailyCount[day.DayName] = make(map[string]int)

            // Check for manual timetable entries
            if manualEntries, ok := manualTimetable[day.DayName]; ok {
                for startTime, entries := range manualEntries {
                    for _, entry := range entries {
                        timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
                        subjectsAssigned[day.DayName][entry.SubjectName] = true
                        // periodsLeft[entry.SubjectName]-- // Decrease the period count as it's already assigned
                        facultyDailyCount[day.DayName][entry.FacultyName]++
                    }
                }
            }

            // Check for timetable skips
            if skips, ok := skipTimetable[day.DayName]; ok {
                for startTime, entry := range skips {
                    timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry)
                    subjectsAssigned[day.DayName][entry.SubjectName] = true
                    periodsLeft[entry.SubjectName]-- // Decrease the period count as it's already assigned
                }
            }
        }

        rand.Seed(time.Now().UnixNano())

        // Generate the timetable with lab subjects first
        for _, day := range days {
            for i := 0; i < len(hours); i++ {
                // Skip periods already assigned by the manual timetable or skips
                startTime := hours[i].StartTime
                if len(timetable[day.DayName][startTime]) > 0 {
                    continue
                }

                for attempts := 0; attempts < maxAttempts; attempts++ {
                    // Filtering lab subjects that haven't been assigned yet
                    var filteredLabSubjects []models.Subject
                    for _, subject := range labSubjects {
                        if periodsLeft[subject.Name] > 0 && !subjectsAssigned[day.DayName][subject.Name] && !labSubjectAssigned[day.DayName] {
                            filteredLabSubjects = append(filteredLabSubjects, subject)
                        }
                    }

                    if len(filteredLabSubjects) == 0 {
                        break
                    }

                    // Assigning the lab subject
                    subjectIndex := rand.Intn(len(filteredLabSubjects))
                    subject := filteredLabSubjects[subjectIndex]

                    if subject.Status == 0 && i < len(hours)-1 {
                        nextStartTime := hours[i+1].StartTime
                        if IsPeriodAvailable(timetable, day.DayName, nextStartTime, "") {
                            facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
                            if facultyName == "" {
                                fmt.Println("Error: No faculty available for lab subject", subject.Name)
                                return nil
                            }
                            if facultyDailyCount[day.DayName][facultyName] >= 2 {
                                continue
                            }

                            classroomName := selectRandomClassroom(classrooms)
                            if classroomName == "" {
                                fmt.Println("Error: No classroom found for lab subject", subject.Name)
                                return nil
                            }

                            // Create lab entries for 2 periods
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

                            // Add entries to timetable
                            timetable[day.DayName][startTime] = append(timetable[day.DayName][startTime], entry1)
                            timetable[day.DayName][nextStartTime] = append(timetable[day.DayName][nextStartTime], entry2)

                            // Mark subject, faculty, and classroom as assigned
                            periodsLeft[subject.Name] -= 2
                            subjectsAssigned[day.DayName][subject.Name] = true
                            labSubjectAssigned[day.DayName] = true
                            facultyAssignments[day.DayName][facultyName]++
                            facultyDailyCount[day.DayName][facultyName] += 2
                            break
                        }
                    }
                }
            }
        }

        // Generate the timetable for non-lab subjects
        for _, day := range days {
            for i := 0; i < len(hours); i++ {
                startTime := hours[i].StartTime
                if len(timetable[day.DayName][startTime]) > 0 {
                    continue
                }

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

                    facultyName := selectRandomFaculty(faculty, subject, facultySubjects, departmentID, semesterID, academicYearID, sectionID)
                    if facultyName == "" {
                        fmt.Println("Error: No faculty available for non-lab subject", subject.Name)
                        return nil
                    }
                    if facultyDailyCount[day.DayName][facultyName] >= 1 {
                        continue
                    }
                    classroomName := selectRandomClassroom(classrooms)
                    if classroomName == "" {
                        fmt.Println("Error: No classroom found for non-lab subject", subject.Name)
                        return nil
                    }

                    // Add the timetable entry for the non-lab subject
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
                    facultyAssignments[day.DayName][facultyName]++
                    facultyDailyCount[day.DayName][facultyName]++
                    break
                }
            }
        }

        return timetable
    }

    // Generate timetable until all periods are filled
    for {
        timetable := generate()
        allPeriodsFilled := true
        for _, day := range days {
            for _, hour := range hours {
                startTime := hour.StartTime
                if len(timetable[day.DayName][startTime]) == 0 {
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