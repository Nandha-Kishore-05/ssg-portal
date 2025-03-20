package section

import (
	"fmt"

	"ssg-portal/config"
	"ssg-portal/models"
)

func GetSectionDetails(sectionID int) ([]models.Section, error) {
	var  sections []models.Section
	rows, err := config.Database.Query("SELECT id, section_name FROM master_section WHERE id = ?",sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying hours: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var section models.Section
		if err := rows.Scan(&section.ID, &section.SectionName); err != nil {
			return nil, fmt.Errorf("error scanning hour: %v", err)
		}
		sections = append(sections, section)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return sections, nil
}
