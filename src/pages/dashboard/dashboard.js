import React, { useEffect, useState } from 'react';
import AppLayout from '../../layout/layout';
import axios from 'axios';
import './styles.css'; 

function Dashboard() {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage] = useState(5); 

  useEffect(() => {
  
    axios.get('http://localhost:8080/saved/deptoptions')
      .then((response) => {
        setData(response.data);
        setFilteredData(response.data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }, []);

  useEffect(() => {

    const results = data.filter((item) =>
      item.department_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.semester_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.classroom.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredData(results);
    setCurrentPage(1); 
  }, [searchTerm, data]);

  const handleDelete = (index) => {
    // Handle deletion of a row
    const newData = [...filteredData];
    newData.splice(index, 1);
    setFilteredData(newData);
    setData(newData); // Also update the original data if needed
  };

  // Pagination logic
  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentRows = filteredData.slice(indexOfFirstRow, indexOfLastRow);
  const totalPages = Math.ceil(filteredData.length / rowsPerPage);

  const handleClick = (pageNumber) => {
    setCurrentPage(pageNumber);
  };

  return (
    <AppLayout
      rId={1}
      title="Dashboard"
      body={
        <>
          {/* <div className="table-header">
            <div className="search-bar">
              <input
                type="text"
                placeholder="Search…"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="search-input"
              />
            </div>
          </div>
          <div className="custom-table">
            <div className="custom-table-body">
              <table>
                <thead>
                  <tr>
                    <td>S.No</td>
                    <td>Department</td>
                    <td>Semester</td>
                    <td>Classroom</td>
                    <td>Action</td>
                  </tr>
                </thead>
                <tbody>
                  {currentRows.map((row, index) => (
                    <tr key={index}>
                      <td>{indexOfFirstRow + index + 1}</td>
                      <td>{row.department_name}</td>
                      <td>{row.semester_name}</td>
                      <td>{row.classroom}</td>
                      <td>
                        <button onClick={() => handleDelete(indexOfFirstRow + index)}>
                          <span className="delete-icon">🗑️</span>
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <div className="pagination">
              <div className="pagination-right">
                {Array.from({ length: totalPages }, (_, index) => (
                  <button
                    key={index}
                    onClick={() => handleClick(index + 1)}
                    className={currentPage === index + 1 ? 'active' : ''}
                  >
                    {index + 1}
                  </button>
                ))}
              </div>
            </div>
          </div> */}
        </>
      }
    />
  );
}

export default Dashboard;
