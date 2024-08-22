
  
// import React, { useState, useEffect } from 'react';
// import axios from 'axios';

// import CustomButton from '../../components/button';
// import * as XLSX from 'xlsx';
// import { Modal, Box, Typography, Button } from '@mui/material';

// const SavedTimetable = (props) => {
//   const [schedule, setSchedule] = useState([]);
//   const [days, setDays] = useState([]);
//   const [times, setTimes] = useState([]);
//   const [venue, setVenue] = useState('');
//   const [selectedPeriod, setSelectedPeriod] = useState(null);
//   const [isModalOpen, setIsModalOpen] = useState(false);
//   const [isEditMode, setIsEditMode] = useState(false);  // Track edit mode

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

//   const handleOpenModal = (day, time) => {
//     if (!isEditMode) return; // Prevent opening modal if not in edit mode
//     setSelectedPeriod({ day, time });
//     setIsModalOpen(true);
//   };

//   const handleCloseModal = () => {
//     setIsModalOpen(false);
//     setSelectedPeriod(null);
//   };

//   const handleConfirmRemove = () => {
//     const updatedSchedule = schedule.filter(item =>
//       !(item.day_name === selectedPeriod.day && `${item.start_time} - ${item.end_time}` === selectedPeriod.time)
//     );
//     setSchedule(updatedSchedule);
//     handleCloseModal();
//   };

//   const handleDownload = () => {
//     const workbook = XLSX.utils.book_new();
//     const worksheetData = [
//       ['Day/Time', ...times]
//     ];

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
//     XLSX.utils.book_append_sheet(workbook, worksheet, `Semester ${props.semesterID}`);
    
//     XLSX.writeFile(workbook, `timetable_semester_${props.semesterID}.xlsx`);
//   };

//   const handleToggleEditMode = () => {
//     setIsEditMode(!isEditMode);
//   };

//   return (
//     <div style={{
//       backgroundColor: '#fff',
//       padding: '20px',
//       borderRadius: '8px',
//       boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)',
//       margin: '20px 0'
//     }}>
//       <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', marginBottom: '13px' }}>
//         <div style={{ display: 'flex', flexDirection: 'row' }}>
//           <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : S{props.semesterID} </h2>
//           <h2 style={{ fontSize: '20px', marginTop: '5px', marginLeft: '15px' }}>Venue: {venue || 'Not Available'}</h2>
//         </div>
//         <div style={{ display: 'flex', flexDirection: 'row', columnGap: 10 }}>
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
//       <table style={{
//         width: '100%',
//         borderCollapse: 'collapse',
//         backgroundColor: '#fff',
//         border: '2px solid #ddd',
//         fontSize: '16px',
//         minHeight: '600px'
//       }}>
//         <thead>
//           <tr style={{ backgroundColor: '#f4f4f4' }}>
//             <th style={{
//               border: '2px solid #ddd',
//               padding: '12px',
//               textAlign: 'center'
//             }}>Day/Time</th>
//             {times.map((time, index) => (
//               <th key={index} style={{
//                 border: '2px solid #ddd',
//                 padding: '12px',
//                 backgroundColor: '#f9f9f9',
//                 textAlign: 'center'
//               }}>
//                 {time}
//               </th>
//             ))}
//           </tr>
//         </thead>
//         <tbody>
//           {days.map((day) => (
//             <tr key={day}>
//               <td style={{
//                 border: '2px solid #ddd',
//                 padding: '12px',
//                 fontWeight: 'bold',
//                 textAlign: 'center'
//               }}>{day}</td>
//               {times.map((time, index) => {
//                 const classes = schedule.filter(item =>
//                   item.day_name === day && `${item.start_time} - ${item.end_time}` === time
//                 );
//                 return (
//                   <td
//                     key={index}
//                     style={{
//                       border: '2px solid #ddd',
//                       padding: '7px',
//                       fontSize: '16px',
//                       fontWeight: 'bold',
//                       textAlign: 'center',
//                       cursor: isEditMode ? 'pointer' : 'default'
//                     }}
//                     onClick={() => handleOpenModal(day, time)}
//                   >
//                     {classes.length > 0 ? (
//                       classes.map((item, idx) => (
//                         <div key={idx}>
//                           <div>{item.subject_name}</div>
//                           <div>{item.faculty_name}</div>
//                         </div>
//                       ))
//                     ) : (
//                       <div>No classes</div>
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
//             Are you sure you want to remove the period on {selectedPeriod?.day} during {selectedPeriod?.time}?
//           </Typography>
//           <Box sx={{ mt: 4, display: 'flex', justifyContent: 'space-between' }}>
//             <Button variant="contained" color="primary" onClick={handleConfirmRemove}>
//               Yes, Remove
//             </Button>
//             <Button variant="outlined" color="secondary" onClick={handleCloseModal}>
//               No, Cancel
//             </Button>
//           </Box>
//         </Box>
//       </Modal>
//     </div>
//   );
// };

// export default SavedTimetable;

  // import React, { useState, useEffect } from 'react';
  // import axios from 'axios';
  // import CustomButton from '../../components/button';
  // import * as XLSX from 'xlsx';
  // import { Modal, Box, Typography, Button, Select, MenuItem, InputLabel, FormControl } from '@mui/material';

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
  //     XLSX.utils.book_append_sheet(workbook, worksheet, `Semester ${props.semesterID}`);
      
  //     XLSX.writeFile(workbook, `timetable_semester_${props.semesterID}.xlsx`);
  //   };

  //   const handleToggleEditMode = () => {
  //     setIsEditMode(!isEditMode);
  //   };

  //   return (
  //     <div style={{
  //       backgroundColor: '#fff',
  //       padding: '20px',
  //       borderRadius: '8px',
  //       boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)',
  //       margin: '20px 0'
  //     }}>
  //       <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', marginBottom: '13px' }}>
  //         <div style={{ display: 'flex', flexDirection: 'row' }}>
  //           <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : S{props.semesterID} </h2>
  //           <h2 style={{ fontSize: '20px', marginTop: '5px', marginLeft: '15px' }}>Venue: {venue || 'Not Available'}</h2>
  //         </div>
  //         <div style={{ display: 'flex', flexDirection: 'row', columnGap: 10 }}>
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
  //       <table style={{
  //         width: '100%',
  //         borderCollapse: 'collapse',
  //         backgroundColor: '#fff',
  //         border: '2px solid #ddd',
  //         fontSize: '16px',
  //         minHeight: '600px'
  //       }}>
  //         <thead>
  //           <tr style={{ backgroundColor: '#f4f4f4' }}>
  //             <th style={{
  //               border: '2px solid #ddd',
  //               padding: '12px',
  //               textAlign: 'center'
  //             }}>Day/Time</th>
  //             {times.map((time, index) => (
  //               <th key={index} style={{
  //                 border: '2px solid #ddd',
  //                 padding: '12px',
  //                 backgroundColor: '#f9f9f9',
  //                 textAlign: 'center'
  //               }}>
  //                 {time}
  //               </th>
  //             ))}
  //           </tr>
  //         </thead>
  //         <tbody>
  //           {days.map((day) => (
  //             <tr key={day}>
  //               <td style={{
  //                 border: '2px solid #ddd',
  //                 padding: '12px',
  //                 fontWeight: 'bold',
  //                 textAlign: 'center'
  //               }}>{day}</td>
  //               {times.map((time, index) => {
  //                 const classes = schedule.filter(item =>
  //                   item.day_name === day && `${item.start_time} - ${item.end_time}` === time
  //                 );
  //                 return (
  //                   <td
  //                     key={index}
  //                     style={{
  //                       border: '2px solid #ddd',
  //                       padding: '7px',
  //                       fontSize: '16px',
  //                       fontWeight: 'bold',
  //                       textAlign: 'center',
  //                       cursor: isEditMode ? 'pointer' : 'default'
  //                     }}
  //                     onClick={() => handleOpenModal(day, time)}
  //                   >
  //                     {classes.length > 0 ? (
  //                       classes.map((item, idx) => (
  //                         <div key={idx}>
  //                           <div>{item.subject_name}</div>
  //                           <div>{item.faculty_name}</div>
  //                         </div>
  //                       ))
  //                     ) : (
  //                       <CustomButton
  //                       width="6"
  //                       label={"show faculty"}
  //                       backgroundColor={"red"}
                    
  //                     />
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
  //           <Box sx={{ mt: 4, display: 'flex', justifyContent: 'space-between' }}>
  //             <Button variant="contained" color="primary" onClick={handleConfirmRemove}>
  //               Yes, Remove
  //             </Button>
  //             <Button variant="outlined" color="secondary" onClick={handleCloseModal}>
  //               No, Cancel
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
import CustomButton from '../../components/button';
import * as XLSX from 'xlsx';
import { Modal, Box, Typography, Button, Select, MenuItem, InputLabel, FormControl } from '@mui/material';
import './save.css'; 

const SavedTimetable = (props) => {
  const [schedule, setSchedule] = useState([]);
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);
  const [venue, setVenue] = useState('');
  const [selectedPeriod, setSelectedPeriod] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [availableFaculty, setAvailableFaculty] = useState([]);
  const [selectedFaculty, setSelectedFaculty] = useState('');

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

  const fetchAvailableFaculty = async () => {
    if (!selectedPeriod) return;

    try {
      const response = await axios.get(`http://localhost:8080/faculty/available/${props.departmentID}/${props.semesterID}/${selectedPeriod.day}/${selectedPeriod.start_time}/${selectedPeriod.end_time}`);
      setAvailableFaculty(response.data);
    } catch (error) {
      console.error('Error fetching available faculty:', error);
    }
  };

  const handleOpenModal = (day, time) => {
    if (!isEditMode) return;

    const [start_time, end_time] = time.split(' - ');
    setSelectedPeriod({ day, start_time, end_time });
    fetchAvailableFaculty();
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedPeriod(null);
    setAvailableFaculty([]);
    setSelectedFaculty('');
  };

  const handleConfirmRemove = () => {
    const updatedSchedule = schedule.filter(item =>
      !(item.day_name === selectedPeriod.day && `${item.start_time} - ${item.end_time}` === `${selectedPeriod.start_time} - ${selectedPeriod.end_time}`)
    );
    setSchedule(updatedSchedule);
    handleCloseModal();
  };

  const handleDownload = () => {
    const workbook = XLSX.utils.book_new();
    const worksheetData = [
      ['Day/Time', ...times]
    ];
  
    // Add rows to the worksheet data
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
  
    // Define styles
    const headerStyle = {
      fill: { fgColor: { rgb: "343a40" } }, // Background color for header
      font: { color: { rgb: "ffffff" }, bold: true, sz: 14 }, // Font color and size for header
      alignment: { horizontal: "center", vertical: "center" },
      border: { 
        top: { style: 'thin', color: { rgb: '000000' } },
        bottom: { style: 'thin', color: { rgb: '000000' } },
        left: { style: 'thin', color: { rgb: '000000' } },
        right: { style: 'thin', color: { rgb: '000000' } }
      }
    };
  
    const timeHeaderStyle = {
      fill: { fgColor: { rgb: "6c757d" } }, // Background color for time headers
      font: { color: { rgb: "ffffff" }, bold: true, sz: 14 }, // Font color and size for time headers
      alignment: { horizontal: "center", vertical: "center" },
      border: { 
        top: { style: 'thin', color: { rgb: '000000' } },
        bottom: { style: 'thin', color: { rgb: '000000' } },
        left: { style: 'thin', color: { rgb: '000000' } },
        right: { style: 'thin', color: { rgb: '000000' } }
      }
    };
  
    const dayCellStyle = {
      fill: { fgColor: { rgb: "6c757d" } }, // Background color for day cells
      font: { color: { rgb: "ffffff" }, bold: true, sz: 12 }, // Font color and size for day cells
      alignment: { horizontal: "center", vertical: "center" },
      border: { 
        top: { style: 'thin', color: { rgb: '000000' } },
        bottom: { style: 'thin', color: { rgb: '000000' } },
        left: { style: 'thin', color: { rgb: '000000' } },
        right: { style: 'thin', color: { rgb: '000000' } }
      }
    };
  
    const subjectCellStyle = {
      border: { 
        top: { style: 'thin', color: { rgb: '000000' } },
        bottom: { style: 'thin', color: { rgb: '000000' } },
        left: { style: 'thin', color: { rgb: '000000' } },
        right: { style: 'thin', color: { rgb: '000000' } }
      }
    };
  
    // Apply styles to header
    worksheet['A1'].s = headerStyle;
    times.forEach((_, index) => {
      worksheet[XLSX.utils.encode_cell({ r: 0, c: index + 1 })].s = timeHeaderStyle;
    });
  
    // Apply styles to each row
    worksheetData.forEach((row, rowIndex) => {
      row.forEach((_, colIndex) => {
        const cellAddress = XLSX.utils.encode_cell({ r: rowIndex, c: colIndex });
        if (rowIndex === 0) {
          // Apply header style to the first row
          worksheet[cellAddress].s = headerStyle;
        } else if (colIndex === 0) {
          // Apply dayCellStyle to the first column
          worksheet[cellAddress].s = dayCellStyle;
        } else {
          // Apply subjectCellStyle to all other cells
          worksheet[cellAddress].s = subjectCellStyle;
        }
      });
    });
  
    // Define column widths
    worksheet['!cols'] = [{ width: 20 }].concat(times.map(() => ({ width: 30 })));
  
    XLSX.utils.book_append_sheet(workbook, worksheet, `Semester ${props.semesterID}`);
    XLSX.writeFile(workbook, `timetable_semester_${props.semesterID}.xlsx`);
  };
  

  
  const handleToggleEditMode = () => {
    setIsEditMode(!isEditMode);
  };

  return (
    <div className="container">
      <div className="header-k">
        <div className="header-info">
          <h2>Semester : S{props.semesterID}</h2>
          <h2>Venue: {venue || 'Not Available'}</h2>
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
            onClick={handleToggleEditMode}
          />
        </div>
      </div>
      <table className="table">
        <thead>
          <tr>
            <th className="day-time">Day/Time</th>
            {times.map((time, index) => (
              <th key={index} className="time">
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
                const classes = schedule.filter(item =>
                  item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                );
                return (
                  <td
                    key={index}
                    className="subject"
                    onClick={() => isEditMode && handleOpenModal(day, time)}
                  >
                    {classes.length > 0 ? (
                      classes.map((item, idx) => (
                        <div key={idx}>
                          <div>{item.subject_name}</div>
                          <div>{item.faculty_name}</div>
                        </div>
                      ))
                    ) : (
                      isEditMode && (
                        <Button
                          variant="contained"
                          color="primary"
                          onClick={() => {
                            setSelectedPeriod({ day, start_time: time.split(' - ')[0], end_time: time.split(' - ')[1] });
                            fetchAvailableFaculty();
                            setIsModalOpen(true);
                          }}
                        >
                          Show Faculty
                        </Button>
                      )
                    )}
                    {availableFaculty.length > 0 && (
                      <FormControl fullWidth sx={{ mt: 2 }}>
                        <InputLabel id="select-faculty-label">Select Faculty</InputLabel>
                        <Select
                          labelId="select-faculty-label"
                          value={selectedFaculty}
                          onChange={(e) => setSelectedFaculty(e.target.value)}
                          displayEmpty
                        >
                          <MenuItem value="" disabled>Select Faculty</MenuItem>
                          {availableFaculty.map(faculty => (
                            <MenuItem key={faculty.id} value={faculty.name}>
                              {faculty.name}
                            </MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                    )}
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>

      <Modal
        open={isModalOpen}
        onClose={handleCloseModal}
      >
        <Box sx={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          width: 400,
          bgcolor: 'background.paper',
          border: '2px solid #000',
          boxShadow: 24,
          p: 4,
        }}>
          <Typography variant="h6" component="h2">
            Confirm Removal
          </Typography>
          <Typography sx={{ mt: 2 }}>
            Are you sure you want to remove the period on {selectedPeriod?.day} during {selectedPeriod?.start_time} - {selectedPeriod?.end_time}?
          </Typography>
          {availableFaculty.length > 0 && (
            <FormControl fullWidth sx={{ mt: 2 }}>
              <InputLabel id="select-faculty-label">Select Faculty</InputLabel>
              <Select
                labelId="select-faculty-label"
                value={selectedFaculty}
                onChange={(e) => setSelectedFaculty(e.target.value)}
                displayEmpty
              >
                <MenuItem value="" disabled>Select Faculty</MenuItem>
                {availableFaculty.map(faculty => (
                  <MenuItem key={faculty.id} value={faculty.name}>
                    {faculty.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}
          <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
            <Button variant="contained" color="primary" onClick={handleConfirmRemove}>
              Confirm
            </Button>
          </Box>
        </Box>
      </Modal>
    </div>
  );
};

export default SavedTimetable;
