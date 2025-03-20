package subject

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetSubjects(departmentID, semesterID, academicYearID, sectionID int) ([]models.Subject, error) {
	var subjects []models.Subject

	query := `
		SELECT 
    s.id AS subject_id,
    s.name,
    fs.status,
    fs.periods,
    s.course_code
FROM 
    faculty_subjects fs
JOIN 
    subjects s ON fs.subject_id = s.id
WHERE 
    fs.department_id = ? AND 
    fs.semester_id = ? AND 
    fs.academic_year_id = ? AND 
    fs.section_id = ?;

	`

	rows, err := config.Database.Query(query, departmentID, semesterID, academicYearID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying subjects: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Status, &subject.Period, &subject.CourseCode); err != nil {
			return nil, fmt.Errorf("error scanning subject: %v", err)
		}
		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return subjects, nil
}
