// import React, { useState, useEffect } from 'react';
// import axios from 'axios';

// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';
// import LabTimetable from './labtable';

// const Lab = () => {


//   const [labOptions, setLabOptions] = useState([]);
//   const [selectedLab, setSelectedLab] = useState("");
//   const [isOpen,setIsOpen] = useState(false)
//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/labOptions')
//       .then(response => {
//         setLabOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching lab subject names:', error);
//       });
//   }, []);

//   const handleViewTimetable = () => {
//     if (selectedLab) {
//       setIsOpen(true);
//     } else {
//       console.error('Please select a lab options');
//     }
//   };
//     return (
//         <AppLayout
//           rId={5}
//           title="Lab Table"
//           body={
//             <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
//                              <div style={{display:'flex',flexDirection:'row',columnGap:10,alignItems:"center"}}>
//             <CustomSelect
//             placeholder="Lab Name"
//             value={selectedLab}
//             onChange={setSelectedLab}
//             options={labOptions}
//           />
        
          
//             <CustomButton
//               width="150"
//               label="View Timetable"
//               onClick={handleViewTimetable}
//             />
        
//           </div>
//           { (selectedLab && isOpen) && 
//           <LabTimetable subjectName ={selectedLab.value} />
             
//           }
//                 </div>
//           }
//           />
//         );
// };

// export default Lab;



import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './lab.css'; // Make sure this CSS file is properly linked
import AppLayout from '../../layout/layout';
import { ArrowBackIosRounded, ArrowForwardIosRounded, VisibilityRounded } from '@mui/icons-material';
import LabTimetable from './labtable';

const Lab = () => {
  const [labData, setLabData] = useState([]);
  const [filteredData, setFilteredData] = useState([]); // Initialized as an empty array
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedLab, setSelectedLab] = useState(null); // State to track selected lab
  const [isOpen, setIsOpen] = useState(false); // State to track if timetable is open

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/timetable/labOptions`);
        setLabData(response.data || []);
        setFilteredData(response.data || []) ;
        setLoading(false);
      } catch (err) {
        setError('Error fetching lab data');
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const results = labData.filter((item) =>
      item.label && item.label.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredData(results);
    setCurrentPage(1); // Reset to the first page when search term changes
  }, [searchTerm, labData]);

  const handleActionClick = (lab) => {
    // If the same lab is clicked again, toggle the timetable display
    if (selectedLab && selectedLab.value === lab.value) {
      setIsOpen(!isOpen);
    } else {
      setSelectedLab(lab);
      setIsOpen(true);
    }
  };

  // Pagination logic
  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = Array.isArray(filteredData) ? filteredData.slice(indexOfFirstRow, indexOfLastRow) : [];
  const totalPages = Math.ceil((filteredData?.length || 0) / rowsPerPage);

  if (selectedLab && isOpen) {
    return (
      <AppLayout
        rId={5}
        title="Lab Table"
        body={<LabTimetable subjectName={selectedLab.value} />}
      />
    );
  }

  return (
    <AppLayout
      rId={5}
      title="Lab Table"
      body={
        <div className="lab-timetable-container">
          <div className="lab-timetable-header">
            <input
              type="text"
              placeholder="Search by lab name..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="lab-timetable-search-input"
            />
          </div>
          <table className="lab-timetable-table">
            <thead className="lab-timetable-head">
              <tr>
                <td>S.No</td>
                <td>Lab Name</td>
                <td>Action</td>
              </tr>
            </thead>
            <tbody className="lab-timetable-body">
              {currentRows.length > 0 ? (
                currentRows.map((item, index) => (
                  <tr key={`${item.value}-${index}`} className="lab-timetable-row">
                    <td className="lab-timetable-cell">{indexOfFirstRow + index + 1}</td>
                    <td className="lab-timetable-cell">{item.label}</td>
                    <td className="lab-timetable-cell">
                      <VisibilityRounded
                        className="dashboard-view-icon"
                        onClick={() => handleActionClick(item)}
                      />
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="3" className="lab-timetable-cell">No data available</td>
                </tr>
              )}
            </tbody>
          </table>
          <div className="lab-timetable-pagination">
            <span className="lab-timetable-pagination-text">
              Page {currentPage} of {totalPages}
            </span>
            <div className="lab-timetable-pagination-right">
              <label htmlFor="rowsPerPage" className="lab-timetable-pagination-text">
                Rows per page:
              </label>
              <select
                id="rowsPerPage"
                value={rowsPerPage}
                onChange={(e) => {
                  setRowsPerPage(parseInt(e.target.value, 10));
                  setCurrentPage(1);
                }}
                className="lab-timetable-pagination-dropdown"
              >
                <option value={8}>8</option>
                <option value={20}>20</option>
                <option value={50}>50</option>
                <option value={100}>100</option>
              </select>
              <button
                onClick={() => setCurrentPage(currentPage - 1)}
                disabled={currentPage === 1}
                className="lab-pagination-button"
                aria-label="Previous Page"
              >
                <ArrowBackIosRounded fontSize="small" />
              </button>
              <button
                onClick={() => setCurrentPage(currentPage + 1)}
                disabled={indexOfLastRow >= filteredData.length}
                className="lab-pagination-button"
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

export default Lab;
