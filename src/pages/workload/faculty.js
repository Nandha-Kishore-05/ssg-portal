

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';
import FacultyTimetable from './facultytable';

const FacTimetable = () => {

  const [facultyOptions, setFacultyOptions] = useState([]);
  const [selectedFaculty, setSelectedFaculty] = useState("");
  const [isOpen,setIsOpen] = useState(false)
  useEffect(() => {
    axios.get('http://localhost:8080/timetable/facultyOptions')
      .then(response => {
        setFacultyOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching faculty names:', error);
      });
  }, []);

  const handleViewTimetable = () => {
    if (selectedFaculty) {
      setIsOpen(true);
    } else {
      console.error('Please select a faculty');
    }
  };

  return (
    <AppLayout
      rId={4}
      title="Faculty Table"
      body={
        <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
                         <div style={{display:'flex',flexDirection:'row',columnGap:10,alignItems:"center"}}>
          <CustomSelect
            placeholder="Faculty Name"
            value={selectedFaculty}
            onChange={setSelectedFaculty}
            options={facultyOptions}
          />
          
            <CustomButton
              width="150"
              label="View Timetable"
              onClick={handleViewTimetable}
            />
            </div>
       
          { (selectedFaculty && isOpen) && 
          <FacultyTimetable facultyName ={selectedFaculty.value} />
             
          }
        </div>
      }
    />
  );
};

export default FacTimetable;
