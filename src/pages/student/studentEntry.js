// import React, { useEffect, useState } from "react";
// import * as XLSX from "xlsx";
// import axios from "axios";
// import { Modal, Fade, Button } from "@mui/material";
// import { ArrowBackIosRounded, ArrowForwardIosRounded } from '@mui/icons-material';
// import AppLayout from "../../layout/layout";
// import CustomSelect from "../../components/select";
// import './studentEntry.css'

// const StudentEntry = () => {
//   const [excelData, setExcelData] = useState([]);
//   const [currentPage, setCurrentPage] = useState(1);
//   const [rowsPerPage, setRowsPerPage] = useState(10); 
//   const [openModal, setOpenModal] = useState(false);
//   const [successMessage, setSuccessMessage] = useState("");
//   const [department, setDepartment] = useState(null);
//   const [deptOptions, setDeptOptions] = useState([]);
//   const [semester, setSemester] = useState(null);
//   const [semOptions, setSemOptions] = useState([]);
//   const [academicYear, setAcademicYear] = useState(null);
//   const [academicsOptions, setAcademicsOptions] = useState([]);

//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/options')
//       .then(response => {
//         setDeptOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching department options:', error);
//       });
//   }, []);

//   useEffect(() => {
//     axios.get('http://localhost:8080/timetable/semoptions')
//       .then(response => {
//         setSemOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching semester options:', error);
//       });
//   }, []);

//   useEffect(() => {
//     axios.get('http://localhost:8080/acdemicYearOptions')
//       .then(response => {
//         setAcademicsOptions(response.data);
//       })
//       .catch(error => {
//         console.error('Error fetching academic year options:', error);
//       });
//   }, []);

//   const handleFileUpload = (e) => {
//     const file = e.target.files[0];
//     const reader = new FileReader();

//     reader.onload = (event) => {
//       const binaryStr = event.target.result;
//       const workbook = XLSX.read(binaryStr, { type: "binary" });
//       const sheetName = workbook.SheetNames[0];
//       const worksheet = XLSX.utils.sheet_to_json(workbook.Sheets[sheetName], { header: 1 });

//       setExcelData(worksheet);
//     };

//     reader.readAsBinaryString(file);
//   };

//   const parseExcelToJson = () => {
//     const headers = excelData[0];
//     const jsonData = excelData.slice(1).map(row => {
//       const rowData = {};
//       headers.forEach((header, index) => {
//         rowData[header] = row[index];
//       });
//       return rowData;
//     });
//     return jsonData;
//   };

//   const sendDataToBackend = () => {
//     const jsonData = parseExcelToJson();
//     console.log("Parsed Excel Data:", jsonData);

//     // Add selected department, semester, and academic year to the request
//     const requestData = {
//       department: department ? department.value : null,
//       semester: semester ? semester.value : null,
//       academicYear: academicYear ? academicYear.value : null,
//       students: jsonData
//     };
//     console.log(requestData)

//     axios.post('http://localhost:8080/studententry/upload', requestData)
//       .then(response => {
//         console.log("Server response:", response);
//         setSuccessMessage("Data uploaded successfully!");
//         setOpenModal(true);
//       })
//       .catch(error => {
//         console.error("Error uploading data:", error);
//         setSuccessMessage("Failed to upload data. Please try again.");
//         setOpenModal(true);
//       });
//   };

//   const handleCloseModal = () => {
//     setOpenModal(false);
//   };

//   const indexOfLastRow = currentPage * rowsPerPage;
//   const indexOfFirstRow = indexOfLastRow - rowsPerPage;
//   const currentRows = excelData.slice(indexOfFirstRow + 1, indexOfLastRow + 1);

//   const totalPages = Math.ceil((excelData.length - 1) / rowsPerPage);

//   return (
//     <AppLayout
//       rId={13}
//       title="Student Allocation"
//       body={
//         <div>
//           <input
//             accept=".xlsx, .xls"
//             className="file-upload-input"
//             type="file"
//             onChange={handleFileUpload}
//           />

//           <div style={{ backgroundColor: "white", padding: 17, marginTop: 20, borderRadius: "10px" }}>
//             <div style={{ display: 'flex', flexDirection: 'row', columnGap: 10, alignItems: "center", justifyContent: "space-between" }}>
//               <CustomSelect
//                 placeholder="DEPARTMENT"
//                 value={department}
//                 onChange={setDepartment}
//                 options={deptOptions}
//               />
//               <CustomSelect
//                 placeholder="SEMESTER"
//                 value={semester}
//                 onChange={setSemester}
//                 options={semOptions}
//               />
//               <CustomSelect
//                 placeholder="ACADEMIC YEAR"
//                 value={academicYear}
//                 onChange={setAcademicYear}
//                 options={academicsOptions}
//               />
//               <button
//                 className="student-upload-button"
//                 onClick={() => document.querySelector('.file-upload-input').click()}
//               >
//                 Upload Excel
//               </button>
//             </div>
//           </div>

//           {excelData.length > 0 && (
//             <div className="table-section">
//               <div className="scrollable-table">
//                 <table className="data-table">
//                   <thead>
//                     <tr>
//                       {excelData[0].map((header, index) => (
//                         <th key={index} className="table-header">{header}</th>
//                       ))}
//                     </tr>
//                   </thead>
//                   <tbody>
//                     {currentRows.map((row, rowIndex) => (
//                       <tr key={rowIndex} className="table-row">
//                         {row.map((cell, cellIndex) => (
//                           <td key={cellIndex} className="table-cell">{cell}</td>
//                         ))}
//                       </tr>
//                     ))}
//                   </tbody>
//                 </table>
//               </div>
//               <div className="dashboard-pagination">
//                 <span className="dashboard-pagination-text">
//                   Page {currentPage} of {totalPages}
//                 </span>
//                 <div className="dashboard-pagination-right">
//                   <label htmlFor="rowsPerPage" className="dashboard-pagination-text">
//                     Rows per page:
//                   </label>
//                   <select
//                     id="rowsPerPage"
//                     value={rowsPerPage}
//                     onChange={(e) => {
//                       setRowsPerPage(parseInt(e.target.value, 10));
//                       setCurrentPage(1);
//                     }}
//                     className="dashboard-pagination-dropdown"
//                   >
//                     <option value={10}>10</option>
//                     <option value={20}>20</option>
//                     <option value={50}>50</option>
//                     <option value={100}>100</option>
//                   </select>
//                   <button
//                     onClick={() => setCurrentPage(currentPage - 1)}
//                     disabled={currentPage === 1}
//                     className="dashboard-pagination-button"
//                     aria-label="Previous Page"
//                   >
//                     <ArrowBackIosRounded fontSize="small" />
//                   </button>
//                   <button
//                     onClick={() => setCurrentPage(currentPage + 1)}
//                     disabled={indexOfLastRow >= excelData.length - 1}
//                     className="dashboard-pagination-button"
//                     aria-label="Next Page"
//                   >
//                     <ArrowForwardIosRounded fontSize="small" />
//                   </button>
//                 </div>
//               </div>
//               <br />
//               <center>
//                 <button
//                   className="submit-button"
//                   onClick={sendDataToBackend}
//                 >
//                   Submit
//                 </button>
//               </center>
//             </div>
//           )}
//           <center>
//             <Modal
//               open={openModal}
//               onClose={handleCloseModal}
//               closeAfterTransition
//             >
//               <Fade in={openModal}>
//                 <div className="modal-content">
//                   <center>
//                     <h2>{successMessage}</h2>
//                     <Button onClick={handleCloseModal}>Close</Button>
//                   </center>
//                 </div>
//               </Fade>
//             </Modal>
//           </center>
//         </div>
//       }
//     />
//   );
// };

// export default StudentEntry;


import React, { useEffect, useState } from "react";
import * as XLSX from "xlsx";
import axios from "axios";
import { Modal, Fade, Button } from "@mui/material";
import { ArrowBackIosRounded, ArrowForwardIosRounded } from '@mui/icons-material';
import AppLayout from "../../layout/layout";
import CustomSelect from "../../components/select";
import './studentEntry.css';

const StudentEntry = () => {
  const [excelData, setExcelData] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10); 
  const [openModal, setOpenModal] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [department, setDepartment] = useState(null);
  const [deptOptions, setDeptOptions] = useState([]);
  const [semester, setSemester] = useState(null);
  const [semOptions, setSemOptions] = useState([]);
  const [academicYear, setAcademicYear] = useState(null);
  const [academicsOptions, setAcademicsOptions] = useState([]);
  const [filteredSemOptions, setFilteredSemOptions] = useState([]);

  useEffect(() => {
    axios.get('http://localhost:8080/timetable/options')
      .then(response => {
        setDeptOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching department options:', error);
      });
  }, []);

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
      const worksheet = XLSX.utils.sheet_to_json(workbook.Sheets[sheetName], { header: 1 });

      setExcelData(worksheet);
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
    console.log("Parsed Excel Data:", jsonData);

    const requestData = {
      department: department ? department.value : null,
      semester: semester ? semester.value : null,
      academicYear: academicYear ? academicYear.value : null,
      students: jsonData
    };
    console.log(requestData)

    axios.post('http://localhost:8080/studententry/upload', requestData)
      .then(response => {
        console.log("Server response:", response);
        setSuccessMessage("Data uploaded successfully!");
        setOpenModal(true);
      })
      .catch(error => {
        console.error("Error uploading data:", error);
        setSuccessMessage("Failed to upload data. Please try again.");
        setOpenModal(true);
      });
  };

  const handleCloseModal = () => {
    setOpenModal(false);
  };

  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = excelData.slice(indexOfFirstRow + 1, indexOfLastRow + 1);
  const totalPages = Math.ceil((excelData.length - 1) / rowsPerPage);

  return (
    <AppLayout
      rId={13}
      title="Student Allocation"
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
              <CustomSelect
                placeholder="DEPARTMENT"
                value={department}
                onChange={setDepartment}
                options={deptOptions}
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
                    disabled={indexOfLastRow >= excelData.length - 1}
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
                  <h2>{successMessage}</h2>
                  <Button onClick={handleCloseModal}>Close</Button>
                </div>
              </Fade>
            </Modal>
          </center>
        </div>
      }
    />
  );
};

export default StudentEntry;
