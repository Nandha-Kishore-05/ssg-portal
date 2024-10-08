import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { ArrowBackIosRounded, ArrowForwardIosRounded, VisibilityRounded } from '@mui/icons-material';
import AppLayout from '../../../layout/layout';
import VenueTimetable from './venueTimetable';

const Venue = () => {
  const [venueData, setVenueData] = useState([]); // Changed from labData to venueData
  const [filteredData, setFilteredData] = useState([]); 
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedVenue, setSelectedVenue] = useState(null); // Changed from selectedLab to selectedVenue
  const [isOpen, setIsOpen] = useState(false); 

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/venueTimetableOptions`);
        setVenueData(response.data || []); // Changed from labData to venueData
        setFilteredData(response.data || []);
        setLoading(false);
      } catch (err) {
        setError('Error fetching venue data'); // Updated error message
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const results = venueData.filter((item) =>
      item.classroom && item.classroom.toLowerCase().includes(searchTerm.toLowerCase()) // Changed from label to classroom
    );
    setFilteredData(results);
    setCurrentPage(1); 
  }, [searchTerm, venueData]);

  const handleActionClick = (venue) => { // Changed from lab to venue
    if (selectedVenue && selectedVenue.classroomValue === venue.classroomValue) { // Changed from selectedLab to selectedVenue
      setIsOpen(!isOpen);
    } else {
      setSelectedVenue(venue); // Changed from selectedLab to selectedVenue
      setIsOpen(true);
    }
  };

  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = Array.isArray(filteredData) ? filteredData.slice(indexOfFirstRow, indexOfLastRow) : [];
  const totalPages = Math.ceil((filteredData?.length || 0) / rowsPerPage);

  if (selectedVenue && isOpen) { // Changed from selectedLab to selectedVenue
    return (
      <AppLayout
        rId={12}
        title="Venue Table"
        body={<VenueTimetable venueName={selectedVenue.classroomValue} />} // Adjust this based on your next steps
      />
    );
  }

  return (
    <AppLayout
      rId={12}
      title="Venue Table"
      body={
        <div className="lab-timetable-container"> {/* Keeping class name unchanged */}
          <div className="lab-timetable-header">
            <input
              type="text"
              placeholder="Search by venue name..." // Changed placeholder text
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="lab-timetable-search-input"
            />
          </div>
          <table className="lab-timetable-table">
            <thead className="lab-timetable-head">
              <tr>
                <td>S.No</td>
                <td>Classroom</td>
                <td>Academic Year</td>
                <td>Action</td>
              </tr>
            </thead>
            <tbody className="lab-timetable-body">
              {currentRows.length > 0 ? (
                currentRows.map((item, index) => (
                  <tr key={`${item.value}-${index}`} className="lab-timetable-row">
                    <td className="lab-timetable-cell">{indexOfFirstRow + index + 1}</td>
                    <td className="lab-timetable-cell">{item.classroom}</td>
                    <td className="lab-timetable-cell">{item.academicyear}</td>
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
                  <td colSpan="4" className="lab-timetable-cell">No data available</td> {/* Adjusted column span */}
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

export default Venue;
