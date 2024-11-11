


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
import fileDownload from 'js-file-download';

function SaveTimetable() {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [DownloadSemId,setDownloadSemId]=useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedDepartment, setSelectedDepartment] = useState('');
  const [selectedSemester, setSelectedSemester] = useState(''); 
  const [selectedSection, setSelectedSection] = useState(''); 
  const [selectedAcademicyear, setSelectedAcademicyear] = useState('');
  const [departments, setDepartments] = useState([]);
  const [semesters, setSemesters] = useState([]);
  const [isDropdownOpen, setIsDropdownOpen] = useState({ department: false, semester: false });
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();


  useEffect(() => {
    axios.get('http://localhost:8080/saved/deptoptions')
      .then((response) => {
        const fetchedData = response.data || [];
        setData(fetchedData);
        setFilteredData(fetchedData);
  
        // Create a map to filter distinct semesters by id and name
        const uniqueSemesters = [
          ...new Map(
            fetchedData.map((item) => [item.semester_id, { id: item.semester_id, name: item.semester_name }])
          ).values(),
        ];
        setSemesters(uniqueSemesters);
  
        // Get unique department names
        setDepartments([...new Set(fetchedData.map(item => item.department_name))]);
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

  // Pagination logic
  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
  const totalPages = Math.ceil(filteredData.length / rowsPerPage);

  // Handle view click for timetable details
  const handleViewClick = (departmentId, semesterId, academicYearID,sectionID) => {
    setSelectedDepartment(departmentId);
    setSelectedSemester(semesterId);
    setSelectedAcademicyear(academicYearID);
    setSelectedSection(sectionID)
    setIsOpen(true);
  };

  // Dropdown toggle logic
  const toggleDropdown = (type) => {
    setIsDropdownOpen((prevState) => ({ ...prevState, [type]: !prevState[type] }));
  };

  // Select department
  const handleDepartmentSelect = (department) => {
    setSelectedDepartment(department);
    setIsDropdownOpen({ department: false, semester: false });
  };

  // Select semester
  const handleSemesterSelect = (sem) => {
    console.log('Selected Semester ID:', sem.id);
    console.log('Selected Semester Name:', sem.name);
    
    setDownloadSemId(sem.id)
    setSelectedSemester(sem.name); 
    setIsDropdownOpen({ department: false, semester: false });
  };
  


  const downloadTimetable = () => {
    if (!selectedSemester) {
      alert('Please select a semester to download the timetable.');
      return;
    }


    axios.get(`http://localhost:8080/downloadTimetable/${DownloadSemId}`, { responseType: 'blob' })
      .then((response) => {
        fileDownload(response.data, `timetable_${selectedSemester}.xlsx`);
      })
      .catch((error) => {
        console.error('Error downloading the file:', error);
      });
  };

  if (selectedDepartment && selectedSemester && selectedAcademicyear && selectedSection && isOpen) {

    return (
      <AppLayout
        rId={3}
        title="Saved Timetable"
        body={
          <SavedTimetable setIsOpen={setIsOpen} departmentID={selectedDepartment} semesterID={selectedSemester} academicYearID={selectedAcademicyear} sectionID={selectedSection} />
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
                 
                  </td>
                  <td>
  Semester
  <FilterListRounded className="filter-icon" onClick={() => toggleDropdown('semester')} />
  {isDropdownOpen.semester && (
    <div className="filter-dropdown semester-filter-dropdown">
      {semesters.map((sem, index) => (
        <div key={index} onClick={() => handleSemesterSelect(sem)}>
          {sem.name} 
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
                          onClick={() => handleViewClick(row.department_id, row.semester_id, row.academic_year_id,row.section_id,row.classroom)}
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
