

//   import React, { useState, useEffect } from 'react';
// import axios from 'axios';
// import CustomButton from '../../components/button';
// import * as XLSX from 'xlsx';
// import { Modal, Box, Typography, Button, Select, MenuItem, InputLabel, FormControl } from '@mui/material';
// import './save.css'; 

// const SavedTimetable = (props) => {
//   const [schedule, setSchedule] = useState([]);
//   const [days, setDays] = useState([]);
//   const [times, setTimes] = useState([]);
//   const [venue, setVenue] = useState('');
//   const [selectedPeriod, setSelectedPeriod] = useState(null);
//   const [isModalOpen, setIsModalOpen] = useState(false);
//   const [isEditMode, setIsEditMode] = useState(false);
//   const [availableFaculty, setAvailableFaculty] = useState([]);
//   const [selectedFaculty, setSelectedFaculty] = useState('');

//   useEffect(() => {
//     const fetchSchedule = async () => {
//       if (!props.departmentID || !props.semesterID) {
//         console.error('Department ID and Semester ID are required');
//         return;
//       }

//       try {
//         const response = await axios.get(`http://localhost:8080/timetable/saved/${props.departmentID}/${props.semesterID}`);
//         const data = response.data;

//         const allDays = new Set();
//         const allTimes = new Set();
//         let venueSet = new Set();

//         data.forEach(item => {
//           allDays.add(item.day_name);
//           allTimes.add(`${item.start_time} - ${item.end_time}`);
//           venueSet.add(item.classroom);
//         });

//         const sortedDays = Array.from(allDays).sort((a, b) => {
//           const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
//           return order.indexOf(a) - order.indexOf(b);
//         });

//         const sortedTimes = Array.from(allTimes).sort((a, b) => a.localeCompare(b, undefined, { numeric: true }));

//         setDays(sortedDays);
//         setTimes(sortedTimes);
//         setSchedule(data);
//         setVenue(Array.from(venueSet).join(', '));
//       } catch (error) {
//         console.error('Error fetching timetable:', error);
//       }
//     };

//     fetchSchedule();
//   }, [props.departmentID, props.semesterID]);

//   const fetchAvailableFaculty = async () => {
//     if (!selectedPeriod) return;

//     try {
//       const response = await axios.get(`http://localhost:8080/faculty/available/${props.departmentID}/${props.semesterID}/${selectedPeriod.day}/${selectedPeriod.start_time}/${selectedPeriod.end_time}`);
//       setAvailableFaculty(response.data);
//     } catch (error) {
//       console.error('Error fetching available faculty:', error);
//     }
//   };

//   const handleOpenModal = (day, time) => {
//     if (!isEditMode) return;

//     const [start_time, end_time] = time.split(' - ');
//     setSelectedPeriod({ day, start_time, end_time });
//     fetchAvailableFaculty();
//     setIsModalOpen(true);
//   };

//   const handleCloseModal = () => {
//     setIsModalOpen(false);
//     setSelectedPeriod(null);
//     setAvailableFaculty([]);
//     setSelectedFaculty('');
//   };

//   const handleConfirmRemove = () => {
//     const updatedSchedule = schedule.filter(item =>
//       !(item.day_name === selectedPeriod.day && `${item.start_time} - ${item.end_time}` === `${selectedPeriod.start_time} - ${selectedPeriod.end_time}`)
//     );
//     setSchedule(updatedSchedule);
//     handleCloseModal();
//   };

//   const handleDownload = () => {
//     const workbook = XLSX.utils.book_new();
//     const worksheetData = [
//       ['Day/Time', ...times]
//     ];
  
//     // Add rows to the worksheet data
//     days.forEach(day => {
//       const row = [day];
//       times.forEach(time => {
//         const classes = schedule.filter(item =>
//           item.day_name === day && `${item.start_time} - ${item.end_time}` === time
//         );
  
//         const classInfo = classes.map(item => `${item.subject_name} (${item.faculty_name})`).join('\n');
//         row.push(classInfo || 'No classes');
//       });
//       worksheetData.push(row);
//     });
  
//     const worksheet = XLSX.utils.aoa_to_sheet(worksheetData);
  
//     // Define styles
//     const headerStyle = {
//       fill: { fgColor: { rgb: "343a40" } }, // Background color for header
//       font: { color: { rgb: "ffffff" }, bold: true, sz: 14 }, // Font color and size for header
//       alignment: { horizontal: "center", vertical: "center" },
//       border: { 
//         top: { style: 'thin', color: { rgb: '000000' } },
//         bottom: { style: 'thin', color: { rgb: '000000' } },
//         left: { style: 'thin', color: { rgb: '000000' } },
//         right: { style: 'thin', color: { rgb: '000000' } }
//       }
//     };
  
//     const timeHeaderStyle = {
//       fill: { fgColor: { rgb: "6c757d" } }, // Background color for time headers
//       font: { color: { rgb: "ffffff" }, bold: true, sz: 14 }, // Font color and size for time headers
//       alignment: { horizontal: "center", vertical: "center" },
//       border: { 
//         top: { style: 'thin', color: { rgb: '000000' } },
//         bottom: { style: 'thin', color: { rgb: '000000' } },
//         left: { style: 'thin', color: { rgb: '000000' } },
//         right: { style: 'thin', color: { rgb: '000000' } }
//       }
//     };
  
//     const dayCellStyle = {
//       fill: { fgColor: { rgb: "6c757d" } }, // Background color for day cells
//       font: { color: { rgb: "ffffff" }, bold: true, sz: 12 }, // Font color and size for day cells
//       alignment: { horizontal: "center", vertical: "center" },
//       border: { 
//         top: { style: 'thin', color: { rgb: '000000' } },
//         bottom: { style: 'thin', color: { rgb: '000000' } },
//         left: { style: 'thin', color: { rgb: '000000' } },
//         right: { style: 'thin', color: { rgb: '000000' } }
//       }
//     };
  
//     const subjectCellStyle = {
//       border: { 
//         top: { style: 'thin', color: { rgb: '000000' } },
//         bottom: { style: 'thin', color: { rgb: '000000' } },
//         left: { style: 'thin', color: { rgb: '000000' } },
//         right: { style: 'thin', color: { rgb: '000000' } }
//       }
//     };
  
//     // Apply styles to header
//     worksheet['A1'].s = headerStyle;
//     times.forEach((_, index) => {
//       worksheet[XLSX.utils.encode_cell({ r: 0, c: index + 1 })].s = timeHeaderStyle;
//     });
  
//     // Apply styles to each row
//     worksheetData.forEach((row, rowIndex) => {
//       row.forEach((_, colIndex) => {
//         const cellAddress = XLSX.utils.encode_cell({ r: rowIndex, c: colIndex });
//         if (rowIndex === 0) {
//           // Apply header style to the first row
//           worksheet[cellAddress].s = headerStyle;
//         } else if (colIndex === 0) {
//           // Apply dayCellStyle to the first column
//           worksheet[cellAddress].s = dayCellStyle;
//         } else {
//           // Apply subjectCellStyle to all other cells
//           worksheet[cellAddress].s = subjectCellStyle;
//         }
//       });
//     });
  
//     // Define column widths
//     worksheet['!cols'] = [{ width: 20 }].concat(times.map(() => ({ width: 30 })));
  
//     XLSX.utils.book_append_sheet(workbook, worksheet, `Semester ${props.semesterID}`);
//     XLSX.writeFile(workbook, `timetable_semester_${props.semesterID}.xlsx`);
//   };
  

  
//   const handleToggleEditMode = () => {
//     setIsEditMode(!isEditMode);
//   };

//   return (
//     <div className="container-1">
//       <div className="header-k">
//         <div className="header-info">
//           <h2>Semester : S{props.semesterID}</h2>
//           <h2>Venue: {venue || 'Not Available'}</h2>
//         </div>
//         <div className="buttons">
//           <CustomButton
//             width="150"
//             label="Download Timetable"
//             onClick={handleDownload}
//           />
//           <CustomButton
//             width="150"
//             label={isEditMode ? "Save Edited Timetable" : "Edit Timetable"}
//             backgroundColor={isEditMode ? "green" : "red"}
//             onClick={handleToggleEditMode}
//           />
//         </div>
//       </div>
//       <table className="table">
//         <thead>
//           <tr>
//             <th className="day-time">Day/Time</th>
//             {times.map((time, index) => (
//               <th key={index} className="time">
//                 {time}
//               </th>
//             ))}
//           </tr>
//         </thead>
//         <tbody>
//           {days.map((day) => (
//             <tr key={day}>
//               <td className="day">{day}</td>
//               {times.map((time, index) => {
//                 const classes = schedule.filter(item =>
//                   item.day_name === day && `${item.start_time} - ${item.end_time}` === time
//                 );
//                 return (
//                   <td
//                     key={index}
//                     className="subject"
//                     onClick={() => isEditMode && handleOpenModal(day, time)}
//                   >
//                     {classes.length > 0 ? (
//                       classes.map((item, idx) => (
//                         <div key={idx}>
//                           <div>{item.subject_name}</div>
//                           <div>{item.faculty_name}</div>
//                         </div>
//                       ))
//                     ) : (
//                       isEditMode && (
//                         <Button
//                           variant="contained"
//                           color="primary"
//                           onClick={() => {
//                             setSelectedPeriod({ day, start_time: time.split(' - ')[0], end_time: time.split(' - ')[1] });
//                             fetchAvailableFaculty();
//                             setIsModalOpen(true);
//                           }}
//                         >
//                           Show Faculty
//                         </Button>
//                       )
//                     )}
//                     {availableFaculty.length > 0 && (
//                       <FormControl fullWidth sx={{ mt: 2 }}>
//                         <InputLabel id="select-faculty-label">Select Faculty</InputLabel>
//                         <Select
//                           labelId="select-faculty-label"
//                           value={selectedFaculty}
//                           onChange={(e) => setSelectedFaculty(e.target.value)}
//                           displayEmpty
//                         >
//                           <MenuItem value="" disabled>Select Faculty</MenuItem>
//                           {availableFaculty.map(faculty => (
//                             <MenuItem key={faculty.id} value={faculty.name}>
//                               {faculty.name}
//                             </MenuItem>
//                           ))}
//                         </Select>
//                       </FormControl>
//                     )}
//                   </td>
//                 );
//               })}
//             </tr>
//           ))}
//         </tbody>
//       </table>

//       <Modal
//         open={isModalOpen}
//         onClose={handleCloseModal}
//       >
//         <Box sx={{
//           position: 'absolute',
//           top: '50%',
//           left: '50%',
//           transform: 'translate(-50%, -50%)',
//           width: 400,
//           bgcolor: 'background.paper',
//           border: '2px solid #000',
//           boxShadow: 24,
//           p: 4,
//         }}>
//           <Typography variant="h6" component="h2">
//             Confirm Removal
//           </Typography>
//           <Typography sx={{ mt: 2 }}>
//             Are you sure you want to remove the period on {selectedPeriod?.day} during {selectedPeriod?.start_time} - {selectedPeriod?.end_time}?
//           </Typography>
//           {availableFaculty.length > 0 && (
//             <FormControl fullWidth sx={{ mt: 2 }}>
//               <InputLabel id="select-faculty-label">Select Faculty</InputLabel>
//               <Select
//                 labelId="select-faculty-label"
//                 value={selectedFaculty}
//                 onChange={(e) => setSelectedFaculty(e.target.value)}
//                 displayEmpty
//               >
//                 <MenuItem value="" disabled>Select Faculty</MenuItem>
//                 {availableFaculty.map(faculty => (
//                   <MenuItem key={faculty.id} value={faculty.name}>
//                     {faculty.name}
//                   </MenuItem>
//                 ))}
//               </Select>
//             </FormControl>
//           )}
//           <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
//             <Button variant="contained" color="primary" onClick={handleConfirmRemove}>
//               Confirm
//             </Button>
//           </Box>
//         </Box>
//       </Modal>
//     </div>
//   );
// };

// export default SavedTimetable;



import React, { useState, useEffect } from 'react';
import axios from 'axios';
import * as XLSX from 'xlsx';
import CustomButton from '../../components/button';
import { Drawer, Box, Typography, List, ListItem, Button, TextField, InputAdornment } from '@mui/material';
import { ListItemAvatar, Avatar, ListItemText, Divider } from '@mui/material';
import PersonIcon from '@mui/icons-material/Person';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import SearchIcon from '@mui/icons-material/Search';
import './save.css';

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

  const handleDownload = () => {
    const workbook = XLSX.utils.book_new();
    const worksheetData = [
      ['Day/Time', ...times]
    ];

    days.forEach(day => {
      const row = [day];
      times.forEach(time => {
        const classes = schedule.filter(item =>
          item.day_name === day && `${item.start_time} - ${item.end_time}` === time
        );

        const classInfo = classes.map(item => `${item.subject_name} (${item.faculty_name})`).join('\n');
        row.push(classInfo || 'No classes');
      });
      worksheetData.push(row);
    });

    const worksheet = XLSX.utils.aoa_to_sheet(worksheetData);
    XLSX.utils.book_append_sheet(workbook, worksheet, `Semester ${props.semesterID}`);
    
    XLSX.writeFile(workbook, `timetable_semester_${props.semesterID}.xlsx`);
  };

  const handleToggleEditMode = () => {
    setIsEditMode(!isEditMode);
  };

  const handleSaveTimetable = async () => {
    console.log("Attempting to save timetable:", schedule);
    try {
      const response = await axios.put('http://localhost:8080/timetable/update', schedule);
      console.log("Response from server:", response.data);
      toast.success('Timetable updated successfully!');
    } catch (error) {
      console.error('Failed to update timetable:', error);
      toast.error('Failed to update timetable');
    }
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
                }}
                secondaryTypographyProps={{
                  color: '#888',
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
      <Button variant="contained" onClick={handleCloseDrawer} color="secondary">
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
