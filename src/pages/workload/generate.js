// import React, { useState, useEffect } from 'react';
// import axios from 'axios';

// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';
// import Timetable from './workload';


// const GenerateTimetable = () => {
 

//   const [department, setDepartment] = useState(null);
//   const [deptOptions, setDeptOptions] = useState([]);
//   const [semester, setSemester] = useState(null);
//   const [semOptions, setSemOptions] = useState([]);
//   const [isOpen, setIsOpen] = useState(false);
//   const [viewedDepartment, setViewedDepartment] = useState(null);
//   const [viewedSemester, setViewedSemester] = useState(null);
//   const [viewedAcademic, setViewedAcademic] = useState(null);
//   const [academicYear, setAcademicYear] = useState(null);
//   const [academicsOptions, setAcademicsOptions] = useState([]);

//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/options')
//       .then(response => {
//         setDeptOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching department options:', error);
//       });
//   }, []);

//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/semoptions')
//       .then(response => {
//         setSemOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching semester options:', error);
//       });
//   }, []);

//   useEffect(() => {
//     axios.get('http://localhost:8080/acdemicYearOptions')
//       .then(response => {
//         setAcademicsOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching semester options:', error);
//       });
//   }, []);

//   const handleViewTimetable = () => {
//     if (department && semester) {
//       setViewedDepartment(department.value);
//       setViewedSemester(semester.value);
//       setViewedAcademic(academicYear.value);
//       setIsOpen(true);
//     } else {
//       console.error('Please select both department and semester');
//     }
//   };

//   return (
//     <AppLayout
//       rId={3}
//       title="Time Table"
//       body={
//         <div style={{ backgroundColor: "white", padding: 17, marginTop: 20, borderRadius: "10px" }}>
//           <div style={{ display: 'flex', flexDirection: 'row', columnGap: 10, alignItems: "center" }}>
           
//               <CustomSelect
//               placeholder="ACADEMIC YEAR"
//               value={academicYear}
//               onChange={setAcademicYear}
//               options={academicsOptions}
//             />
//             <CustomSelect
//               placeholder="SEMESTER"
//               value={semester}
//               onChange={setSemester}
//               options={semOptions}
//             />
//             <CustomSelect
//               placeholder="DEPARTMENT"
//               value={department}
//               onChange={setDepartment}
//               options={deptOptions}
//             />
          
//             <CustomButton
//               width="150"
//               label="Generate Timetable"
//               onClick={handleViewTimetable}
//               backgroundColor="#0878d3"
//             />
//           </div>

//           {(viewedDepartment && viewedSemester && viewedAcademic && isOpen) && 
//             <Timetable departmentID={viewedDepartment} semesterID={viewedSemester} academicYearID = {viewedAcademic} />
//           }
//         </div>
//       }
//     />
//   );
// };

// export default GenerateTimetable;

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
  // Fetch department options
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
    if (department && semester && academicYear && section && daysCount) {
      setViewedDepartment(department.value);
      setViewedSemester(semester.value);
      setViewedAcademic(academicYear.value);
      setViewedSection(section.value)
      setVieweddaysCount(daysCount)
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
               <input
              type="text"
              placeholder="Enter the Number Of days"
              value={daysCount}
              onChange={(e) => setdaysCount(e.target.value)}
              className="generate-search-input"
            />
            <CustomButton
              width="150"
              label="Generate Timetable"
              onClick={handleViewTimetable}
              backgroundColor="#0878d3"
            />
           
          </div>

          {(viewedDepartment && viewedSemester && viewedAcademic && viewedSection && vieweddaysCount &&  isOpen) && 
            <Timetable departmentID={viewedDepartment} semesterID={viewedSemester} academicYearID={viewedAcademic} sectionID = {viewedSection} day = {vieweddaysCount}  />
          }
        </div>
      }
    />
  );
};

export default GenerateTimetable;
