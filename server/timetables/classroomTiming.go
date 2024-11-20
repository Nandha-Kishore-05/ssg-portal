package timetables

import (
	"net/http"
	"ssg-portal/config"

	"github.com/gin-gonic/gin"
)

func GetAvailableTimingsForFacultyAndClassroom(c *gin.Context) {
	facultyName := c.Param("facultyName")
	day := c.Param("day")
	classroomName := c.Param("classroomName")
	academicYearID := c.Param("academicYearID") // Retrieve academic year ID from parameters

	availableTimings, err := FacultyAndClassroomAvailableTimings(facultyName, day, classroomName, academicYearID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availableTimings)
}

func FacultyAndClassroomAvailableTimings(facultyName, day, classroomName, academicYearID string) ([]map[string]string, error) {
	var availableTimings []map[string]string

	query := `
	SELECT available_slots.start_time, available_slots.end_time
	FROM (
		SELECT t1.start_time, t1.end_time
		FROM hours AS t1
		WHERE NOT EXISTS (
			SELECT 1
			FROM (
				SELECT faculty_name, day_name, start_time, end_time 
				FROM timetable
				WHERE academic_year = ?
				UNION ALL
				SELECT faculty_name, day_name, start_time, end_time 
				FROM manual_timetable
				WHERE academic_year = ?
				UNION ALL
				SELECT faculty_name, day_name, start_time, end_time 
				FROM timetable_skips
				WHERE academic_year = ?
			) AS combined_faculty
			WHERE combined_faculty.faculty_name = ?
			  AND combined_faculty.day_name = ?
			  AND (
				  t1.start_time < combined_faculty.end_time
				  AND t1.end_time > combined_faculty.start_time
			  )
		)
		AND NOT EXISTS (
			SELECT 1
			FROM (
				SELECT classroom AS classroom, day_name, start_time, end_time 
				FROM timetable
				WHERE academic_year = ?
				UNION ALL
				SELECT classroom AS classroom, day_name, start_time, end_time 
				FROM manual_timetable
				WHERE academic_year = ?
				UNION ALL
				SELECT classroom AS classroom, day_name, start_time, end_time 
				FROM timetable_skips
				WHERE academic_year = ?
			) AS combined_classroom
			WHERE combined_classroom.classroom = ?
			  AND combined_classroom.day_name = ?
			  AND (
				  t1.start_time < combined_classroom.end_time
				  AND t1.end_time > combined_classroom.start_time
			  )
		)
	) AS available_slots
	ORDER BY available_slots.start_time;
	`

	rows, err := config.Database.Query(
		query,
		academicYearID, academicYearID, academicYearID, // For combined_faculty
		facultyName, day,
		academicYearID, academicYearID, academicYearID, // For combined_classroom
		classroomName, day,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var startTime, endTime string
		if err := rows.Scan(&startTime, &endTime); err != nil {
			return nil, err
		}
		availableTimings = append(availableTimings, map[string]string{
			"day_name":   day,
			"start_time": startTime,
			"end_time":   endTime,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return availableTimings, nil
}
