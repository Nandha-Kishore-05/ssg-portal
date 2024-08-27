

// import React, { useState, useEffect } from 'react';
// import axios from 'axios';

// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';
// import FacultyTimetable from './facultytable';

// const FacTimetable = () => {

//   const [facultyOptions, setFacultyOptions] = useState([]);
//   const [selectedFaculty, setSelectedFaculty] = useState("");
//   const [isOpen,setIsOpen] = useState(false)
//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/facultyOptions')
//       .then(response => {
//         setFacultyOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching faculty names:', error);
//       });
//   }, []);

//   const handleViewTimetable = () => {
//     if (selectedFaculty) {
//       setIsOpen(true);
//     } else {
//       console.error('Please select a faculty');
//     }
//   };

//   return (
//     <AppLayout
//       rId={4}
//       title="Faculty Table"
//       body={
//         <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
//                          <div style={{display:'flex',flexDirection:'row',columnGap:10,alignItems:"center"}}>
//           <CustomSelect
//             placeholder="Faculty Name"
//             value={selectedFaculty}
//             onChange={setSelectedFaculty}
//             options={facultyOptions}
//           />
          
//             <CustomButton
//               width="150"
//               label="View Timetable"
//               onClick={handleViewTimetable}
//             />
//             </div>
       
//           { (selectedFaculty && isOpen) && 
//           <FacultyTimetable facultyName ={selectedFaculty.value} />
             
//           }
//         </div>
//       }
//     />
//   );
// };

// export default FacTimetable;

// import React, { useEffect, useState } from 'react';
// import axios from 'axios';
// import './Fac.css'; // Make sure this CSS file is properly linked
// import AppLayout from '../../layout/layout';
// import { ArrowBackIosRounded, ArrowForwardIosRounded } from '@mui/icons-material';

// const FacTimetable = () => {
//   const [facultyData, setFacultyData] = useState([]);
//   const [filteredData, setFilteredData] = useState([]);
//   const [loading, setLoading] = useState(true);
//   const [error, setError] = useState(null);
//   const [searchTerm, setSearchTerm] = useState('');
//   const [currentPage, setCurrentPage] = useState(1);
//   const [rowsPerPage, setRowsPerPage] = useState(10);

//   useEffect(() => {
//     const fetchData = async () => {
//       try {
//         const response = await axios.get(`http://localhost:8080/timetable/facultyOptions`);
//         setFacultyData(response.data);
//         setFilteredData(response.data);
//         setLoading(false);
//       } catch (err) {
//         setError('Error fetching faculty data');
//         setLoading(false);
//       }
//     };

//     fetchData();
//   }, []);

//   useEffect(() => {
//     const results = facultyData.filter((item) =>
//       item.name && item.name.toLowerCase().includes(searchTerm.toLowerCase())
//     );
//     setFilteredData(results);
//     setCurrentPage(1); // Reset to first page when search term changes
//   }, [searchTerm, facultyData]);
  

//   if (loading) return <div>Loading...</div>;
//   if (error) return <div>{error}</div>;

//   // Pagination logic
//   const indexOfLastRow = currentPage * rowsPerPage;
//   const indexOfFirstRow = indexOfLastRow - rowsPerPage;
//   const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
//   const totalPages = Math.ceil(filteredData.length / rowsPerPage);

//   return (
//     <AppLayout
//       rId={4}
//       title="Faculty Table"
//       body={
//         <div className="faculty-timetable-container">
//           <div className="faculty-timetable-header">
//             <input
//               type="text"
//               placeholder="Search by faculty name..."
//               value={searchTerm}
//               onChange={(e) => setSearchTerm(e.target.value)}
//               className="faculty-timetable-search-input"
//             />
//           </div>
//           <table className="faculty-timetable-table">
//             <thead className="faculty-timetable-head">
//               <tr>
//                 <td>S.No</td>
//                 <td>Faculty Name</td>
//                 <td>Action</td>
//               </tr>
//             </thead>
//             <tbody className="faculty-timetable-body">
//               {currentRows.length > 0 ? (
//                 currentRows.map((item, index) => (
//                   <tr key={item.id} className="faculty-timetable-row">
//                     <td className="faculty-timetable-cell">{indexOfFirstRow + index + 1}</td>
//                     <td className="faculty-timetable-cell">{item.name}</td>
//                     <td className="faculty-timetable-cell">
//                       <span className="faculty-timetable-action">View</span>
//                     </td>
//                   </tr>
//                 ))
//               ) : (
//                 <tr>
//                   <td colSpan="3" className="faculty-timetable-cell">No data available</td>
//                 </tr>
//               )}
//             </tbody>
//           </table>
//           <div className="faculty-timetable-pagination">
           
//             <span className="faculty-timetable-pagination-text">
//               Page {currentPage} of {totalPages}
//             </span>
            
//             <div className="faculty-timetable-pagination-right">
//               <label htmlFor="rowsPerPage" className="faculty-timetable-pagination-text">
//                 Rows per page:
//               </label>
//               <select
//                 id="rowsPerPage"
//                 value={rowsPerPage}
//                 onChange={(e) => {
//                   setRowsPerPage(parseInt(e.target.value, 10));
//                   setCurrentPage(1);
//                 }}
//                 className="faculty-timetable-pagination-dropdown"
//               >
//                 <option value={10}>10</option>
//                 <option value={20}>20</option>
//                 <option value={50}>50</option>
//                 <option value={100}>100</option>
//               </select>
//               <button
//                   onClick={() => setCurrentPage(currentPage - 1)}
//                   disabled={currentPage === 1}
//                   className="dashboard-pagination-button"
//                   aria-label="Previous Page"
//                 >
//                   <ArrowBackIosRounded fontSize="small" />
//                 </button>
//                 <button
//                   onClick={() => setCurrentPage(currentPage + 1)}
//                   disabled={indexOfLastRow >= filteredData.length}
//                   className="dashboard-pagination-button"
//                   aria-label="Next Page"
//                 >
//                   <ArrowForwardIosRounded fontSize="small" />
//                 </button>
//             </div>
//           </div>
//         </div>
//       }
//     />
//   );
// };

// export default FacTimetable;

// import React, { useEffect, useState } from 'react';
// import axios from 'axios';
// import './Fac.css'; // Make sure this CSS file is properly linked
// import AppLayout from '../../layout/layout';
// import { ArrowBackIosRounded, ArrowForwardIosRounded, VisibilityRounded } from '@mui/icons-material';

// const FacTimetable = () => {
//   const [facultyData, setFacultyData] = useState([]);
//   const [filteredData, setFilteredData] = useState([]);
//   const [loading, setLoading] = useState(true);
//   const [error, setError] = useState(null);
//   const [searchTerm, setSearchTerm] = useState('');
//   const [currentPage, setCurrentPage] = useState(1);
//   const [rowsPerPage, setRowsPerPage] = useState(8);

//   useEffect(() => {
//     const fetchData = async () => {
//       try {
//         const response = await axios.get(`http://localhost:8080/timetable/facultyOptions`);
//         setFacultyData(response.data);
//         setFilteredData(response.data);
//         setLoading(false);
//       } catch (err) {
//         setError('Error fetching faculty data');
//         setLoading(false);
//       }
//     };

//     fetchData();
//   }, []);

//   useEffect(() => {
//     const results = facultyData.filter((item) =>
//       item.label && item.label.toLowerCase().includes(searchTerm.toLowerCase())
//     );
//     setFilteredData(results);
//     setCurrentPage(1); // Reset to first page when search term changes
//   }, [searchTerm, facultyData]);

//   if (loading) return <div>Loading...</div>;
//   if (error) return <div>{error}</div>;

//   // Pagination logic
//   const indexOfLastRow = currentPage * rowsPerPage;
//   const indexOfFirstRow = indexOfLastRow - rowsPerPage;
//   const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
//   const totalPages = Math.ceil(filteredData.length / rowsPerPage);

//   return (
//     <AppLayout
//       rId={4}
//       title="Faculty Table"
//       body={
//         <div className="faculty-timetable-container">
//           <div className="faculty-timetable-header">
//             <input
//               type="text"
//               placeholder="Search by faculty name..."
//               value={searchTerm}
//               onChange={(e) => setSearchTerm(e.target.value)}
//               className="faculty-timetable-search-input"
//             />
//           </div>
//           <table className="faculty-timetable-table">
//             <thead className="faculty-timetable-head">
//               <tr>
//                 <td>S.No</td>
//                 <td>Faculty Name</td>
//                 <td>Action</td>
//               </tr>
//             </thead>
//             <tbody className="faculty-timetable-body">
//               {currentRows.length > 0 ? (
//                 currentRows.map((item, index) => (
//                   <tr key={`${item.value}-${index}`} className="faculty-timetable-row">
//                     <td className="faculty-timetable-cell">{indexOfFirstRow + index + 1}</td>
//                     <td className="faculty-timetable-cell">{item.label}</td>
//                     <td className="faculty-timetable-cell">
//                     <VisibilityRounded
//                         className="dashboard-view-icon"
                    
//                       />
//                     </td>
//                   </tr>
//                 ))
//               ) : (
//                 <tr>
//                   <td colSpan="3" className="faculty-timetable-cell">No data available</td>
//                 </tr>
//               )}
//             </tbody>
//           </table>
//           <div className="faculty-timetable-pagination">
           
//             <span className="faculty-timetable-pagination-text">
//               Page {currentPage} of {totalPages}
//             </span>
            
//             <div className="faculty-timetable-pagination-right">
//               <label htmlFor="rowsPerPage" className="faculty-timetable-pagination-text">
//                 Rows per page:
//               </label>
//               <select
//                 id="rowsPerPage"
//                 value={rowsPerPage}
//                 onChange={(e) => {
//                   setRowsPerPage(parseInt(e.target.value, 10));
//                   setCurrentPage(1);
//                 }}
//                 className="faculty-timetable-pagination-dropdown"
//               >
//                 <option value={10}>10</option>
//                 <option value={20}>20</option>
//                 <option value={50}>50</option>
//                 <option value={100}>100</option>
//               </select>
//               <button
//                   onClick={() => setCurrentPage(currentPage - 1)}
//                   disabled={currentPage === 1}
//                   className="dashboard-pagination-button"
//                   aria-label="Previous Page"
//                 >
//                   <ArrowBackIosRounded fontSize="small" />
//                 </button>
//                 <button
//                   onClick={() => setCurrentPage(currentPage + 1)}
//                   disabled={indexOfLastRow >= filteredData.length}
//                   className="dashboard-pagination-button"
//                   aria-label="Next Page"
//                 >
//                   <ArrowForwardIosRounded fontSize="small" />
//                 </button>
//             </div>
//           </div>
//         </div>
//       }
//     />
//   );
// };

// export default FacTimetable;

import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './Fac.css'; // Make sure this CSS file is properly linked
import AppLayout from '../../layout/layout';
import { ArrowBackIosRounded, ArrowForwardIosRounded, VisibilityRounded } from '@mui/icons-material';
import FacultyTimetable from './facultytable';

const FacTimetable = () => {
  const [facultyData, setFacultyData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [selectedFaculty, setSelectedFaculty] = useState(null); // State to track selected faculty
  const [isOpen, setIsOpen] = useState(false); // State to track if timetable is open

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/timetable/facultyOptions`);
        setFacultyData(response.data);
        setFilteredData(response.data);
        setLoading(false);
      } catch (err) {
        setError('Error fetching faculty data');
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const results = facultyData.filter((item) =>
      item.label && item.label.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredData(results);
    setCurrentPage(1); // Reset to first page when search term changes
  }, [searchTerm, facultyData]);

  const handleActionClick = (faculty) => {
    // If the same faculty is clicked again, toggle the timetable display
    if (selectedFaculty && selectedFaculty.value === faculty.value) {
      setIsOpen(!isOpen);
    } else {
      setSelectedFaculty(faculty);
      setIsOpen(true);
    }
  };

  

  // Pagination logic
  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
  const totalPages = Math.ceil(filteredData.length / rowsPerPage);


  if (selectedFaculty  && isOpen) {
    return (
      <AppLayout
        rId={4}
       title="Faculty Table"
        body={
          <FacultyTimetable facultyName={selectedFaculty.value} />
        }
      />
    );
  }
  return (
    <AppLayout
      rId={4}
      title="Faculty Table"
      body={
        <div className="faculty-timetable-container">
          <div className="faculty-timetable-header">
            <input
              type="text"
              placeholder="Search by faculty name..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="faculty-timetable-search-input"
            />
          </div>
          <table className="faculty-timetable-table">
            <thead className="faculty-timetable-head">
              <tr>
                <td>S.No</td>
                <td>Faculty Name</td>
                <td>Action</td>
              </tr>
            </thead>
            <tbody className="faculty-timetable-body">
              {currentRows.length > 0 ? (
                currentRows.map((item, index) => (
                  <tr key={`${item.value}-${index}`} className="faculty-timetable-row">
                    <td className="faculty-timetable-cell">{indexOfFirstRow + index + 1}</td>
                    <td className="faculty-timetable-cell">{item.label}</td>
                    <td className="faculty-timetable-cell">
                      <VisibilityRounded
                        className="dashboard-view-icon"
                        onClick={() => handleActionClick(item)}
                      />
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="3" className="faculty-timetable-cell">No data available</td>
                </tr>
              )}
            </tbody>
          </table>
          <div className="faculty-timetable-pagination">
            <span className="faculty-timetable-pagination-text">
              Page {currentPage} of {totalPages}
            </span>
            <div className="faculty-timetable-pagination-right">
              <label htmlFor="rowsPerPage" className="faculty-timetable-pagination-text">
                Rows per page:
              </label>
              <select
                id="rowsPerPage"
                value={rowsPerPage}
                onChange={(e) => {
                  setRowsPerPage(parseInt(e.target.value, 10));
                  setCurrentPage(1);
                }}
                className="faculty-timetable-pagination-dropdown"
              >
                <option value={10}>10</option>
                <option value={20}>20</option>
                <option value={50}>50</option>
                <option value={100}>100</option>
              </select>
              <button
                  onClick={() => setCurrentPage(currentPage - 1)}
                  disabled={currentPage === 1}
                  className="dashboard-pagination-button"
                  aria-label="Previous Page"
                >
                  <ArrowBackIosRounded fontSize="small" />
                </button>
                <button
                  onClick={() => setCurrentPage(currentPage + 1)}
                  disabled={indexOfLastRow >= filteredData.length}
                  className="dashboard-pagination-button"
                  aria-label="Next Page"
                >
                  <ArrowForwardIosRounded fontSize="small" />
                </button>
            </div>
          </div>
          
     
        </div>
      }
    />
  );
};

export default FacTimetable;
