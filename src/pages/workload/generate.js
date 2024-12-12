

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';
import Timetable from './workload';
import InputBox from '../../components/input';

const GenerateTimetable = () => {
  const [department, setDepartment] = useState(null);
  const [deptOptions, setDeptOptions] = useState([]);
  const [semester, setSemester] = useState(null);
  const [semOptions, setSemOptions] = useState([]);
  const [filteredSemOptions, setFilteredSemOptions] = useState([]);
  const [isOpen, setIsOpen] = useState(false);
  const [viewedDepartment, setViewedDepartment] = useState(null);
  const [viewedSemester, setViewedSemester] = useState(null);
  const [viewedAcademic, setViewedAcademic] = useState(null);
  const [viewedSection, setViewedSection] = useState(null);
  const [academicYear, setAcademicYear] = useState(null);
  const [academicsOptions, setAcademicsOptions] = useState([]);
  const [section, setSection] = useState(null);
  const [sectionOptions, setSectionOptions] = useState([]);
  const [vieweddaysCount, setVieweddaysCount] = useState(null);
  const [daysCount, setdaysCount] = useState('');
  const [startDate, setstartDate] = useState(null);
  const [startDateOptions, setstartDateOptions] = useState([]);
  const [viewedstartDate, setViewedstartDate] = useState(null);
  const [endDate, setendDate] = useState(null);
  const [endDateOptions, setendDateOptions] = useState([]);
  const [viewedendDate, setViewedendDate] = useState(null);
  // Fetch department options

  useEffect(() => {
    const fetchStartDate = async () => {
      try {
        const response = await axios.get('http://localhost:8080/workingDayoptions');
        setstartDateOptions(response.data);
      } catch (error) {
        console.error('Error fetching Start Date options:', error);
 
      }
    };

    fetchStartDate();
  }, []);

  useEffect(() => {
    const fetchEndDate = async () => {
      try {
        const response = await axios.get('http://localhost:8080/workingDayoptions');
        setendDateOptions(response.data);
      } catch (error) {
        console.error('Error fetching Start Date options:', error);
 
      }
    };

    fetchEndDate();
  }, []);

  useEffect(() => {
    const fetchDepartments = async () => {
      try {
        const response = await axios.get('http://localhost:8080/timetable/options');
        setDeptOptions(response.data);
      } catch (error) {
        console.error('Error fetching department options:', error);
 
      }
    };

    fetchDepartments();
  }, []);

  useEffect(() => {
    const fetchSection = async () => {
      try {
        const response = await axios.get('http://localhost:8080/timetable/sectionoptions');
        setSectionOptions(response.data);
      } catch (error) {
        console.error('Error fetching section options:', error);
    
      }
    };

    fetchSection();
  }, []);

  // Fetch semester options
  useEffect(() => {
    const fetchSemesters = async () => {
      try {
        const response = await axios.get('http://localhost:8080/timetable/semoptions');
        setSemOptions(response.data);
      } catch (error) {
        console.error('Error fetching semester options:', error);
   
      }
    };

    fetchSemesters();
  }, []);

  // Fetch academic year options
  useEffect(() => {
    const fetchAcademicYears = async () => {
      try {
        const response = await axios.get('http://localhost:8080/acdemicYearOptions');
        setAcademicsOptions(response.data);
      } catch (error) {
        console.error('Error fetching academic year options:', error);
   
      }
    };

    fetchAcademicYears();
  }, []);

  // Function to filter semesters based on academic year label
  useEffect(() => {
    if (academicYear && academicYear.label) {
      const isOddYear = academicYear.label.includes("ODD"); // Check if the academic year label contains 'ODD'

      const filteredSemOptions = semOptions.filter(sem => {
        const semNumber = parseInt(sem.label.replace(/^\D+/g, ''), 10); // Extract the number from the semester label
        return isOddYear ? semNumber % 2 !== 0 : semNumber % 2 === 0;
      });

      setFilteredSemOptions(filteredSemOptions);
    } else {
      setFilteredSemOptions(semOptions); // Reset to show all if no academic year is selected
    }
  }, [academicYear, semOptions]);

  const handleViewTimetable = () => {
    if (department && semester && academicYear && section && startDate && endDate) {
      setViewedDepartment(department.value);
      setViewedSemester(semester.value);
      setViewedAcademic(academicYear.value);
      setViewedSection(section.value)
      setViewedstartDate(startDate.value)
      setViewedendDate(endDate.value)
      setIsOpen(true);
    } else {
      console.error('Please select all required options (department, semester, academic year).');
    }
  };

 

  return (
    <AppLayout
      rId={3}
      title="Time Table"
      body={
        <div style={{ backgroundColor: "white", padding: 17, marginTop: 20, borderRadius: "10px" }}>
          <div style={{ display: 'flex', flexDirection: 'row', columnGap: 10, alignItems: "center" }}>
            <CustomSelect
              placeholder="ACADEMIC YEAR"
              value={academicYear}
              onChange={setAcademicYear}
              options={academicsOptions}
            />
            <CustomSelect
              placeholder="SEMESTER"
              value={semester}
              onChange={setSemester}
              options={filteredSemOptions} // Use filtered semester options
            />
            <CustomSelect
              placeholder="DEPARTMENT"
              value={department}
              onChange={setDepartment}
              options={deptOptions}
            />
             <CustomSelect
              placeholder="SECTION"
              value={section}
              onChange={setSection}
              options={sectionOptions}
            />
             <CustomSelect
              placeholder="START DATE"
              value={startDate}
              onChange={setstartDate}
              options={startDateOptions}
            />
             <CustomSelect
              placeholder="END DATE"
              value={endDate}
              onChange={setendDate}
              options={endDateOptions}
            />
               
            <CustomButton
              width="150"
              label="Generate Timetable"
              onClick={handleViewTimetable}
              backgroundColor="#0878d3"
            />
           
          </div>

          {(viewedDepartment && viewedSemester && viewedAcademic && viewedSection && viewedstartDate && viewedendDate &&  isOpen) && 
            <Timetable departmentID={viewedDepartment} semesterID={viewedSemester} academicYearID={viewedAcademic} sectionID = {viewedSection} startDate = {viewedstartDate} endDate = {viewedendDate} />
          }
        </div>
      }
    />
  );
};

export default GenerateTimetable;
