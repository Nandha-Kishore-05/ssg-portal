import React, { useState, useEffect } from 'react';
import axios from 'axios';

import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';
import SavedTimetable from './timetable';

const SaveTimetable = () => {
 
  const [department, setDepartment] = useState(null);
  const [semester, setSemester] = useState(null);
  const [deptOptions, setDeptOptions] = useState([]);
  const [semOptions, setSemOptions] = useState([]);
  const [isOpen,setIsOpen] = useState(false)
  useEffect(() => {
    axios.get('http://localhost:8080/timetable/options')
      .then(response => {
        setDeptOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching faculty names:', error);
      });
  }, []);

  useEffect(() => {
    axios.get('http://localhost:8080/timetable/semoptions')
      .then(response => {
        setSemOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching faculty names:', error);
      });
  }, []);

  const handleViewTimetable = () => {
    if (department && semester) {
      setIsOpen(true);
    } else {
      console.error('Please select both department and semester');
    }
  };
    return (
        <AppLayout
          rId={3}
          title="Venue Table"
          body={
            <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
                          <div style={{display:'flex',flexDirection:'row',columnGap:10,alignItems:"center"}}>
            <CustomSelect
            placeholder="DEPARTMENT"
            value={department}
            onChange={setDepartment}
            options={deptOptions}
          />
      
          <CustomSelect
            placeholder="SEMESTER"
            value={semester}
            onChange={setSemester}
            options={semOptions}
           
          />
         
            <CustomButton
              width="150"
              label="View Timetable"
              onClick={handleViewTimetable}
            />
      
          </div>
          { (department && semester && isOpen) && 
          <SavedTimetable departmentID={department.value} semesterID = {semester.value} />
             
          }
                </div>
          }
          />

        );
};

export default SaveTimetable;