

import React, { useState } from "react";
import * as XLSX from "xlsx";
import axios from "axios";
import { Modal, Fade, Button } from "@mui/material";
import { ArrowBackIosRounded, ArrowForwardIosRounded } from '@mui/icons-material';
import AppLayout from "../layout/layout";



const Bulkentry = () => {
  const [excelData, setExcelData] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10); // Default rows per page
  const [openModal, setOpenModal] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");

  const handleFileUpload = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();

    reader.onload = (event) => {
        const binaryStr = event.target.result;
        const workbook = XLSX.read(binaryStr, { type: "binary" });
        const sheetName = workbook.SheetNames[0];
        const worksheet = workbook.Sheets[sheetName];

        // Parse worksheet to JSON and handle dates/times
        const jsonData = XLSX.utils.sheet_to_json(worksheet, {
            header: 1, // Keep headers
            raw: false, // Parse dates/times automatically
            dateNF: 'yyyy-mm-dd hh:mm:ss', // Custom format for date/time if necessary
        });

        // Manually adjust date/time fields if necessary
        const processedData = jsonData.map((row, rowIndex) => {
            if (rowIndex > 0) { // Skip header row
                row.forEach((cell, colIndex) => {
                    const header = jsonData[0][colIndex];
                    if (header && header.toLowerCase().includes('time')) {
                        // Assume time columns contain 'time' in their header
                        row[colIndex] = XLSX.SSF.format('hh:mm:ss', cell); // Format to time string
                    }
                });
            }
            return row;
        });

        setExcelData(processedData);
    };

    reader.readAsBinaryString(file);
};


  const parseExcelToJson = () => {
    const headers = excelData[0];
    const jsonData = excelData.slice(1).map(row => {
      const rowData = {};
      headers.forEach((header, index) => {
        rowData[header] = row[index];
      });
      return rowData;
    });
    return jsonData;
  };
  const sendDataToBackend = () => {
    const jsonData = parseExcelToJson();
    const formattedData = {
      entries: jsonData.map(entry => ({
        subject_name: entry.Subject,
        department_name: entry.Department,
        semester_id: parseInt(entry.Semester, 10),
        day_name: entry["dayNmae"],  // Ensure this matches what backend expects
        start_time: entry["start time"],
        end_time: entry["end time"],
        faculty_name: entry.Faculty,
        classroom : entry["classroom"],
        status: entry["Lab-subject"] === "NO" ? 1 : 0
      }))
    };
  
    axios.post('http://localhost:8080/manual/bulksubmit', formattedData, {
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(response => {
        console.log("Server response:", response);
        setSuccessMessage("Data uploaded successfully!");
        setOpenModal(true);
      })
      .catch(error => {
        console.error("Error uploading data:", error.response ? error.response.data : error.message);
        setSuccessMessage("Failed to upload data. Please try again.");
        setOpenModal(true);
      });
  };
  
  
  const handleCloseModal = () => {
    setOpenModal(false);
  };

  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = excelData.slice(indexOfFirstRow + 1, indexOfLastRow + 1); // Adjusting for headers

  const totalPages = Math.ceil((excelData.length - 1) / rowsPerPage); // Adjusting for headers

  return (
    <AppLayout
      rId={7}
      title="Manual Entry"
      body={
        <div>
          <input
            accept=".xlsx, .xls"
            className="file-upload-input"
            type="file"
            onChange={handleFileUpload}
          />
          <div className="upload-section">
            <center><br />
              <h2>Here you can upload the Period Allocation list</h2>
              <button
                className="upload-button"
                onClick={() => document.querySelector('.file-upload-input').click()}
              >
                Upload Excel
              </button>
            </center>
          </div>
          {excelData.length > 0 && (
            <div className="table-section">
              <div className="scrollable-table">
                <table className="data-table">
                  <thead>
                    <tr>
                      {excelData[0].map((header, index) => (
                        <th key={index} className="table-header">{header}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {currentRows.map((row, rowIndex) => (
                      <tr key={rowIndex} className="table-row">
                        {row.map((cell, cellIndex) => (
                          <td key={cellIndex} className="table-cell">{cell}</td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <div className="dashboard-pagination">
                <span className="dashboard-pagination-text">
                  Page {currentPage} of {totalPages}
                </span>
                <div className="dashboard-pagination-right">
                  <label htmlFor="rowsPerPage" className="dashboard-pagination-text">
                    Rows per page:
                  </label>
                  <select
                    id="rowsPerPage"
                    value={rowsPerPage}
                    onChange={(e) => {
                      setRowsPerPage(parseInt(e.target.value, 10));
                      setCurrentPage(1);
                    }}
                    className="dashboard-pagination-dropdown"
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
                    disabled={indexOfLastRow >= excelData.length - 1} // Adjusting for headers
                    className="dashboard-pagination-button"
                    aria-label="Next Page"
                  >
                    <ArrowForwardIosRounded fontSize="small" />
                  </button>
                </div>
              </div>
              <br />
              <center>
                <button
                  className="submit-button"
                  onClick={sendDataToBackend}
                >
                  Submit
                </button>
              </center>
            </div>
          )}
          <center>
            <Modal
              open={openModal}
              onClose={handleCloseModal}
              closeAfterTransition
            >
              <Fade in={openModal}>
                <div className="modal-content">
                  <center>
                    <h2>{successMessage}</h2>
                    <Button onClick={handleCloseModal}>Close</Button>
                  </center>
                </div>
              </Fade>
            </Modal>
          </center>
        </div>
      }
    />
  );
};

export default Bulkentry ;
