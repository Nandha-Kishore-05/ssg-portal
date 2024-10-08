import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import axios from 'axios'; // Import axios for API call
import AppLayout from '../../layout/layout';
import CustomCard from '../../components/card';
import SavedTimetable from '../workload/timetable';

import './mastertimetable.css';

function MasterSemester() {
  const location = useLocation();

  const { yearId, deptId } = location.state || {};

  const [isOpen, setIsOpen] = useState(false);
  const [semOptions, setSemOptions] = useState([]); // State for storing semester options
  const [selectedSemester, setSelectedSemester] = useState(null); // State for the selected semester

  // Fetch semester options from the API
  useEffect(() => {
    axios.get('http://localhost:8080/timetable/semoptions')
      .then(response => {
        setSemOptions(response.data); // Update state with API response
      })
      .catch(error => {
        console.error('Error fetching semester options:', error);
      });
  }, []);

  const handleCardClick = (semester) => {
    setSelectedSemester(semester); // Set the selected semester
    setIsOpen(true); // Show the SavedTimetable component
  };

  // Conditionally render SavedTimetable
  if (yearId && deptId && selectedSemester && isOpen) {
    return (
      <AppLayout
        rId={11}
        title="Saved Timetable"
        body={
          <SavedTimetable 
            setIsOpen={setIsOpen} 
            departmentID={deptId} 
            semesterID={selectedSemester.value} // Use the value (ID) from the API
            academicYearID={yearId} 
          />
        }
      />
    );
  }

  // Default rendering
  return (
    <AppLayout
      rId={11}
      title="Master Semester"
      body={
        <div className="cards-container">
          {semOptions.map((semesterObj) => (
            <CustomCard 
              key={semesterObj.value} // Use the value as the key
              year={semesterObj.label} // Display the label (semester name)
              title={`Semester ${semesterObj.label}`} 
              onCardClick={() => handleCardClick(semesterObj)} 
            />
          ))}
        </div>
      }
    />
  );
}

export default MasterSemester;
