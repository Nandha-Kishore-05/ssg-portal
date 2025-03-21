

import React, { useEffect, useState } from "react";
import * as XLSX from "xlsx";
import axios from "axios";
import { Modal, Fade, Button } from "@mui/material";
import { ArrowBackIosRounded, ArrowForwardIosRounded } from '@mui/icons-material';
import AppLayout from "../layout/layout";
import CustomSelect from "../components/select";



const Bulkentry = () => {
  const [excelData, setExcelData] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10); // Default rows per page
  const [openModal, setOpenModal] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [semester, setSemester] = useState(null);
  const [semOptions, setSemOptions] = useState([]);
  const [academicYear, setAcademicYear] = useState(null);
  const [academicsOptions, setAcademicsOptions] = useState([]);
  const [filteredSemOptions, setFilteredSemOptions] = useState([]);

 

  useEffect(() => {
    axios.get('http://localhost:8080/timetable/semoptions')
      .then(response => {
        setSemOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching semester options:', error);
      });
  }, []);

  useEffect(() => {
    axios.get('http://localhost:8080/acdemicYearOptions')
      .then(response => {
        setAcademicsOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching academic year options:', error);
      });
  }, []);

  useEffect(() => {
    if (academicYear && academicYear.label) {
      // Check if the academic year label contains 'ODD' or 'EVEN' (case-insensitive)
      const isOddYear = /ODD/i.test(academicYear.label); // Check for 'ODD' in a case-insensitive manner

      const filteredSemOptions = semOptions.filter(sem => {
        const semNumber = parseInt(sem.label.replace(/^\D+/g, ''), 10); // Extract the number from the semester label
        return isOddYear ? [1, 3, 5, 7].includes(semNumber) : [2, 4, 6, 8].includes(semNumber);
      });

      setFilteredSemOptions(filteredSemOptions);
    } else {
      setFilteredSemOptions(semOptions); // Reset to show all if no academic year is selected
    }
  }, [academicYear, semOptions]);

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
            dateNF: 'yyyy-mm-dd', // Custom format for date/time if necessary
        });

        // Manually adjust date/time fields if necessary
        const processedData = jsonData.map((row, rowIndex) => {
          if (rowIndex > 0) { // Skip header row
              row.forEach((cell, colIndex) => {
                  const header = jsonData[0][colIndex];
      
                  if (header) {
                      if (header.toLowerCase().includes('day') || header.toLowerCase().includes('date')) {
                          // Convert to YYYY-MM-DD format
                          if (cell instanceof Date) {
                              row[colIndex] = cell.toISOString().split("T")[0]; // Ensures YYYY-MM-DD format
                          } else if (typeof cell === "string") {
                              row[colIndex] = cell.replace(/\//g, "-"); // Convert 2024/12/16 → 2024-12-16
                          }
                      }
      
                      if (header.toLowerCase().includes('time')) {
                          // Ensure time remains formatted correctly
                          row[colIndex] = XLSX.SSF.format('hh:mm:ss', cell);
                      }
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
    const jsonData = parseExcelToJson(); // Replace with your function to parse Excel data.
  console.log(jsonData)
    const formattedData = {
        entries: jsonData.map(entry => ({
            day_name: entry["DAY"] ? entry["DAY"].trim().toUpperCase() : "",
            period: entry["PERIOD"] ? entry["PERIOD"].trim().split(",").map(Number) : [],
            classroom: entry["VENUE"] ? entry["VENUE"].trim().toUpperCase() : "",
            semester_id: semester ? semester.value : null,
            department_name: entry["Department"] ? entry["Department"].trim().toUpperCase().split(",") : [],
            subject_name: entry["Course Name"] ? entry["Course Name"].trim() : "",
            faculty_name: entry["Faculty Name"] ? entry["Faculty Name"].trim() : "",
            subject_type: entry["Subject Type"] ? entry["Subject Type"].trim() : "",
            academic_year: academicYear ? academicYear.value : null,
            course_code: entry["Course Code"] ? entry["Course Code"].trim() : "",
            section: entry["Section"] ? entry["Section"].trim() : "",
        }))
    };

    axios.post('http://localhost:8080/manual/bulksubmit', formattedData, {
        headers: {
            'Content-Type': 'application/json',
        },
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
      rId={1}
      title="Manual Entry"
      body={
        <div>
          <input
            accept=".xlsx, .xls"
            className="file-upload-input"
            type="file"
            onChange={handleFileUpload}
          />

          <div style={{ backgroundColor: "white", padding: 17, marginTop: 20, borderRadius: "10px" }}>
            <div style={{ display: 'flex', flexDirection: 'row', columnGap: 10, alignItems: "center", justifyContent: "space-between" }}>
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
                options={filteredSemOptions} // Use filtered options here
              />
           
              <button
                className="student-upload-button"
                onClick={() => document.querySelector('.file-upload-input').click()}
              >
                Upload Excel
              </button>
            </div>
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
