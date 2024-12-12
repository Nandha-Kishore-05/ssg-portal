import React, { useState } from 'react';
import AppLayout from '../../layout/layout';

const ExcelUpload = () => {
  const [file, setFile] = useState(null);
  const [allocations, setAllocations] = useState([]); // State to store the table data

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
      const response = await fetch("http://localhost:8080/upload", {
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

  // Render the table from allocations
  const renderTable = () => {
    if (allocations.length === 0) {
      return <p>No data available</p>;
    }

    return (
      <table>
        <thead>
          <tr>
            <th>Department</th>
            <th>Date</th>
            <th>Period</th>
            <th>Venue</th>
            <th>Subject</th>
            <th>Course Code</th>
            <th>Section</th>
            <th>Faculty</th>
          </tr>
        </thead>
        <tbody>
          {allocations.map((departmentAlloc, index) => (
            <React.Fragment key={index}>
              {departmentAlloc.allocations.map((allocation, idx) => (
                <tr key={idx}>
                  <td>{departmentAlloc.department}</td>
                  <td>{allocation.date}</td>
                  <td>{allocation.period}</td>
                  <td>{allocation.venue}</td>
                  <td>{allocation.subject}</td>
                  <td>{allocation.course_code}</td>
                  <td>{allocation.section}</td>
                  <td>{allocation.faculty}</td>
                </tr>
              ))}
            </React.Fragment>
          ))}
        </tbody>
      </table>
    );
  };

  return (
    <AppLayout
      rId={1}
      title="Dashboard"
      body={
        <>
          <form onSubmit={handleSubmit}>
            <input type="file" onChange={handleFileChange} />
            <button type="submit">Upload Excel</button>
          </form>
          <div>{renderTable()}</div>
        </>
      }
    />
  );
};

export default ExcelUpload;
