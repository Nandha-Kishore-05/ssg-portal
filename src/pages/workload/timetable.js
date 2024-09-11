



import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CustomButton from '../../components/button';
import { Drawer, Box, Typography, List, ListItem, Button, TextField, InputAdornment } from '@mui/material';
import { ListItemAvatar, Avatar, ListItemText, Divider } from '@mui/material';
import PersonIcon from '@mui/icons-material/Person';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import SearchIcon from '@mui/icons-material/Search';
import './save.css';
import { useNavigate } from 'react-router-dom';
import ExcelJS from 'exceljs';
import { saveAs } from 'file-saver';

const SavedTimetable = (props) => {
  const navigate = useNavigate();
  const [schedule, setSchedule] = useState([]);
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);
  const [venue, setVenue] = useState('');
  const [selectedPeriod, setSelectedPeriod] = useState(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [availableFaculty, setAvailableFaculty] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    const fetchSchedule = async () => {
      if (!props.departmentID || !props.semesterID) {
        console.error('Department ID and Semester ID are required');
        return;
      }

      try {
        const response = await axios.get(`http://localhost:8080/timetable/saved/${props.departmentID}/${props.semesterID}`);
        const data = response.data;

        const allDays = new Set();
        const allTimes = new Set();
        let venueSet = new Set();

        data.forEach(item => {
          allDays.add(item.day_name);
          allTimes.add(`${item.start_time} - ${item.end_time}`);
          venueSet.add(item.classroom);
        });

        const sortedDays = Array.from(allDays).sort((a, b) => {
          const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
          return order.indexOf(a) - order.indexOf(b);
        });

        const sortedTimes = Array.from(allTimes).sort((a, b) => a.localeCompare(b, undefined, { numeric: true }));

        setDays(sortedDays);
        setTimes(sortedTimes);
        setSchedule(data);
        setVenue(Array.from(venueSet).join(', '));
      } catch (error) {
        console.error('Error fetching timetable:', error);
      }
    };

    fetchSchedule();
  }, [props.departmentID, props.semesterID]);

  const fetchAvailableFaculty = async (day, time) => {
    try {
      const [startTime, endTime] = time.split(' - ');
      const response = await axios.get(`http://localhost:8080/faculty/available/${props.departmentID}/${props.semesterID}/${day}/${startTime}/${endTime}`);
      setAvailableFaculty(response.data);
    } catch (error) {
      console.error('Error fetching available faculty:', error);
    }
  };

  const handleOpenDrawer = (day, time) => {
    if (!isEditMode) return;
    setSelectedPeriod({ day, time });
    setDrawerOpen(true);
    fetchAvailableFaculty(day, time);
  };

  const handleCloseDrawer = () => {
    setDrawerOpen(false);
    setSelectedPeriod(null);
    setAvailableFaculty([]);
  };

  const handleConfirmAssignFaculty = async (faculty) => {
    try {
      const updatedSchedule = schedule.map(item => {
        if (item.day_name === selectedPeriod.day && `${item.start_time} - ${item.end_time}` === selectedPeriod.time) {
          return { ...item, faculty_name: faculty.name, subject_name: faculty.subject_name };
        }
        return item;
      });

      setSchedule(updatedSchedule);
      handleCloseDrawer();
      toast.success(`Assigned ${faculty.name} for ${selectedPeriod.day} at ${selectedPeriod.time}`);
    } catch (error) {
      console.error('Failed to assign faculty:', error);
      toast.error('Failed to assign faculty');
    }
  };
  const handleDownload = async () => {
    const workbook = new ExcelJS.Workbook();
    const worksheet = workbook.addWorksheet(`Semester ${props.semesterID}`);
  
  
    worksheet.columns = [
      { header: 'Day/Time', key: 'day', width: 30 }, 
      ...times.map(time => ({ header: time, key: time, width: 30 })), 
    ];
  
 
    days.forEach(day => {
      const rowData = { day };
      times.forEach(time => {
        const classes = schedule.filter(item =>
          item.day_name === day && `${item.start_time} - ${item.end_time}` === time
        );
        const classInfo = classes.map(item => `${item.subject_name}\n${item.faculty_name}`).join('\n'); 
        rowData[time] = classInfo || 'No classes';
      });
      worksheet.addRow(rowData);
    });
  

    worksheet.eachRow(row => {
      row.height = 65; 
    });
  

    const headerRow = worksheet.getRow(1);
    headerRow.eachCell(cell => {
      cell.font = { 
        bold: true, 
        name: 'Segoe UI Variable Display Semib', 
        size: 12,
        color: { argb: 'FFEFEFEF' } 
      };
      cell.fill = {
        type: 'pattern',
        pattern: 'solid',
        fgColor: { argb: 'FF6C757D' } 
      };
      cell.alignment = { horizontal: 'center', vertical: 'middle' };
      cell.border = {
        top: { style: 'thin', color: { argb: 'FF000000' } }, 
        left: { style: 'thin', color: { argb: 'FF000000' } },
        bottom: { style: 'thin', color: { argb: 'FF000000' } },
        right: { style: 'thin', color: { argb: 'FF000000' } },
      };
    });
  
  
    worksheet.eachRow((row, rowNumber) => {
      row.eachCell(cell => {
        cell.font = {
          name: 'Segoe UI Variable Display Semib', 
          size: 12, 
          bold: true
        };
        cell.alignment = { horizontal: 'center', vertical: 'middle', wrapText: true };
        cell.border = {
          top: { style: 'thin', color: { argb: 'FF000000' } }, 
          left: { style: 'thin', color: { argb: 'FF000000' } },
          bottom: { style: 'thin', color: { argb: 'FF000000' } },
          right: { style: 'thin', color: { argb: 'FF000000' } },
        };
  
        if (rowNumber > 1) {
      
          if (rowNumber % 2 === 0) {
            cell.fill = {
              type: 'pattern',
              pattern: 'solid',
              fgColor: { argb: 'FFF8F9FA' } 
            };
          }
        }
      });
    });
  

    const buffer = await workbook.xlsx.writeBuffer();
    const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    saveAs(blob, `timetable_semester_${props.semesterID}.xlsx`);
  };
  
  const handleToggleEditMode = () => {
    setIsEditMode(!isEditMode);
  };

  const handleSaveTimetable = async () => {
    console.log("Attempting to save timetable:", schedule);
  
    await axios.put('http://localhost:8080/timetable/update', schedule)
      .then(response => {
        console.log("Response from server:", response.data);
        toast.success('Timetable updated successfully!');
      props.setIsOpen(false)
      })
      .catch(error => {
        console.error('Failed to update timetable:', error);
        toast.error('Failed to update timetable');
      });
  };
  

  const filteredFaculty = availableFaculty.filter(faculty =>
    faculty.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    faculty.subject_name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="container-1">
       <div className="header-k">
       <div className="header-info">
          <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : S{props.semesterID} </h2>
          <h2 style={{ fontSize: '20px', marginTop: '5px', marginLeft: '15px' }}>Venue: {venue || 'Not Available'}</h2>
        </div>
        <div className="buttons">
          <CustomButton
            width="150"
            label="Download Timetable"
            onClick={handleDownload}
          />
          <CustomButton
            width="150" 
            label={isEditMode ? "Save Edited Timetable" : "Edit Timetable"}
            backgroundColor={isEditMode ? "green" : "red"}
            onClick={isEditMode ? handleSaveTimetable : handleToggleEditMode}
          />
        </div>
      </div>
      <table className="table">
        <thead>
          <tr>
            <th className="day-time" >Day/Time</th>
            {times.map((time, index) => (
              <th key={index} className="time" >
                {time}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {days.map((day) => (
            <tr key={day}>
              <td className="day" >{day}</td>
              {times.map((time, index) => {
                const classes = schedule.filter(item =>
                  item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                );
                const isActive = selectedPeriod?.day === day && selectedPeriod?.time === time;
                return (
                  <td
                    key={index}
                    style={{
                      border: '2px solid #ddd',
                      padding: '7px',
                      fontSize: '16px',
                      fontWeight: 'bold',
                      textAlign: 'center',
                      cursor: isEditMode ? 'pointer' : 'default',
                      backgroundColor: isActive ? '#dff0d8' : '#fff'
                    }}
                    onClick={() => handleOpenDrawer(day, time)}
                  >
                    {classes.length > 0 ? (
                      classes.map((item, idx) => (
                        <div key={idx}>
                          <div>{item.subject_name}</div>
                          <div>{item.faculty_name}</div>
                        </div>
                      ))
                    ) : (
                      <div>No classes</div>
                    )}
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
      <Drawer
  anchor="left"
  open={drawerOpen}
  onClose={handleCloseDrawer}
  PaperProps={{
    style: {
      width: '350px',
      backgroundColor: '#fff',
      borderTopRightRadius: '20px', // Top right border-radius
      borderBottomRightRadius: '20px', // Bottom right border-radius
      padding: '10px',
      height: '100vh',
      overflowY: filteredFaculty.length > 6 ? 'auto' : 'hidden',
    color:'black'
    },
  }}
>
  <Box><br />
    <Typography variant="h5" gutterBottom style={{ textAlign: 'center',fontWeight:'bolder' }}>
      Select Faculty
    </Typography><br />
    <TextField
      variant="standard"
      fullWidth
      placeholder="Search faculty..."
      value={searchQuery}
      onChange={(e) => setSearchQuery(e.target.value)}
      InputProps={{
        startAdornment: (
          <InputAdornment position="start">
            <SearchIcon />
          </InputAdornment>
        ),
        style: {
          borderBottom: '1px solid #ccc',
          paddingBottom: '5px',
          fontFamily: 'Nunito, sans-serif',
        },
        disableUnderline: true,
      }}
      style={{ marginBottom: '20px' }}
    />
    <List style={{ padding: '0', marginTop: '10px' }}>
      {filteredFaculty.length === 0 ? (
        <Typography variant="body1" style={{ textAlign: 'center', marginTop: '20px' }}>
          No faculty found.
        </Typography>
      ) : (
        filteredFaculty.map((faculty, index) => (
          <React.Fragment key={index}>
            <ListItem
              button
              onClick={() => handleConfirmAssignFaculty(faculty)}
              style={{
                borderRadius: '12px',
                marginBottom: '10px',
                border: '1px solid #ddd',
                boxShadow: '0px 4px 6px rgba(0.3, 0.9, 0.5, 1.5)',
                backgroundColor: '#fefefe',
                transition: 'transform 0.3s, box-shadow 0.3s',
                padding: '10px',
                display: 'flex',
                alignItems: 'center',
                '&:hover': {
                  backgroundColor: '#f0f0f0',
                  transform: 'scale(1.02)',
                  boxShadow: '0px 8px 12px rgba(0, 0, 0, 0.2)',
                   color:'black'
                },
              }}
            >
              <ListItemAvatar>
                <Avatar style={{ backgroundColor: 'gery', color: 'black' }}>
                  <PersonIcon />
                </Avatar>
              </ListItemAvatar>
              <ListItemText
                primary={faculty.name}
                secondary={faculty.subject_name}
                primaryTypographyProps={{
                  fontWeight: 'bold',
                  fontFamily: 'Nunito, sans-serif',
                   color:'black'
                }}
                secondaryTypographyProps={{
                  color: 'grey',
                  fontFamily: 'Nunito, sans-serif',
                }}
              />
            </ListItem>
            {index < filteredFaculty.length - 1 && <Divider style={{ margin: '10px 0' }} />}
          </React.Fragment>
        ))
      )}
    </List>
    <Box mt={2} display="flex" justifyContent="center">
      <Button variant="contained" onClick={handleCloseDrawer}   sx={{ backgroundColor: 'red', '&:hover': { backgroundColor: 'red' } }}>
        Cancel
      </Button>
    </Box>
  </Box>
</Drawer>


      <ToastContainer />
    </div>
  );
};

export default SavedTimetable;
