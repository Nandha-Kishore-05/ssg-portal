



import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CustomButton from '../../components/button';
import { Drawer, Box, Typography, List, ListItem, Button, TextField, InputAdornment, IconButton, Grid } from '@mui/material';
import { ListItemAvatar, Avatar, ListItemText, Divider } from '@mui/material';
import PersonIcon from '@mui/icons-material/Person';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import SearchIcon from '@mui/icons-material/Search';
import './save.css';
import { Modal, Dialog, DialogTitle, DialogActions, DialogContent } from '@mui/material';
import ExcelJS from 'exceljs';
import { saveAs } from 'file-saver';
import CloseIcon from '@mui/icons-material/Close';
import SubjectIcon from '@mui/icons-material/MenuBook';

import RoomIcon from '@mui/icons-material/Room';
import AccessTimeIcon from '@mui/icons-material/AccessTime'

const SavedTimetable = (props) => {

  const [schedule, setSchedule] = useState([]);
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);
  const [venue, setVenue] = useState('');
  const [selectedPeriod, setSelectedPeriod] = useState(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [availableFaculty, setAvailableFaculty] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [facultyName, setFacultyName] = useState('');
  const [day, setDay] = useState('');
  const [availableTimings, setAvailableTimings] = useState([]);
  const [editTypeModalOpen, setEditTypeModalOpen] = useState(false);
const [editType, setEditType] = useState(null); // 'swap' or 'manual'
const [manualEditModalOpen, setManualEditModalOpen] = useState(false);
const [selectedPeriodDetails, setSelectedPeriodDetails] = useState(null);


const modalStyle = {
  position: 'absolute',
  top: '45%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 450,
  bgcolor: 'rgba(255, 255, 255, 0.9)',
  borderRadius: '16px',
  boxShadow: '0 8px 24px rgba(0, 0, 0, 0.2)',
  backdropFilter: 'blur(12px)',
  p:6
}



  useEffect(() => {
    const fetchSchedule = async () => {
      if (!props.departmentID || !props.semesterID || !props.academicYearID) {
        console.error('Department ID and Semester ID are required');
        return;
      }

      try {
        const response = await axios.get(`http://localhost:8080/timetable/saved/${props.departmentID}/${props.semesterID}/${props.academicYearID}/${props.sectionID}`);
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
  }, [props.departmentID, props.semesterID,props.academicYearID,props.sectionID]);

  const fetchAvailableFaculty = async (day, time) => {
    try {
      const [startTime, endTime] = time.split(' - ');
      const response = await axios.get(`http://localhost:8080/faculty/available/${props.departmentID}/${props.semesterID}/${day}/${startTime}/${endTime}/${props.academicYearID}/${props.sectionID}`);
      setAvailableFaculty(response.data);
    } catch (error) {
      console.error('Error fetching available faculty:', error);
    }
  };

  const fetchAvailableTimings = async (day,faculty) => {
    try {
      const response = await axios.get(`http://localhost:8080/available-timings/${faculty}/${day}`);
      setAvailableTimings(response.data);
      console.log(response.data)
    } catch (error) {
      console.error('Error fetching available timings:', error);
      toast.error('Failed to fetch available timings');
    }
  };

  // const handleOpenDrawer = (day, time,faculty) => {
  //   if (!isEditMode) return;
  //   setSelectedPeriod({ day, time });
  //   setDrawerOpen(true);
  //   fetchAvailableFaculty(day, time);
  //   fetchAvailableTimings(day,faculty)
  //   console.log("Selected faculty:", faculty);

  // };

  const handleCloseDrawer = () => {
    setDrawerOpen(false);
    setSelectedPeriod(null);
    setAvailableFaculty([]);
    setAvailableTimings([]);
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

        setTimeout(() => {
          window.location.reload();  // Refresh the page immediately
        }, 0); 
        props.setIsOpen(false);  // No delay (effectively microseconds)
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

  const handleEditTimetableClick = () => {
  setEditTypeModalOpen(true);
};

const handleSelectEditType = (type) => {
  setEditType(type);
  setEditTypeModalOpen(false);
  setIsEditMode(true);
};

const handleOpenDrawer = (day, time, faculty) => {
  if (!isEditMode) return;

  if (editType === 'manual') {
    const periodData = schedule.find(item => item.day_name === day && `${item.start_time} - ${item.end_time}` === time);
    if (periodData) {
      setSelectedPeriodDetails(periodData);
      setManualEditModalOpen(true);
    }
  } else {
    // existing drawer logic for swap
    setSelectedPeriod({ day, time });
    setDrawerOpen(true);
    fetchAvailableFaculty(day, time);
    fetchAvailableTimings(day, faculty);
  }
};



  return (
    <div className="container-3">
       <div className="header-i">
       <div className="header-info">
          <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : S{props.semesterID} </h2>
        
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
  onClick={isEditMode ? handleSaveTimetable : handleEditTimetableClick}
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
      <td className="day">{day}</td>
      {times.map((time, index) => {
        const classes = schedule.filter(
          (item) => item.day_name === day && `${item.start_time} - ${item.end_time}` === time
        );
        const isActive = selectedPeriod?.day === day && selectedPeriod?.time === time;
        const faculty = classes.length > 0 ? classes.map((item) => item.faculty_name) : null;

        // Match available timings
        const isAvailableTiming = availableTimings.some(
          (timing) =>
            timing.day_name === day &&
            timing.start_time === time.split(' - ')[0] &&
            timing.end_time === time.split(' - ')[1]
        );

        console.log({
          day,
          time,
          startTime: time.split(' - ')[0],
          endTime: time.split(' - ')[1],
          isAvailableTiming,
        });

        const cellStyle = {
          // border: '2px solid #ddd',
          // padding: '7px',
          // fontSize: '16px',
          // fontWeight: 'bold',
          // textAlign: 'center',
          cursor: isEditMode ? 'pointer' : 'default',
          backgroundColor: isAvailableTiming
            ? '#dff0d8' // Light green if timing is available
            : isActive
            ? '#f0e68c' // Light yellow if active
            : '#fff', // Default white background
        };

        return (
          <td
            key={index}
            style={cellStyle}
            onClick={() => handleOpenDrawer(day, time, faculty)}
          >
            {classes.length > 0 ? (
              classes.map((item, idx) => (
                <div key={idx}>
                  <div>{item.subject_name}</div>
                  <div>{item.faculty_name}</div>
                  <div>{item.classroom}</div>
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
          borderBottom: '1px solid #cccc',
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
     <Modal
  open={editTypeModalOpen}
  onClose={() => setEditTypeModalOpen(false)}
  aria-labelledby="edit-type-modal-title"
  aria-describedby="edit-type-modal-description"
>
  <Box sx={modalStyle}>
    {/* Close icon */}
    <IconButton
      onClick={() => setEditTypeModalOpen(false)}
      sx={{
        position: 'absolute',
        top: 16,
        right: 16,
        color: '#444',
        '&:hover': { color: '#000' },
      }}
    >
      <CloseIcon fontSize="medium" />
    </IconButton>

    {/* Title */}

   
   

    {/* Subtitle */}
     <center>
                   
                      <h3 style={{marginTop:'20px',marginBottom:'20px'}}> Are you sure you want to edit the timetable?</h3>
                  
                  </center>

    {/* Buttons */}
    <Box sx={{ display: 'flex', justifyContent: 'center', gap: 3 }}>
      <CustomButton
        label="Swap"
        backgroundColor="#1976d2"
        onClick={() => handleSelectEditType('swap')}
        style={{
          width: '160px',
          height: '50px',
          backgroundColor: '#E9F5FE',
          border: '2px solid #1976d2',
          color: '#1976d2',
          fontWeight: '600',
          borderRadius: '10px',
          transition: '0.3s',
          boxShadow: '0 2px 6px rgba(25, 118, 210, 0.2)',
        }}
        hoverStyle={{
          backgroundColor: '#0C7FDA',
          color: '#fff',
        }}
      />
      <CustomButton
        label="Edit Manually"
        onClick={() => handleSelectEditType('manual')}
        style={{
          width: '160px',
          height: '50px',
          backgroundColor: '#1976d2',
          color: '#fff',
          fontWeight: '600',
          borderRadius: '10px',
          boxShadow: '0 2px 8px rgba(25, 118, 210, 0.3)',
          transition: '0.3s',
        }}
        hoverStyle={{
          backgroundColor: '#0d47a1',
        }}
      />
    </Box>
  </Box>
</Modal>

<Dialog
  open={manualEditModalOpen}
  onClose={() => setManualEditModalOpen(false)}
  fullWidth
  maxWidth="sm"
  PaperProps={{
    sx: {
      bgcolor: '#ffffff',
      borderRadius: 4,
      boxShadow: '0 10px 30px rgba(0,0,0,0.15)',
      fontFamily: 'inherit',
    },
  }}
>
  <Box
    sx={{
      backgroundColor: '#ffffff',
      color: 'black',
      fontWeight: 700,
      fontSize: '1.6rem',
      borderBottom: '1px solid #e0e0e0',
      px: 4,
      py: 2.5,
      textAlign: 'center',
    }}
  >
    Edit Period Details
  </Box>
<br />
  <DialogContent sx={{ px: 4, py: 3 }}>
    {selectedPeriodDetails && (
      <Grid container spacing={3}>
        {[
          { label: 'Subject', name: 'subject_name' },
          { label: 'Faculty', name: 'faculty_name' },
          { label: 'Venue', name: 'classroom' },
        ].map((field) => (
          <Grid item xs={12} key={field.name}>
            <TextField
              fullWidth
              label={field.label}
              name={field.name}
              value={selectedPeriodDetails[field.name]}
              onChange={(e) =>
                setSelectedPeriodDetails({
                  ...selectedPeriodDetails,
                  [field.name]: e.target.value,
                })
              }
              variant="outlined"
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: 2,
                  fontSize: '1rem',
                  px: 1.5,
                  fontFamily: 'inherit',
                  boxShadow: 'inset 0 1px 3px rgba(0,0,0,0.1)',
                },
                '& .MuiInputLabel-root': {
                  fontWeight: 800,
                  fontSize:'1.2rem',
                  color: '#555',
                  fontFamily: 'inherit',
                },
              }}
            />
          </Grid>
        ))}

        {[
          { label: 'Start Time', name: 'start_time' },
          { label: 'End Time', name: 'end_time' },
        ].map((field) => (
          <Grid item xs={6} key={field.name}>
            <TextField
              fullWidth
              type="time"
              label={field.label}
              name={field.name}
              value={selectedPeriodDetails[field.name]}
              onChange={(e) =>
                setSelectedPeriodDetails({
                  ...selectedPeriodDetails,
                  [field.name]: e.target.value,
                })
              }
              InputLabelProps={{ shrink: true }}
              variant="outlined"
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: 2,
                  fontSize: '1rem',
                  px: 1.5,
                  fontFamily: 'inherit',
                  boxShadow: 'inset 0 1px 3px rgba(0,0,0,0.1)',
                },
                '& .MuiInputLabel-root': {
                  fontWeight: 800,
                     fontSize:'1.2rem',
                  color: '#555',
                  fontFamily: 'inherit',
                },
              }}
            />
          </Grid>
        ))}
      </Grid>
    )}
  </DialogContent>

 <DialogActions sx={{ px: 4, pb: 3, pt: 2 }}>
    <Button
      onClick={() => setManualEditModalOpen(false)}
      variant="outlined"
      color="error"
      sx={{
        borderRadius: 2,
        px: 3.5,
        py: 1.2,
        fontWeight: 600,
        fontSize: '1rem',
        fontFamily: 'inherit',
        textTransform: 'none',
        lineHeight: 1.4,
        '&:hover': {
          backgroundColor: '#ffe6e6',
        },
      }}
    >
      Cancel
    </Button>
    <Button
      onClick={() => {
        console.log('Saved:', selectedPeriodDetails);
        toast.success('Period updated successfully!');
        setManualEditModalOpen(false);
      }}
      variant="contained"
      color="primary"
      sx={{
        borderRadius: 2,
        px: 3.5,
        py: 1.2,
        fontWeight: 600,
        fontSize: '1rem',
        fontFamily: 'inherit',
        textTransform: 'none',
        lineHeight: 1.4,
        boxShadow: '0 3px 6px rgba(0,0,0,0.1)',
      }}
    >
      Save Changes
    </Button>
  </DialogActions>

  <ToastContainer position="bottom-center" autoClose={2500} />
</Dialog>



    </div>
  );
};

export default SavedTimetable;
