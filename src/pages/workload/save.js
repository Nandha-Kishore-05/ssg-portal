


// import React, { useEffect, useState } from 'react';
// import AppLayout from '../../layout/layout';
// import axios from 'axios';
// import './save.css';
// import ArrowBackIosRounded from '@mui/icons-material/ArrowBackIosRounded';
// import ArrowForwardIosRounded from '@mui/icons-material/ArrowForwardIosRounded';
// import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
// import SavedTimetable from '../workload/timetable';
// import CustomButton from '../../components/button';
// import { useNavigate } from 'react-router-dom';

// function SaveTimetable() {
//   const [data, setData] = useState([]);
//   const [filteredData, setFilteredData] = useState([]);
//   const [searchTerm, setSearchTerm] = useState('');
//   const [currentPage, setCurrentPage] = useState(1);
//   const [rowsPerPage, setRowsPerPage] = useState(10);
//   const [selectedDepartment, setSelectedDepartment] = useState(null);
//   const [selectedSemester, setSelectedSemester] = useState(null);
//   const [isOpen, setIsOpen] = useState(false);
//   const navigate = useNavigate();

//   useEffect(() => {
//     axios.get('http://localhost:8080/saved/deptoptions')
//       .then((response) => {
//         setData(response.data || []);  
//         setFilteredData(response.data || []); 
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
//     setSelectedDepartment(departmentId);
//     setSelectedSemester(semesterId);
//     setIsOpen(true);
//   };

//   if (selectedDepartment && selectedSemester && isOpen) {
//     return (
//       <AppLayout
//         rId={3}
//         title="Saved Timetable"
//         body={
//           <SavedTimetable setIsOpen={setIsOpen} departmentID={selectedDepartment} semesterID={selectedSemester} />
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
//                      <div className="buttons">
//                  <CustomButton
//                 width="150"
//                 label="Download Timetable"
              
//               />
//               <CustomButton
//                 width="150"
//                 label="Generate Timetable"
//                 onClick={() => navigate('/timetable')}
//               />
//               </div>
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
//                 {currentRows.length > 0 ? (
//                   currentRows.map((row, index) => (
//                     <tr key={index} className="dashboard-table-row">
//                       <td className="dashboard-table-cell">{indexOfFirstRow + index + 1}</td>
//                       <td className="dashboard-table-cell">{row.department_name}</td>
//                       <td className="dashboard-table-cell">{row.semester_name}</td>
//                       <td className="dashboard-table-cell">{row.classroom}</td>
//                       <td className="dashboard-table-cell">
//                         <VisibilityRounded
//                           className="dashboard-view-icon"
//                           onClick={() => handleViewClick(row.department_id, row.semester_id)}
//                         />
//                       </td>
//                     </tr>
//                   ))
//                 ) : (
//                   <tr>
//                     <td colSpan="5" className="dashboard-table-no-data">
//                       No data available
//                     </td>
//                   </tr>
//                 )}
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

// import React, { useEffect, useState } from 'react';
// import AppLayout from '../../layout/layout';
// import axios from 'axios';
// import './save.css';
// import ArrowBackIosRounded from '@mui/icons-material/ArrowBackIosRounded';
// import ArrowForwardIosRounded from '@mui/icons-material/ArrowForwardIosRounded';
// import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
// import FilterListRounded from '@mui/icons-material/FilterListRounded';  // Filter icon
// import SavedTimetable from '../workload/timetable';
// import CustomButton from '../../components/button';
// import { useNavigate } from 'react-router-dom';

// function SaveTimetable() {
//   const [data, setData] = useState([]);
//   const [filteredData, setFilteredData] = useState([]);
//   const [searchTerm, setSearchTerm] = useState('');
//   const [currentPage, setCurrentPage] = useState(1);
//   const [rowsPerPage, setRowsPerPage] = useState(10);
//   const [selectedDepartment, setSelectedDepartment] = useState('');
//   const [selectedSemester, setSelectedSemester] = useState('');
//   const [departments, setDepartments] = useState([]); // List of departments
//   const [semesters, setSemesters] = useState([]); // List of semesters
//   const [isDropdownOpen, setIsDropdownOpen] = useState({ department: false, semester: false });
//   const [isOpen, setIsOpen] = useState(false);
//   const navigate = useNavigate();

//   useEffect(() => {
//     // Fetch data and department/semester options
//     axios.get('http://localhost:8080/saved/deptoptions')
//       .then((response) => {
//         setData(response.data || []);  
//         setFilteredData(response.data || []); 
//         setDepartments([...new Set(response.data.map(item => item.department_name))]); // Unique departments
//         setSemesters([...new Set(response.data.map(item => item.semester_name))]); // Unique semesters
//       })
//       .catch((error) => {
//         console.error('Error fetching data:', error);
//       });
//   }, []);

//   // Filter data based on search term, department, and semester
//   useEffect(() => {
//     const results = data.filter((item) =>
//       (selectedDepartment ? item.department_name === selectedDepartment : true) &&
//       (selectedSemester ? item.semester_name === selectedSemester : true) &&
//       (item.department_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
//       item.semester_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
//       item.classroom.toLowerCase().includes(searchTerm.toLowerCase()))
//     );
//     setFilteredData(results);
//     setCurrentPage(1);
//   }, [searchTerm, selectedDepartment, selectedSemester, data]);

//   const indexOfLastRow = currentPage * rowsPerPage;
//   const indexOfFirstRow = indexOfLastRow - rowsPerPage;
//   const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
//   const totalPages = Math.ceil(filteredData.length / rowsPerPage);

//   const handleViewClick = (departmentId, semesterId) => {
//     setSelectedDepartment(departmentId);
//     setSelectedSemester(semesterId);
//     setIsOpen(true);
//   };

//   const toggleDropdown = (type) => {
//     setIsDropdownOpen((prevState) => ({ ...prevState, [type]: !prevState[type] }));
//   };

//   const handleDepartmentSelect = (department) => {
//     setSelectedDepartment(department);
//     setIsDropdownOpen({ department: false, semester: false });
//   };

//   const handleSemesterSelect = (semester) => {
//     setSelectedSemester(semester);
//     setIsDropdownOpen({ department: false, semester: false });
//   };

//   if (selectedDepartment && selectedSemester && isOpen) {
//     return (
//       <AppLayout
//         rId={3}
//         title="Saved Timetable"
//         body={
//           <SavedTimetable setIsOpen={setIsOpen} departmentID={selectedDepartment} semesterID={selectedSemester} />
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
//               <div className="buttons">
//                 <CustomButton width="150" label="Download Timetable" />
//                 <CustomButton width="150" label="Generate Timetable" onClick={() => navigate('/timetable')} />
//               </div>
//             </div>
//             <table className="dashboard-table">
//               <thead className="dashboard-table-head">
//                 <tr>
//                   <td>S.No</td>
//                   <td>
//                  Department
//   <FilterListRounded className="filter-icon" onClick={() => toggleDropdown('department')} />
//   {isDropdownOpen.department && (
//     <div className="filter-dropdown department-filter-dropdown">  
//       {departments.map((dept, index) => (
//         <div key={index} onClick={() => handleDepartmentSelect(dept)}>
//           {dept}
//         </div>
//       ))}
//     </div>
//   )}
// </td>
// <td>
//   Semester
//   <FilterListRounded className="filter-icon" onClick={() => toggleDropdown('semester')} />
//   {isDropdownOpen.semester && (
//     <div className="filter-dropdown semester-filter-dropdown">  
//       {semesters.map((sem, index) => (
//         <div key={index} onClick={() => handleSemesterSelect(sem)}>
//           {sem}
//         </div>
//       ))}
//     </div>
//   )}
// </td>

//                   <td>Classroom</td>
//                   <td>Action</td>
//                 </tr>
//               </thead>
//               <tbody className="dashboard-table-body">
//                 {currentRows.length > 0 ? (
//                   currentRows.map((row, index) => (
//                     <tr key={index} className="dashboard-table-row">
//                       <td className="dashboard-table-cell">{indexOfFirstRow + index + 1}</td>
//                       <td className="dashboard-table-cell">{row.department_name}</td>
//                       <td className="dashboard-table-cell">{row.semester_name}</td>
//                       <td className="dashboard-table-cell">{row.classroom}</td>
//                       <td className="dashboard-table-cell">
//                         <VisibilityRounded
//                           className="dashboard-view-icon"
//                           onClick={() => handleViewClick(row.department_id, row.semester_id)}
//                         />
//                       </td>
//                     </tr>
//                   ))
//                 ) : (
//                   <tr>
//                     <td colSpan="5" className="dashboard-table-no-data">
//                       No data available
//                     </td>
//                   </tr>
//                 )}
//               </tbody>
//             </table>
//             <div className="dashboard-pagination">
//               <h4 className="dashboard-pagination-text">
//                 Page {currentPage} of {totalPages}
//               </h4>
//               <div className="dashboard-pagination-right">
//                 <h4 className="dashboard-pagination-text">Rows per page:</h4>
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
import FilterListRounded from '@mui/icons-material/FilterListRounded';
import SavedTimetable from '../workload/timetable';
import CustomButton from '../../components/button';
import { useNavigate } from 'react-router-dom';
import ExcelJS from 'exceljs';
import { saveAs } from 'file-saver';

function SaveTimetable() {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedDepartment, setSelectedDepartment] = useState('');
  const [selectedSemester, setSelectedSemester] = useState('');
  const [selectedAcademicyear, setSelectedAcademicyear] = useState('');
  const [departments, setDepartments] = useState([]);
  const [semesters, setSemesters] = useState([]);
  const [isDropdownOpen, setIsDropdownOpen] = useState({ department: false, semester: false });
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {

    axios.get('http://localhost:8080/saved/deptoptions')
      .then((response) => {
        setData(response.data || []);  
        setFilteredData(response.data || []); 
        setDepartments([...new Set(response.data.map(item => item.department_name))]); 
        setSemesters([...new Set(response.data.map(item => item.semester_name))]); 
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

 
  useEffect(() => {
    const results = data.filter((item) =>
      (selectedDepartment ? item.department_name === selectedDepartment : true) &&
      (selectedSemester ? item.semester_name === selectedSemester : true) &&
      (item.department_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.semester_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.classroom.toLowerCase().includes(searchTerm.toLowerCase()))
    );
    setFilteredData(results);
    setCurrentPage(1);
  }, [searchTerm, selectedDepartment, selectedSemester, data]);

  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
  const totalPages = Math.ceil(filteredData.length / rowsPerPage);

  const handleViewClick = (departmentId, semesterId,academicYearID) => {
    setSelectedDepartment(departmentId);
    setSelectedSemester(semesterId);
    setSelectedAcademicyear(academicYearID)
    setIsOpen(true);
  };

  const toggleDropdown = (type) => {
    setIsDropdownOpen((prevState) => ({ ...prevState, [type]: !prevState[type] }));
  };

  const handleDepartmentSelect = (department) => {
    setSelectedDepartment(department);
    setIsDropdownOpen({ department: false, semester: false });
  };

  const handleSemesterSelect = (semester) => {
    setSelectedSemester(semester);
    setIsDropdownOpen({ department: false, semester: false });
  };

  const downloadTimetable = async () => {
    try {
      const workbook = new ExcelJS.Workbook();
      const worksheet = workbook.addWorksheet('Timetable');
  
      worksheet.columns = [
        { header: 'Department', key: 'department_name', width: 20 },
        { header: 'Semester', key: 'semester_name', width: 20 },
        { header: 'Classroom', key: 'classroom', width: 20 },
        // Add other columns as needed
      ];
  
      filteredData.forEach((row) => {
        worksheet.addRow(row);
      });
  
      worksheet.eachRow((row, rowNumber) => {
        if (rowNumber === 1) {
          row.font = { bold: true };
          row.alignment = { horizontal: 'center' };
        }
        row.eachCell({ includeEmpty: true }, (cell, colNumber) => {
          cell.border = {
            top: { style: 'thin' },
            left: { style: 'thin' },
            bottom: { style: 'thin' },
            right: { style: 'thin' },
          };
          cell.alignment = { horizontal: 'center' };
        });
      });
  
      const buffer = await workbook.xlsx.writeBuffer();
      saveAs(new Blob([buffer], { type: 'application/octet-stream' }), 'Timetable.xlsx');
    } catch (error) {
      console.error('Error exporting to Excel:', error);
    }
  };
  

  if (selectedDepartment && selectedSemester && selectedAcademicyear && isOpen) {
    return (
      <AppLayout
        rId={3}
        title="Saved Timetable"
        body={
          <SavedTimetable setIsOpen={setIsOpen} departmentID={selectedDepartment} semesterID={selectedSemester} academicYearID = {selectedAcademicyear} />
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
              <div className="buttons">
                <CustomButton width="150" label="Download Timetable" onClick={downloadTimetable} />
                <CustomButton width="150" label="Generate Timetable" onClick={() => navigate('/timetable')} />
              </div>
            </div>
            <table className="dashboard-table">
              <thead className="dashboard-table-head">
                <tr>
                  <td>S.No</td>
                  <td>
                    Department
                    <FilterListRounded className="filter-icon" onClick={() => toggleDropdown('department')} />
                    {isDropdownOpen.department && (
                      <div className="filter-dropdown department-filter-dropdown">
                        {departments.map((dept, index) => (
                          <div key={index} onClick={() => handleDepartmentSelect(dept)}>
                            {dept}
                          </div>
                        ))}
                      </div>
                    )}
                  </td>
                  <td>
                    Semester
                    <FilterListRounded className="filter-icon" onClick={() => toggleDropdown('semester')} />
                    {isDropdownOpen.semester && (
                      <div className="filter-dropdown semester-filter-dropdown">
                        {semesters.map((sem, index) => (
                          <div key={index} onClick={() => handleSemesterSelect(sem)}>
                            {sem}
                          </div>
                        ))}
                      </div>
                    )}
                  </td>
                  <td>Classroom</td>
                  <td>Academic-Year</td>
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
                      <td className="dashboard-table-cell">{row.academic_year_name}</td>
                      <td className="dashboard-table-cell">
                        <VisibilityRounded
                          className="dashboard-view-icon"
                          onClick={() => handleViewClick(row.department_id, row.semester_id,row.academic_year_id)}
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
                <h4 className="dashboard-pagination-text">Rows per page:</h4>
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
