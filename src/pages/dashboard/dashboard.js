import React, { useState } from 'react';
import AppLayout from '../../layout/layout';
import { Table, TableBody, TableCell, TableContainer, TableHead, TablePagination, TableRow, Paper, Button } from '@mui/material';
import ExcelJS from 'exceljs';
import { saveAs } from 'file-saver';

const ExcelUpload = () => {
  const [file, setFile] = useState(null);
  const [allocations, setAllocations] = useState([]); // State to store the table data
  const [page, setPage] = useState(0); // Pagination state
  const [rowsPerPage, setRowsPerPage] = useState(10); // Rows per page for pagination

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!file) {
      alert("Please select a file");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch("http://localhost:8080/upload-lab", {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        const result = await response.json();
        setAllocations(result); // Store the result in state
      } else {
        console.error("Error uploading file");
      }
    } catch (error) {
      console.error("Error uploading file:", error);
    }
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0); // Reset to first page on rows per page change
  };

  // Function to download the table as an Excel file
  const downloadExcel = () => {
    const workbook = new ExcelJS.Workbook();
    const worksheet = workbook.addWorksheet("Allocations");

    // Add header row
    worksheet.columns = [
      { header: "Department", key: "department", width: 20 },
      { header: "Date", key: "date", width: 20 },
      { header: "Period", key: "period", width: 20 },
      { header: "Venue", key: "venue", width: 20 },
      { header: "Subject", key: "subject", width: 20 },
      { header: "Course Code", key: "course_code", width: 20 },
      { header: "Section", key: "section", width: 10 },
      { header: "Faculty", key: "faculty", width: 20 },
    ];

    // Add data rows
    allocations.forEach((departmentAlloc) => {
      departmentAlloc.allocations.forEach((allocation) => {
        worksheet.addRow({
          department: departmentAlloc.department,
          date: allocation.date,
          period: allocation.period,
          venue: allocation.venue,
          subject: allocation.subject,
          course_code: allocation.course_code,
          section: allocation.section,
          faculty: allocation.faculty,
        });
      });
    });

    // Write the Excel file to a buffer
    workbook.xlsx.writeBuffer().then((buffer) => {
      const blob = new Blob([buffer], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" });
      saveAs(blob, "allocations.xlsx");
    });
  };

  const renderTable = () => {
    if (allocations.length === 0) {
      return <p>No data available</p>;
    }

    const startIdx = page * rowsPerPage;
    const endIdx = startIdx + rowsPerPage;
    const currentAllocations = allocations.slice(startIdx, endIdx);

    return (
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Department</TableCell>
              <TableCell>Date</TableCell>
              <TableCell>Period</TableCell>
              <TableCell>Venue</TableCell>
              <TableCell>Subject</TableCell>
              <TableCell>Course Code</TableCell>
              <TableCell>Section</TableCell>
              <TableCell>Faculty</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {currentAllocations.map((departmentAlloc, index) => (
              <React.Fragment key={index}>
                {departmentAlloc.allocations.map((allocation, idx) => (
                  <TableRow key={idx}>
                    <TableCell>{departmentAlloc.department}</TableCell>
                    <TableCell>{allocation.date}</TableCell>
                    <TableCell>{allocation.period}</TableCell>
                    <TableCell>{allocation.venue}</TableCell>
                    <TableCell>{allocation.subject}</TableCell>
                    <TableCell>{allocation.course_code}</TableCell>
                    <TableCell>{allocation.section}</TableCell>
                    <TableCell>{allocation.faculty}</TableCell>
                  </TableRow>
                ))}
              </React.Fragment>
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={allocations.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </TableContainer>
    );
  };

  return (
    <AppLayout
      rId={1}
      title="Dashboard"
      body={
        <>
       
        
        </>
      }
    />
  );
};

export default ExcelUpload;
