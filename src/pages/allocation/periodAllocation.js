import AppLayout from "../../layout/layout";
import React, { useState } from "react";
import * as XLSX from "xlsx";
import axios from "axios";
import { Modal, Fade, Button } from "@mui/material";
import "./style.css";

const PeriodAllocation = () => {
  const [excelData, setExcelData] = useState([]);
  const [openModal, setOpenModal] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");

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
  
    axios.post('http://localhost:8080/upload', jsonData)
      .then(response => {
        // Log the response for debugging
        console.log("Server response:", response);
  
        // Assuming the backend responds with a success message
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

  return (
    <AppLayout
      rId={6}
      title="Period Allocation"
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
              <h3>Here you can upload the Period Allocation list</h3>
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
                    {excelData.slice(1).map((row, rowIndex) => (
                      <tr key={rowIndex} className="table-row">
                        {row.map((cell, cellIndex) => (
                          <td key={cellIndex} className="table-cell">{cell}</td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div><br />
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
}

export default PeriodAllocation;
