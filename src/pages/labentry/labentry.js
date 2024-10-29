import React, { useState, useEffect } from 'react';
import './labentry.css';
import axios from 'axios';

import { Modal, Box, Typography } from '@mui/material'; // Import Modal components
import AppLayout from '../../layout/layout';
import InputBox from '../../components/input';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';

function LabEntry() {
    const [departments, setDepartments] = useState([]);
    const [deptOptions, setDeptOptions] = useState([]);
    const [semester, setSemester] = useState([]);
    const [semOptions, setSemOptions] = useState([]);
    const [filteredSemOptions, setFilteredSemOptions] = useState([]); // State for filtered semesters
    const [day, setDay] = useState(null);
    const [dayOptions, setDayOptions] = useState([]);
    const [startTime, setStartTime] = useState(null);
    const [labstartTime, setLabStartTime] = useState(null);
    const [startTimeOptions, setStartTimeOptions] = useState([]);
    const [endTime, setEndTime] = useState(null);
    const [labendTime, setLabEndTime] = useState(null);
    const [endTimeOptions, setEndTimeOptions] = useState([]);
    const [subject, setSubject] = useState('');
    const [courseCode, setCourseCode] = useState('');
    const [faculty, setFaculty] = useState(null);
    const [facultyOptions, setFacultyOptions] = useState([]);
    const [academicYear, setAcademicYear] = useState(null);
    const [academicsOptions, setAcademicsOptions] = useState([]);
    const [venue, setVenue] = useState(null);
    const [venueOptions, setVenueOptions] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false); // State for modal visibility
    const [selectedOption, setSelectedOption] = useState(null);
    const [section, setSection] = useState(null);
    const [sectionOptions, setSectionOptions] = useState([]);

    useEffect(() => {
        // Fetching initial options from the backend
        const fetchOptions = async () => {
            try {
                const response = await axios.get('http://localhost:8080/manual/options');
                setDayOptions(response.data.dayOptions);
                setStartTimeOptions(response.data.startTimeOptions);
                setEndTimeOptions(response.data.endTimeOptions);
                setFacultyOptions(response.data.facultyOptions);
            } catch (error) {
                console.error('Error fetching options:', error);
            }
        };
        fetchOptions();
    }, []);

    useEffect(() => {
        const fetchDeptOptions = async () => {
            try {
                const response = await axios.get('http://localhost:8080/timetable/options');
                setDeptOptions(response.data);
            } catch (error) {
                console.error('Error fetching department options:', error);
            }
        };
        fetchDeptOptions();
    }, []);

    useEffect(() => {
        const fetchSemOptions = async () => {
            try {
                const response = await axios.get('http://localhost:8080/timetable/semoptions');
                setSemOptions(response.data);
            } catch (error) {
                console.error('Error fetching semester options:', error);
            }
        };
        fetchSemOptions();
    }, []);

    useEffect(() => {
        const fetchAcademicYears = async () => {
            try {
                const response = await axios.get('http://localhost:8080/acdemicYearOptions');
                setAcademicsOptions(response.data);
            } catch (error) {
                console.error('Error fetching academic year options:', error);
            }
        };
        fetchAcademicYears();
    }, []);

    useEffect(() => {
        const fetchSection = async () => {
          try {
            const response = await axios.get('http://localhost:8080/timetable/sectionoptions');
            setSectionOptions(response.data);
          } catch (error) {
            console.error('Error fetching section options:', error);
        
          }
        };
    
        fetchSection();
      }, []);

      // Effect to update semesters based on the selected academic year
      useEffect(() => {
        if (academicYear) {
            const yearLabel = academicYear.label.toUpperCase();
            const isOdd = /ODD/.test(yearLabel); // Check if it contains 'ODD'
            const filteredSemesters = semOptions.filter(sem => {
                // Show only odd or even semesters based on the selected academic year
                return isOdd ? /S[1357]/i.test(sem.label) : /S[2468]/i.test(sem.label);
            });
            setFilteredSemOptions(filteredSemesters);
        } else {
            // Reset semesters if no academic year is selected
            setFilteredSemOptions(semOptions);
        }
    }, [academicYear, semOptions]);

    useEffect(() => {
        const fetchClassroomOptions = async () => {
            try {
                const response = await axios.get('http://localhost:8080/classroomDetailsOptions');
                setVenueOptions(response.data);
            } catch (error) {
                console.error('Error fetching classroom options:', error);
            }
        };
        fetchClassroomOptions();
    }, []);
    

    const handleSubmit = async () => {
        try {
            for (const sem of semester) {
                for (const dept of departments) {
            let data;
    
            if (selectedOption.value === 0) { // First option
                data = [{
                    subject_name: subject,
                    department_id: dept.value,
                    semester_id: sem.value,
                    day_name: day ? day.value : null,
                    start_time: startTime ? startTime.value : null,
                    end_time: endTime ? endTime.value : null,
                    faculty_name: faculty ? faculty.value : null,
                    classroom: venue ? venue.value : null,
                    academic_year: academicYear ? academicYear.value : null,
                    course_code: courseCode,
                    status: selectedOption.value,
                    section_id: section.value,
                },
                {
                    subject_name: subject,
                    department_id: dept.value,
                        semester_id: sem.value,
                    day_name: day ? day.value : null,
                    faculty_name: faculty ? faculty.value : null,
                    classroom: venue ? venue.value : null,
                    academic_year: academicYear ? academicYear.value : null,
                    course_code: courseCode,
                    status: selectedOption.value,
                    section_id: section.value,
                    start_time: labstartTime ? labstartTime.value : null, 
                    end_time: labendTime ? labendTime.value : null, 
                }];
            } else if (selectedOption.value === 1) { // Second option
                data = [{
                    subject_name: subject,
                    department_id: dept.value,
                        semester_id: sem.value,
                    day_name: day ? day.value : null,
                    start_time: startTime ? startTime.value : null,
                    end_time: endTime ? endTime.value : null,
                    faculty_name: faculty ? faculty.value : null,
                    classroom: venue ? venue.value : null,
                    academic_year: academicYear ? academicYear.value : null,
                    course_code: courseCode,
                    status: selectedOption.value,
                    section_id: section.value,
                }];
            }
    
            console.log('Data to be sent for ', data);
    
            // Sending the data to the backend
            await axios.post('http://localhost:8080/labentry', data);
            // console.log('Form submitted successfully for department', dept.value, 'and semester', sem.value);
            setIsModalOpen(true); 
        }
    }// Open the modal upon successful submission
        } catch (error) {
            console.error('Error submitting form:', error);
        }
    };
    
    const handleCloseModal = () => {
        setIsModalOpen(false);
    };


    return (
        <AppLayout
            rId={16}
            title="Lab Entry"
            body={
                <>
                    <div className="manual-container">
                        <center>
                            <h1>Here you can upload the Lab entry</h1>
                        </center>
                        <br />
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
        label="Choose an Option"
        options={[
            { label: "Lab subject", value: 0 },
            { label: "Non-Lab Subject", value: 1 },
        
          ]}
        placeholder="Select an option"
    
        onChange={setSelectedOption}
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
                        <div className="form-group">
                            <CustomSelect
                                label="SEMESTER"
                                placeholder="SEMESTER"
                                value={semester}
                                onChange={setSemester}
                                options={filteredSemOptions} // Use filtered options
                                isMulti={true}
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
                        label="SECTION"
              placeholder="SECTION"
              value={section}
              onChange={setSection}
              options={sectionOptions}
            />
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
                        {selectedOption && selectedOption.value === 0 && (
        <div className='row'>
          <CustomSelect
            label="START TIME"
            placeholder="START TIME"
            value={labstartTime}
            onChange={setLabStartTime}
            options={startTimeOptions}
          />
          <CustomSelect
            label="END TIME"
            placeholder="END TIME"
            value={labendTime}
            onChange={setLabEndTime}
            options={endTimeOptions}
          />
        </div>
      )}
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
                                <Typography id="modal-description" className="modal-description">
                                    Your manual entry has been submitted successfully.
                                </Typography>
                                <CustomButton
                                    label="Close"
                                    onClick={handleCloseModal}
                                />
                            </Box>
                        </Modal>
                    </div>
                </>
            }
        />
    );
}

export default LabEntry;
