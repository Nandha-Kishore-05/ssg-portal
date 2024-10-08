import React, { useState, useEffect } from 'react';
import './entry.css';

import axios from 'axios';
import CustomButton from '../components/button';
import InputBox from '../components/input';
import AppLayout from '../layout/layout';
import CustomSelect from '../components/select';
import { Modal, Box, Typography } from '@mui/material'; // Import Modal components

function ManualEntry() {
    const [departments, setDepartments] = useState([]); // Update state to handle multiple departments
    const [deptOptions, setDeptOptions] = useState([]);
    const [semester, setSemester] = useState([]);
    const [semOptions, setSemOptions] = useState([]);
    const [day, setDay] = useState(null);
    const [dayOptions, setDayOptions] = useState([]);
    const [startTime, setStartTime] = useState(null);
    const [startTimeOptions, setStartTimeOptions] = useState([]);
    const [endTime, setEndTime] = useState(null);
    const [endTimeOptions, setEndTimeOptions] = useState([]);
    const [subject, setSubject] = useState('');
    const [courseCode, setCourseCode] = useState('');
    const [faculty, setFaculty] = useState(null);
    const [facultyOptions, setFacultyOptions] = useState([]);
    const [academicYear, setAcademicYear] = useState(null);
    const [academicsOptions, setAcademicsOptions] = useState(null);
    const [venue, setVenue] = useState(null);
    const [venueOptions, setVenueOptions] = useState(null);
    const [isModalOpen, setIsModalOpen] = useState(false); // State for modal visibility
    

    useEffect(() => {
        axios.get('http://localhost:8080/manual/options')
            .then(response => {
                setDayOptions(response.data.dayOptions);
                setStartTimeOptions(response.data.startTimeOptions);
                setEndTimeOptions(response.data.endTimeOptions);
                setFacultyOptions(response.data.facultyOptions);
            })
            .catch(error => {
                console.error('Error fetching options:', error);
            });
    }, []);

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
            console.error('Error fetching semester options:', error);
          });
      }, []);

      useEffect(() => {
        axios.get('http://localhost:8080/classroomOptions')
          .then(response => {
            setVenueOptions(response.data);
          })
          .catch(error => {
            console.error('Error fetching semester options:', error);
          });
      }, []);

    const handleSubmit = () => {
        if (semester.length === 0) {
            console.error("No semesters selected");
            return;
        }

        // Loop through each semester and each department and submit the form data separately
        semester.forEach((sem) => {
            departments.forEach((dept) => {
               
                const data = {
                    subject_name: subject,
                    department_id: dept.value,
                    semester_id: sem.value,
                    // status: status,
                    day_name: day ? day.value : null,
                    start_time: startTime ? startTime.value : null,
                    end_time: endTime ? endTime.value : null,
                    faculty_name: faculty ? faculty.value : null,
                    classroom : venue ? venue.value : null,
                    academic_year: academicYear? academicYear.value : null,
                    course_code : courseCode,
                };

                console.log('Data to be sent for department:', dept.value, 'and semester:', sem.value, data);

                axios.post('http://localhost:8080/manual/submit', data)
                    .then(response => {
                        console.log('Form submitted successfully for department', dept.value, 'and semester', sem.value, response.data);
                        setIsModalOpen(true); // Open the modal upon successful submission
                    })
                    .catch(error => {
                        console.error('Error submitting form for department', dept.value, 'and semester', sem.value, error);
                    });
            });
        });
    };

    

    const handleCloseModal = () => {
        setIsModalOpen(false); 
    };

    return (
        <AppLayout
            rId={7}
            title="Manual Entry"
            body={
                <>
                    <div className="manual-container">
                        <center>
                            <h1>Here you can upload the Manual entry</h1>
                        </center>
                        <br />
                   
                        <div className="form-group">
                            {/* <div className='bulk-button'>
                                <CustomButton
                                    width="150px"
                                    label="Bulk Edit"
                                    backgroundColor="#0878d3"
                                    onClick={handleBulkEditClick} // Trigger navigation on click
                                />
                            </div> */}
                            <div className="form-group">
                                <InputBox
                                    label="SUBJECT NAME"
                                    placeholder="SUBJECT NAME"
                                    value={subject}
                                    onChange={setSubject}
                                />
                            </div>
                            <div className="form-group">
                                <InputBox
                                    label="COURSE CODE"
                                    placeholder="COURSE CODE"
                                    value={courseCode}
                                    onChange={setCourseCode}
                                />
                            </div>
                            <div className="form-group">
                                <CustomSelect
                                    label="DEPARTMENT"
                                    placeholder="DEPARTMENT"
                                    value={departments}
                                    onChange={setDepartments}
                                    options={deptOptions}
                                    isMulti={true} // Enable multi-select
                                />
                                      </div>
                                      <div className="form-group">
                                <CustomSelect
                                    label="SEMESTER"
                                    placeholder="SEMESTER"
                                    value={semester}
                                    onChange={setSemester}
                                    options={semOptions}
                                    isMulti={true} 
                                />
                            </div>
                        </div>
                        <div className="form-group">
                            <CustomSelect
                                label="CLASSROOM"
                                placeholder="CLASSROOM"
                                value={venue}
                                onChange={setVenue}
                                options={venueOptions}
                            />
                        </div>
                        <div className="form-group">
                            <CustomSelect
                                label="FACULTY"
                                placeholder="FACULTY"
                                value={faculty}
                                onChange={setFaculty}
                                options={facultyOptions}
                            />
                        </div>
                        
                        <div className='row'>
                            <CustomSelect
                                label="START TIME"
                                placeholder="START TIME"
                                value={startTime}
                                onChange={setStartTime}
                                options={startTimeOptions}
                            />
                            <CustomSelect
                                label="END TIME"
                                placeholder="END TIME"
                                value={endTime}
                                onChange={setEndTime}
                                options={endTimeOptions}
                            />
                        </div>
                        <div className="form-group">
                            <CustomSelect
                                label="DAY"
                                placeholder="DAY"
                                value={day}
                                onChange={setDay}
                                options={dayOptions}
                            />
                        </div>
                        <div className="form-group">
                            <CustomSelect
                                label="ACADEMIC YEAR"
                                placeholder="ACADEMIC YEAR"
                                value={academicYear}
                                onChange={setAcademicYear}
                                options={academicsOptions}
                            />
                        </div>
                       

                        <div className="center-button">
                            <CustomButton
                                width="150px"
                                label="Submit"
                                backgroundColor="#0878d3"
                                onClick={handleSubmit}
                            />
                        </div>

                   
                        <Modal
                            open={isModalOpen}
                            onClose={handleCloseModal}
                            aria-labelledby="modal-title"
                            aria-describedby="modal-description"
                        >
                            <Box className="modal-box">
                                <Typography id="modal-title" variant="h5" component="h1" className="modal-title">
                                    Submission Successful!
                                </Typography>
                                <Typography id="modal-description" sx={{ mt: 2 }} className="modal-description">
                                    Your manual entry has been submitted successfully.
                                </Typography>
                                <div className="center-button">
                                    <CustomButton
                                        width="150px"
                                        label="Close"
                                        backgroundColor="#0878d3"
                                        onClick={handleCloseModal}
                                    />
                                </div>
                            </Box>
                        </Modal>
                    </div>
                </>
            }
        />
    );
}

export default ManualEntry;
