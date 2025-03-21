package  facultyDetails

import (
	"net/http"
	"ssg-portal/config"
	"github.com/gin-gonic/gin"
)

func GetAvailableTimingsForFaculty(c *gin.Context) {
	facultyName := c.Param("facultyName")
	day := c.Param("day")

	availableTimings, err := FacultyAvailableTimings(facultyName, day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availableTimings)
}

func FacultyAvailableTimings(facultyName, day string) ([]map[string]string, error) {
	var availableTimings []map[string]string

	query := `
	SELECT available_slots.start_time, available_slots.end_time
	FROM (
		SELECT t1.start_time, t1.end_time
		FROM hours AS t1
		WHERE NOT EXISTS (
			SELECT 1
			FROM timetable AS t2
			WHERE t2.faculty_name = ?
			  AND t2.day_name = ?
			  AND (
					(t1.start_time < t2.end_time AND t1.end_time > t2.start_time)
				)
		)
	) AS available_slots
	ORDER BY available_slots.start_time;
	`

	rows, err := config.Database.Query(query, facultyName, day)
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
			"day_name": day,
			"start_time": startTime,
			"end_time":   endTime,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return availableTimings, nil
}
