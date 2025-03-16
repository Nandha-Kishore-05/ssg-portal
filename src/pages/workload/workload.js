

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CustomButton from '../../components/button';
import * as XLSX from 'xlsx';
import Modal from '@mui/material/Modal';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Skeleton from '@mui/material/Skeleton';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useNavigate } from 'react-router-dom';
import './workload.css';
import { colors } from '@mui/material';

const style = {
  position: 'absolute',
  top: '45%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: 'background.paper',
  border: 'none',
  borderRadius: '10px',
  boxShadow: '0px 4px 20px rgba(0, 0, 0, 0.1)',
  p: 4,
};

const enhancedStyle = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: '420px',
  bgcolor: '#F9FAFB', // Light neutral background
  border: '1px solid #E0E0E0', // Subtle border for a clean look
  boxShadow: '0px 8px 24px rgba(0, 0, 0, 0.12)', // Soft shadow for depth
  borderRadius: '12px', // Smooth rounded corners
  p: 4,
};



const Timetable = (props) => {
  const navigate = useNavigate();
  const [schedule, setSchedule] = useState({});
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);
  const [venue, setVenue] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [showErrorModal,  setShowErrorModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  const [loading, setLoading] = useState(true);  // State for loading

  useEffect(() => {
    const fetchSchedule = async () => {
      setLoading(true);  // Start loading
      if (!props.departmentID || !props.semesterID) {
        console.error('Department ID and Semester ID are required');
        return;
      }

      try {
        const response = await axios.get(`http://localhost:8080/timetable/${props.departmentID}/${props.semesterID}/${props.academicYearID}/${props.sectionID}`);
        const data = response.data;

        console.log('Fetched data:', data);
        const allDays = new Set();
        const allTimes = new Set();
        const classrooms = new Set();

        Object.values(data).forEach(facultyDays => {
          Object.entries(facultyDays).forEach(([day, subjects]) => {
            if (Array.isArray(subjects)) {
              subjects.forEach(subject => {
                allDays.add(subject.day_name);
                allTimes.add(`${subject.start_time} - ${subject.end_time}`);
                classrooms.add(subject.classroom );
              });
            } else {
              console.warn('Subjects is not an array for day:', day, subjects);
            }
          });
        });

        const sortedDays = Array.from(allDays).sort((a, b) => {
          const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
          return order.indexOf(a) - order.indexOf(b);
        });

        const sortedTimes = Array.from(allTimes).sort((a, b) => {
          return a.localeCompare(b, undefined, { numeric: true });
        });

        const firstClassroom = Array.from(classrooms)[0] || 'Not Available';

        setVenue(firstClassroom);
        setDays(sortedDays);
        setTimes(sortedTimes);
        setSchedule(data);
      }catch (error) {
        console.error('Error fetching timetable:', error);
        if (error.response && error.response.data && error.response.data.error) {
            setErrorMessage(error.response.data.error); // Set the error message
            setShowErrorModal(true); // Show the modal
        }
    } finally {
        setLoading(false);  // Stop loading
    }
    };

    fetchSchedule();
  }, [props.departmentID, props.semesterID,props.academicYearID,props.sectionID ]);

  const handleSaveTimetable = async (timetableData) => {
    try {
      await axios.post('http://localhost:8080/timetable/save', timetableData);
    } catch (err) {
      console.error("Error saving timetable:", err);
    }
  };

  const handleSave = async () => {
    const timetableData = [];

    days.forEach(day => {
      times.forEach(time => {
        const entries = Object.values(schedule).flatMap(faculty =>
          Object.values(faculty).flatMap(subjects =>
            Array.isArray(subjects) ? subjects.filter(
              item => item.day_name === day && `${item.start_time} - ${item.end_time}` === time
            ) : []
          )
        );

        entries.forEach(entry => {
          console.log(entry.lab_name)
          const data = {
            day_name: entry.day_name,
            start_time: entry.start_time,
            end_time: entry.end_time,
            subject_name: entry.subject_name,
            faculty_name: entry.faculty_name,
            classroom: entry.classroom || entry.lab_name,
            status: entry.status,
            semester_id: entry.semester_id,
            department_id: entry.department_id,
            academic_year_id : entry.academic_year_id,
            course_code: entry.course_code,
            section_id : entry.section_id,
           
 
          };
          timetableData.push(data);
        });
      });
    });

    console.log("Final Timetable Data to Save:", timetableData);

    await handleSaveTimetable(timetableData);
  };

  const handleDownload = () => {
    const wsData = [
      ["Day/Time", ...times]
    ];

    days.forEach(day => {
      const row = [day];
      times.forEach(time => {
        const cellData = Object.values(schedule).flatMap(faculty =>
          Object.values(faculty).flatMap(subjects =>
            Array.isArray(subjects) ? subjects.filter(
              item => item.day_name === day && `${item.start_time} - ${item.end_time}` === time
            ) : []
          )
        ).map((item) => `${item.subject_name} (${item.faculty_name})`).join(', ');
        row.push(cellData || '');
      });
      wsData.push(row);
    });

    const worksheet = XLSX.utils.aoa_to_sheet(wsData);
    const workbook = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(workbook, worksheet, "Timetable");

    XLSX.writeFile(workbook, `Timetable_S${props.semesterID}.xlsx`);
  };

  const handleOpenModal = () => {
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
  };
  const handleCloseErrorModal = () => {
    setShowErrorModal(false);
    setErrorMessage(''); // Clear the error message
};


  const handleConfirmSave = async () => {
    await handleSave();
    handleCloseModal();
    
    toast.success('Timetable saved successfully!', {
      position: "top-center",
      autoClose: 5000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      theme: "light",
    });
    navigate('/timetable/saved');
  };

  return (
    <div className="container">
      {loading ? (
        <div>
          <Skeleton variant="rectangular" width="100%" height={40} style={{ marginBottom: '10px' }} />
          <Skeleton variant="rectangular" width="100%" height={600} style={{ borderRadius: '8px' }} />
        </div>
      ) : (
        <>
          <div className="header-1">
            <div className="header-info">
              <h2 className="semester">Semester : S{props.semesterID}</h2>
              <h2 className="venue">Venue : {venue}</h2>
              <h2 className="venue">Academic Year : {props.academicYearID}</h2>
            </div>
            <div className="buttons">
              <CustomButton
                width="150"
                label="Download Timetable"
                onClick={handleDownload}
              />
              <CustomButton
                width="150"
                label="Save Timetable"
                onClick={handleOpenModal}
                backgroundColor="red"
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
              {days.map((day, dayIndex) => (
                <tr key={day}>
                  <td className="day">{day}</td>
                  {times.map((time, index) => (
                    <td key={index}>
                      {Object.values(schedule).flatMap(faculty =>
                        Object.values(faculty).flatMap(subjects =>
                          Array.isArray(subjects) ? subjects.filter(
                            item => item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                          ) : []
                        )
                      ).map((item, idx) => (
                        <div key={idx} className="subject">
                          <div>{item.subject_name}</div>
                          <div>{item.faculty_name}</div>
                          <div>{item.classroom}  {item.lab_name}</div>
                        </div>
                      ))}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
         
          <ToastContainer
            position="top-center"
            autoClose={5000}
            hideProgressBar={false}
            newestOnTop={false}
            closeOnClick
            rtl={false}
            pauseOnFocusLoss
            draggable
            pauseOnHover
            theme="light"
          />

          <Modal
            open={showModal}
            onClose={handleCloseModal}
            aria-labelledby="save-timetable-modal-title"
            aria-describedby="save-timetable-modal-description"
          >
            <Box sx={style}>
              <center>
                <Typography id="save-timetable-modal-description" sx={{ mt: 2 }}>
                  <h3> Are you sure you want to save the timetable?</h3>
                </Typography>
              </center>
              <div className='saveButton' >
                <CustomButton
         
                  label="Save Timetable"
                  onClick={handleConfirmSave}
                />
                <CustomButton
                
                  label="Cancel"
                  onClick={handleCloseModal}
                  backgroundColor="red"
                />
              </div>
            </Box>
          </Modal>
          <Modal
    open={showErrorModal}
    onClose={handleCloseErrorModal}
    aria-labelledby="modal-title"
    aria-describedby="modal-description"
>
    <Box sx={enhancedStyle}>
        <Typography
            id="modal-title"
            variant="h5"
            component="h2"
            sx={{
                color: '#333333', // Dark gray for the title
                fontWeight: 'bold',
                textAlign: 'center',
                fontSize: '24px', // Larger, attractive font
                marginBottom: '10px',
            }}
        >
            Something Went Wrong
        </Typography>
        <Typography
            id="modal-description"
            sx={{
                color: '#555555', // Medium gray for description
                textAlign: 'center',
                fontSize: '16px', // Readable and elegant
                lineHeight: '1.8',
                marginBottom: '20px',
            }}
        >
            {errorMessage}
        </Typography>
        <Box
            sx={{
                display: 'flex',
                justifyContent: 'center',
            }}
        >
            <CustomButton
                label="Close"
                onClick={handleCloseErrorModal}
                width="150"
                backgroundColor="#0056D2" // Professional blue for button
                textColor="#FFFFFF"
                fontSize="16px"
                borderRadius="8px"
            />
        </Box>
    </Box>
</Modal>


        </>
      )}
    </div>
  );
};

export default Timetable;

