import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { ArrowBackIosRounded, ArrowForwardIosRounded, VisibilityRounded } from '@mui/icons-material';
import AppLayout from '../../layout/layout';
import StudentTimetable from './studentTimetable';


const StudentTable = () => {
  const [studentData, setStudentData] = useState([]);
  const [filteredData, setFilteredData] = useState([]); 
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedStudent, setSelectedStudent] = useState(null); 
  const [isOpen, setIsOpen] = useState(false); 

useEffect(() => {
  const fetchData = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/studentoptions`);
      console.log(response.data); 
      setStudentData(response.data || []);
      setFilteredData(response.data || []);
      setLoading(false);
    } catch (err) {
      setError('Error fetching venue data'); 
      setLoading(false);
    }
  };
  fetchData();
}, []);


  useEffect(() => {
    const results = studentData.filter((item) =>
      item.student_name && item.student_name.toLowerCase().includes(searchTerm.toLowerCase()) 
    );
    setFilteredData(results);
    setCurrentPage(1); 
  }, [searchTerm, studentData]);

  const handleActionClick = (student) => {
    if (selectedStudent && selectedStudent.student_id === student.student_id) { 
      setIsOpen(!isOpen);
    } else {
      setSelectedStudent(student); 
      setIsOpen(true);
    }
  };

  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = Array.isArray(filteredData) ? filteredData.slice(indexOfFirstRow, indexOfLastRow) : [];
  const totalPages = Math.ceil((filteredData?.length || 0) / rowsPerPage);

  if (selectedStudent && isOpen) { 
    return (
      <AppLayout
        rId={14}
        title="Student TimeTable"
        body={<StudentTimetable Student={selectedStudent.student_id} StudentName ={selectedStudent.student_name} />} // Adjust this based on your next steps
      />
    );
  }

  return (
    <AppLayout
      rId={14}
      title="Student Table"
      body={
        <div className="lab-timetable-container"> 
          <div className="lab-timetable-header">
            <input
              type="text"
              placeholder="Search by venue name..." 
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="lab-timetable-search-input"
            />
          </div>
          <table className="lab-timetable-table">
            <thead className="lab-timetable-head">
              <tr>
                <td>S.No</td>
                <td>Name</td>
                <td>Roll No</td>
                <td>Departement</td>
                <td>Semester</td>
                <td>Academic Year</td>
                <td>Action</td>
              </tr>
            </thead>
            <tbody className="lab-timetable-body">
  {currentRows.length > 0 ? (
    currentRows.map((item, index) => {
      return (
        <tr key={item.student_id} className="lab-timetable-row"> 
          <td className="lab-timetable-cell">{indexOfFirstRow + index + 1}</td>
          <td className="lab-timetable-cell">{item.student_name}</td> 
          <td className="lab-timetable-cell">{item.roll_no}</td> 
          <td className="lab-timetable-cell">{item.department_name}</td> 
          <td className="lab-timetable-cell">{item.semester_name}</td>
          <td className="lab-timetable-cell">{item.academic_year}</td> 
          <td className="lab-timetable-cell">
            <VisibilityRounded
              className="dashboard-view-icon"
              onClick={() => handleActionClick(item)}
            />
          </td>
        </tr>
      );
    })
  ) : (
    <tr>
      <td colSpan="7" className="lab-timetable-cell">No data available</td>
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

export default StudentTable;
