import React, { useEffect, useState } from 'react';
import AppLayout from '../../layout/layout';
import axios from 'axios';
import './allocation.css';
import ArrowBackIosRounded from '@mui/icons-material/ArrowBackIosRounded';
import ArrowForwardIosRounded from '@mui/icons-material/ArrowForwardIosRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
import { FormControl, MenuItem, Modal, Select, TextField } from '@mui/material';
import CustomButton from '../../components/button';
import { toast, ToastContainer } from 'react-toastify';

function SubjectAllocation() {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [openModal, setOpenModal] = useState(false);
  const [modalData, setModalData] = useState(null);
  const [editedPeriods, setEditedPeriods] = useState('');
  const [editedSubjectName, setEditedSubjectName] = useState('');
  const [editedFacultyName, setEditedFacultyName] = useState('');
  const [editedLabType, setEditedLabType] = useState('');

  useEffect(() => {
    axios.get('http://localhost:8080/periodallocation')
      .then((response) => {
        setData(response.data || []);  
        setFilteredData(response.data || []); 
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {
    const results = data.filter((item) =>
      item.department_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.semester_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.subject_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.faculty_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.status.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredData(results);
    setCurrentPage(1);
  }, [searchTerm, data]);

  const handleOpenModal = (row) => {
    setModalData(row);
    setEditedPeriods(row.periods);
    setEditedSubjectName(row.subject_name);
    setEditedFacultyName(row.faculty_name);
    setEditedLabType(row.status);
    console.log(row);  // This will log the modal data to the console
    setOpenModal(true);
  };
  

  const handleCloseModal = () => {
    setOpenModal(false);
    setModalData(null);
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    if (name === 'periods') {
      setEditedPeriods(value);
    } else if (name === 'subject_name') {
      setEditedSubjectName(value);
    } else if (name === 'faculty_name') {
      setEditedFacultyName(value);
    } else if (name === 'lab_type') {
      setEditedLabType(value);
    }
  };

  // const handleSaveChanges = () => {
  //   if (modalData && editedPeriods > 0) {
  //     axios.put('http://localhost:8080/periodallocationedit', {
  //       id: modalData.id,
  //       periods: parseInt(editedPeriods, 10),
  //       subject_name: editedSubjectName,
  //       old_subject_name: modalData.subject_name, // Sending old subject name
  //       department_id: modalData.department_id,
  //       semester_id: modalData.semester_id,
  //       faculty_name: editedFacultyName,
  //       old_faculty_name: modalData.faculty_name, // Sending old faculty name
  //       faculty_id: modalData.faculty_id,
  //     })
      
  //     .then(() => {
  //       toast.success('Updated successfully');
  //       setData(prevData => prevData.map(item =>
  //         item.id === modalData.id
  //           ? { 
  //               ...item, 
  //               periods: parseInt(editedPeriods, 10), 
  //               subject_name: editedSubjectName, 
  //               faculty_name: editedFacultyName, 
  //               status: editedLabType,
  //               old_subject_name: editedSubjectName, // Update the old subject name
  //               old_faculty_name: editedFacultyName  // Update the old faculty name
  //             }
  //           : item
  //       ));
  //       setFilteredData(prevData => prevData.map(item =>
  //         item.id === modalData.id
  //           ? { 
  //               ...item, 
  //               periods: parseInt(editedPeriods, 10), 
  //               subject_name: editedSubjectName, 
  //               faculty_name: editedFacultyName, 
  //               status: editedLabType,
  //               old_subject_name: editedSubjectName, // Update the old subject name
  //               old_faculty_name: editedFacultyName  // Update the old faculty name
  //             }
  //           : item
  //       ));
  //       handleCloseModal();
  //     })
      
  //     .catch(error => {
  //       console.error('Error updating data:', error);
  //       toast.error('Error updating data');
  //     });
  //   } else {
  //     toast.error('Periods must be greater than 0');
  //   }
  // };

  const handleSaveChanges = () => {
    if (modalData && editedPeriods > 0) {
      axios.put('http://localhost:8080/periodallocationedit', {
        id: modalData.id,
        periods: parseInt(editedPeriods, 10),
        subject_name: editedSubjectName,
        old_subject_name: modalData.subject_name, // Sending old subject name
        department_id: modalData.department_id,
        semester_id: modalData.semester_id,
        faculty_name: editedFacultyName,
        old_faculty_name: modalData.faculty_name, // Sending old faculty name
        faculty_id: modalData.faculty_id,
      })
      
      .then(() => {
        toast.success('Updated successfully');
  
        // Log the old and new data
        console.log('Old data:', modalData);
        console.log('New data:', {
          id: modalData.id,
          periods: parseInt(editedPeriods, 10),
          subject_name: editedSubjectName,
          old_subject_name: modalData.subject_name,
          department_id: modalData.department_id,
          semester_id: modalData.semester_id,
          faculty_name: editedFacultyName,
          old_faculty_name: modalData.faculty_name,
          faculty_id: modalData.faculty_id,
          status: editedLabType
        });
  
        // Update state
        setData(prevData => prevData.map(item =>
          item.id === modalData.id
            ? { 
                ...item, 
                periods: parseInt(editedPeriods, 10), 
                subject_name: editedSubjectName, 
                faculty_name: editedFacultyName, 
                status: editedLabType,
                old_subject_name: editedSubjectName, // Update the old subject name
                old_faculty_name: editedFacultyName  // Update the old faculty name
              }
            : item
        ));
        setFilteredData(prevData => prevData.map(item =>
          item.id === modalData.id
            ? { 
                ...item, 
                periods: parseInt(editedPeriods, 10), 
                subject_name: editedSubjectName, 
                faculty_name: editedFacultyName, 
                status: editedLabType,
                old_subject_name: editedSubjectName, // Update the old subject name
                old_faculty_name: editedFacultyName  // Update the old faculty name
              }
            : item
        ));
        handleCloseModal();
      })
      .catch(error => {
        console.error('Error updating data:', error);
        toast.error('Error updating data');
      });
    } else {
      toast.error('Periods must be greater than 0');
    }
  };
  
  // Pagination
  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
  const totalPages = Math.ceil(filteredData.length / rowsPerPage);

  return (
    <AppLayout
      rId={9}
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
            </div>
            <table className="dashboard-table">
              <thead className="dashboard-table-head">
                <tr>
                  <td>S.No</td>
                  <td>Department</td>
                  <td>Semester</td>
                  <td>Subject</td>
                  <td>Faculty</td>
                  <td>Periods</td>
                  <td>Lab Type</td>
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
                      <td className="dashboard-table-cell">{row.subject_name}</td>
                      <td className="dashboard-table-cell">{row.faculty_name}</td>
                      <td className="dashboard-table-cell">{row.periods}</td>
                      <td className="dashboard-table-cell">{row.status}</td>
                      <td className="dashboard-table-cell">
                        <VisibilityRounded
                          className="dashboard-view-icon"
                          onClick={() => handleOpenModal(row)}
                        />
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="8" className="dashboard-table-no-data">
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
          <ToastContainer />
    
          <Modal
            open={openModal}
            onClose={handleCloseModal}
            aria-labelledby="modal-title"
            aria-describedby="modal-description"
          >
            <div className="modal-box">
              <h2 id="modal-title">Edit Data</h2>
              {modalData && (
                <>
                  <TextField
                    label="Periods"
                    name="periods"
                    value={editedPeriods}
                    onChange={handleEditChange}
                    fullWidth
                    margin="normal"
                    type="number"
                  />
                  <TextField
                    label="Subject Name"
                    name="subject_name"
                    value={editedSubjectName}
                    onChange={handleEditChange}
                    fullWidth
                    margin="normal"
                  />
                  <TextField
                    label="Faculty Name"
                    name="faculty_name"
                    value={editedFacultyName}
                    onChange={handleEditChange}
                    fullWidth
                    margin="normal"
                  />
                  <FormControl fullWidth margin="normal">
                    <Select
                      labelId="lab-type-label"
                      name="lab_type"
                      value={editedLabType}
                      onChange={handleEditChange}
                    >
                      <MenuItem value="Lab Subject">Lab Subject</MenuItem>
                      <MenuItem value="Non-Lab Subject">Non-Lab Subject</MenuItem>
                    </Select>
                  </FormControl>
                  <div className="modal-actions">
                    <CustomButton
                      label="Save"
                      onClick={handleSaveChanges}
                    />
                    <CustomButton
                      label="Cancel"
                      onClick={handleCloseModal}
                      backgroundColor="red"
                    />
                  </div>
                </>
              )}
            </div>
          </Modal>
        </>
      }
    />
  );
}

export default SubjectAllocation;
