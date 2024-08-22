import React, { useState, useEffect } from 'react';
import axios from 'axios';

import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';
import Timetable from './workload';

const GenerateTimetable = () => {
  const [department, setDepartment] = useState(null);
  const [deptOptions, setDeptOptions] = useState([]);
  const [semester, setSemester] = useState(null);
  const [semOptions, setSemOptions] = useState([]);
  const [isOpen, setIsOpen] = useState(false);
  const [viewedDepartment, setViewedDepartment] = useState(null);
  const [viewedSemester, setViewedSemester] = useState(null);

  useEffect(() => {
    axios.get('http://localhost:8080/timetable/options')
      .then(response => {
        setDeptOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching department options:', error);
      });
  }, []);

  useEffect(() => {
    axios.get('http://localhost:8080/timetable/semoptions')
      .then(response => {
        setSemOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching semester options:', error);
      });
  }, []);

  const handleViewTimetable = () => {
    if (department && semester) {
      setViewedDepartment(department.value);
      setViewedSemester(semester.value);
      setIsOpen(true);
    } else {
      console.error('Please select both department and semester');
    }
  };

  return (
    <AppLayout
      rId={2}
      title="Time Table"
      body={
        <div style={{ backgroundColor: "white", padding: 17, marginTop: 20, borderRadius: "10px" }}>
          <div style={{ display: 'flex', flexDirection: 'row', columnGap: 10, alignItems: "center" }}>
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
              label="Generate Timetable"
              onClick={handleViewTimetable}
              backgroundColor="#0878d3"
            />
          </div>

          {(viewedDepartment && viewedSemester && isOpen) && 
            <Timetable departmentID={viewedDepartment} semesterID={viewedSemester} />
          }
        </div>
      }
    />
  );
};

export default GenerateTimetable;
