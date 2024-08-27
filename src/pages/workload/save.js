// import React, { useState, useEffect } from 'react';
// import axios from 'axios';

// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';
// import SavedTimetable from './timetable';

// const SaveTimetable = () => {
 
//   const [department, setDepartment] = useState(null);
//   const [semester, setSemester] = useState(null);
//   const [deptOptions, setDeptOptions] = useState([]);
//   const [semOptions, setSemOptions] = useState([]);
//   const [isOpen,setIsOpen] = useState(false)
//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/options')
//       .then(response => {
//         setDeptOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching faculty names:', error);
//       });
//   }, []);

//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/semoptions')
//       .then(response => {
//         setSemOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching faculty names:', error);
//       });
//   }, []);

//   const handleViewTimetable = () => {
//     if (department && semester) {
//       setIsOpen(true);
//     } else {
//       console.error('Please select both department and semester');
//     }
//   };
//     return (
//         <AppLayout
//           rId={3}
//           title="Venue Table"
//           body={
//             <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
//                           <div style={{display:'flex',flexDirection:'row',columnGap:10,alignItems:"center"}}>
//             <CustomSelect
//             placeholder="DEPARTMENT"
//             value={department}
//             onChange={setDepartment}
//             options={deptOptions}
//           />
      
//           <CustomSelect
//             placeholder="SEMESTER"
//             value={semester}
//             onChange={setSemester}
//             options={semOptions}
           
//           />
         
//             <CustomButton
//               width="150"
//               label="View Timetable"
//               onClick={handleViewTimetable}
//             />
      
//           </div>
//           { (department && semester && isOpen) && 
//           <SavedTimetable departmentID={department.value} semesterID = {semester.value} />
             
//           }
//                 </div>
//           }
//           />

//         );
// };

// export default SaveTimetable;


// import React, { useEffect, useState } from 'react';
// import AppLayout from '../../layout/layout';
// import axios from 'axios';
// import './save.css';
// import ArrowBackIosRounded from '@mui/icons-material/ArrowBackIosRounded';
// import ArrowForwardIosRounded from '@mui/icons-material/ArrowForwardIosRounded';
// import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
// import SavedTimetable from '../workload/timetable';
// import CustomButton from '../../components/button';
// import { useNavigate } from 'react-router-dom'; // Add this import


// function SaveTimetable() {
//   const [data, setData] = useState([]);
//   const [filteredData, setFilteredData] = useState([]);
//   const [searchTerm, setSearchTerm] = useState('');
//   const [currentPage, setCurrentPage] = useState(1);
//   const [rowsPerPage, setRowsPerPage] = useState(10);
//   const [selectedDepartment, setSelectedDepartment] = useState(null); // State for department ID
//   const [selectedSemester, setSelectedSemester] = useState(null); // State for semester ID
//   const [isOpen, setIsOpen] = useState(false); // State to track if SavedTimetable should be open
//   const navigate = useNavigate();

//   useEffect(() => {
//     axios.get('http://localhost:8080/saved/deptoptions')
//       .then((response) => {
//         setData(response.data);
//         setFilteredData(response.data);
//       })
//       .catch((error) => {
//         console.error('Error fetching data:', error);
//       });
//   }, []);

//   useEffect(() => {
//     const results = data.filter((item) =>
//       item.department_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
//       item.semester_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
//       item.classroom.toLowerCase().includes(searchTerm.toLowerCase())
//     );
//     setFilteredData(results);
//     setCurrentPage(1);
//   }, [searchTerm, data]);

//   const indexOfLastRow = currentPage * rowsPerPage;
//   const indexOfFirstRow = indexOfLastRow - rowsPerPage;
//   const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
//   const totalPages = Math.ceil(filteredData.length / rowsPerPage);

//   const handleViewClick = (departmentId, semesterId) => {
//     console.log('View clicked', departmentId, semesterId);
//     setSelectedDepartment(departmentId); // Set the selected department ID
//     setSelectedSemester(semesterId); // Set the selected semester ID
//     setIsOpen(true); // Open the SavedTimetable component
//   };

//   if (selectedDepartment && selectedSemester && isOpen) {
//     return (
//       <AppLayout
//         rId={3}
//        title="Saved Timetable"
//         body={
//           <SavedTimetable departmentID={selectedDepartment} semesterID={selectedSemester} />
//         }
//       />
//     );
//   }

//   return (
//     <AppLayout
//       rId={3}
//       title="Time Table"
//       body={
//         <>
//           <div className='dashboard-container'>
//             <div className="dashboard-header">
//               <input
//                 type="text"
//                 placeholder="Search..."
//                 value={searchTerm}
//                 onChange={(e) => setSearchTerm(e.target.value)}
//                 className="dashboard-search-input"
//               />
//                <CustomButton
//             width="150"
//             label="Generate Timetable"
//             onClick={() => navigate('/timetable')}
//           />
//             </div>
//             <table className="dashboard-table">
//               <thead className="dashboard-table-head">
//                 <tr>
//                   <td>S.No</td>
//                   <td>Department</td>
//                   <td>Semester</td>
//                   <td>Classroom</td>
//                   <td>Action</td>
//                 </tr>
//               </thead>
//               <tbody className="dashboard-table-body">
//                 {currentRows.map((row, index) => (
//                   <tr key={index} className="dashboard-table-row">
//                     <td className="dashboard-table-cell">{indexOfFirstRow + index + 1}</td>
//                     <td className="dashboard-table-cell">{row.department_name}</td>
//                     <td className="dashboard-table-cell">{row.semester_name}</td>
//                     <td className="dashboard-table-cell">{row.classroom}</td>
//                     <td className="dashboard-table-cell">
//                       <VisibilityRounded
//                         className="dashboard-view-icon"
//                         onClick={() => handleViewClick(row.department_id, row.semester_id)}
//                       />
//                     </td>
//                   </tr>
//                 ))}
//               </tbody>
//             </table>
//             <div className="dashboard-pagination">
//               <h4 className="dashboard-pagination-text">
//                 Page {currentPage} of {totalPages}
//               </h4>
//               <div className="dashboard-pagination-right">
//                 <h4 className="dashboard-pagination-text">
//                   Rows per page:
//                 </h4>
//                 <div className="dashboard-pagination-dropdown">
//                   <select
//                     value={rowsPerPage}
//                     onChange={(e) => {
//                       setRowsPerPage(parseInt(e.target.value, 10));
//                       setCurrentPage(1);
//                     }}
//                   >
//                     <option value={10}>10</option>
//                     <option value={20}>20</option>
//                     <option value={50}>50</option>
//                     <option value={100}>100</option>
//                   </select>
//                 </div>
//                 <button
//                   onClick={() => setCurrentPage(currentPage - 1)}
//                   disabled={currentPage === 1}
//                   className="dashboard-pagination-button"
//                 >
//                   <ArrowBackIosRounded fontSize="small" />
//                 </button>
//                 <button
//                   onClick={() => setCurrentPage(currentPage + 1)}
//                   disabled={indexOfLastRow >= filteredData.length}
//                   className="dashboard-pagination-button"
//                 >
//                   <ArrowForwardIosRounded fontSize="small" />
//                 </button>
//               </div>
//             </div>
//           </div>
//         </>
//       }
//     />
//   );
// }

// export default SaveTimetable;


import React, { useEffect, useState } from 'react';
import AppLayout from '../../layout/layout';
import axios from 'axios';
import './save.css';
import ArrowBackIosRounded from '@mui/icons-material/ArrowBackIosRounded';
import ArrowForwardIosRounded from '@mui/icons-material/ArrowForwardIosRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
import SavedTimetable from '../workload/timetable';
import CustomButton from '../../components/button';
import { useNavigate } from 'react-router-dom';

function SaveTimetable() {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedDepartment, setSelectedDepartment] = useState(null);
  const [selectedSemester, setSelectedSemester] = useState(null);
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    axios.get('http://localhost:8080/saved/deptoptions')
      .then((response) => {
        setData(response.data || []);  // Ensure data is an array
        setFilteredData(response.data || []);  // Ensure filteredData is an array
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    const results = data.filter((item) =>
      item.department_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.semester_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.classroom.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredData(results);
    setCurrentPage(1);
  }, [searchTerm, data]);

  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
  const totalPages = Math.ceil(filteredData.length / rowsPerPage);

  const handleViewClick = (departmentId, semesterId) => {
    console.log('View clicked', departmentId, semesterId);
    setSelectedDepartment(departmentId);
    setSelectedSemester(semesterId);
    setIsOpen(true);
  };

  if (selectedDepartment && selectedSemester && isOpen) {
    return (
      <AppLayout
        rId={3}
        title="Saved Timetable"
        body={
          <SavedTimetable setIsOpen={setIsOpen} departmentID={selectedDepartment} semesterID={selectedSemester} />
        }
      />
    );
  }

  return (
    <AppLayout
      rId={3}
      title="Time Table"
      body={
        <>
          <div className='dashboard-container'>
            <div className="dashboard-header">
              <input
                type="text"
                placeholder="Search..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="dashboard-search-input"
              />
              <CustomButton
                width="150"
                label="Generate Timetable"
                onClick={() => navigate('/timetable')}
              />
            </div>
            <table className="dashboard-table">
              <thead className="dashboard-table-head">
                <tr>
                  <td>S.No</td>
                  <td>Department</td>
                  <td>Semester</td>
                  <td>Classroom</td>
                  <td>Action</td>
                </tr>
              </thead>
              <tbody className="dashboard-table-body">
                {currentRows.length > 0 ? (
                  currentRows.map((row, index) => (
                    <tr key={index} className="dashboard-table-row">
                      <td className="dashboard-table-cell">{indexOfFirstRow + index + 1}</td>
                      <td className="dashboard-table-cell">{row.department_name}</td>
                      <td className="dashboard-table-cell">{row.semester_name}</td>
                      <td className="dashboard-table-cell">{row.classroom}</td>
                      <td className="dashboard-table-cell">
                        <VisibilityRounded
                          className="dashboard-view-icon"
                          onClick={() => handleViewClick(row.department_id, row.semester_id)}
                        />
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="5" className="dashboard-table-no-data">
                      No data available
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
            <div className="dashboard-pagination">
              <h4 className="dashboard-pagination-text">
                Page {currentPage} of {totalPages}
              </h4>
              <div className="dashboard-pagination-right">
                <h4 className="dashboard-pagination-text">
                  Rows per page:
                </h4>
                <div className="dashboard-pagination-dropdown">
                  <select
                    value={rowsPerPage}
                    onChange={(e) => {
                      setRowsPerPage(parseInt(e.target.value, 10));
                      setCurrentPage(1);
                    }}
                  >
                    <option value={10}>10</option>
                    <option value={20}>20</option>
                    <option value={50}>50</option>
                    <option value={100}>100</option>
                  </select>
                </div>
                <button
                  onClick={() => setCurrentPage(currentPage - 1)}
                  disabled={currentPage === 1}
                  className="dashboard-pagination-button"
                >
                  <ArrowBackIosRounded fontSize="small" />
                </button>
                <button
                  onClick={() => setCurrentPage(currentPage + 1)}
                  disabled={indexOfLastRow >= filteredData.length}
                  className="dashboard-pagination-button"
                >
                  <ArrowForwardIosRounded fontSize="small" />
                </button>
              </div>
            </div>
          </div>
        </>
      }
    />
  );
}

export default SaveTimetable;
