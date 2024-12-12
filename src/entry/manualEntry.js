import React, { useState, useEffect } from 'react';
import './entry.css';
import axios from 'axios';
import CustomButton from '../components/button';
import InputBox from '../components/input';
import AppLayout from '../layout/layout';
import CustomSelect from '../components/select';
import { Modal, Box, Typography } from '@mui/material'; // Import Modal components

function ManualEntry() {
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
    const [subjectOptions, setSubjectOptions] = useState([]);
    const [courseCode, setCourseCode] = useState('');
    const [courseCodeOptions, setCourseCodeOptions] = useState([]);
    const [faculty, setFaculty] = useState(null);
    const [facultyOptions, setFacultyOptions] = useState([]);
    const [academicYear, setAcademicYear] = useState(null);
    const [academicsOptions, setAcademicsOptions] = useState([]);
    const [venue, setVenue] = useState(null);
    const [venueOptions, setVenueOptions] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false); // State for modal visibility
    const [subjectType, setSubjectType] = useState('');
    const [subjectTypeOption, setsubjectTypeOption] = useState([]);
    const [section, setSection] = useState([]);
    const [sectionOptions, setSectionOptions] = useState([]);
    const [errorMessage, setErrorMessage] = useState('');

    useEffect(() => {
        const fetchSubjectTypeOptions = async () => {
            try {
                const response = await axios.get('http://localhost:8080/subjectTypeoptions');
                setsubjectTypeOption(response.data);
            } catch (error) {
                console.error('Error fetching subject  options:', error);
            }
        };
        fetchSubjectTypeOptions();
    }, []);



    useEffect(() => {
        const fetchSubjectOptions = async () => {
            try {
                const response = await axios.get('http://localhost:8080/subjectoptions', {
                    params: { subject_type_id: subjectType?.value || '' }, // Safeguard for undefined `subjectType`
                });
                setSubjectOptions(response.data || []); // Default to an empty array
            } catch (error) {
                console.error('Error fetching subject options:', error);
                setSubjectOptions([]); // Ensure fallback state on error
            }
        };
        fetchSubjectOptions();
    }, [subjectType]); // Re-fetch when `subjectType` changes
    
    
  
    

    

    useEffect(() => {
        const fetchCourseCodeOptions = async () => {
            if (!subject) return; 
            try {
                const response = await axios.get('http://localhost:8080/course-code', {
              
                    params: { subject_name: subject.label },
                });
                setCourseCodeOptions(Array.isArray(response.data) ? response.data : []);
             
            } catch (error) {
                console.error('Error fetching course code options:', error);
            }
        };
        fetchCourseCodeOptions();
    }, [subject]);


    useEffect(() => {
       
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

     
      useEffect(() => {
        if (academicYear) {
            const yearLabel = academicYear.label.toUpperCase();
            const isOdd = /ODD/.test(yearLabel);
            const filteredSemesters = semOptions.filter(sem => {
         
                return isOdd ? /S[1357]/i.test(sem.label) : /S[2468]/i.test(sem.label);
            });
            setFilteredSemOptions(filteredSemesters);
        } else {
          
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
        console.log(subjectType.value)
        try {
            const data = [];
    
            for (const sem of semester) {
                for (const dept of departments) {
                    for (const sec of section) { // Loop through selected sections
                        if (subjectType.value === 1) { // Lab Subject
                            data.push(
                                {
                                    subject_name: subject.label,
                                    department_id: dept.value,
                                    semester_id: sem.value,
                                    section_id: sec.value,
                                    day_name: day?.value,
                                    start_time: startTime?.value,
                                    end_time: endTime?.value,
                                    faculty_name: faculty?.value,
                                    classroom: venue?.value,
                                    academic_year: academicYear?.value,
                                    course_code: courseCode?.value,
                                    status: subjectType.value,
                                },
                                {
                                    subject_name: subject.label,
                                    department_id: dept.value,
                                    semester_id: sem.value,
                                    section_id: sec.value,
                                    day_name: day?.value,
                                    start_time: labstartTime?.value,
                                    end_time: labendTime?.value,
                                    faculty_name: faculty?.value,
                                    classroom: venue?.value,
                                    academic_year: academicYear?.value,
                                    course_code: courseCode?.value,
                                    status: subjectType.value,
                                }
                            );
                        } else if (subjectType.value === 2 || 3 || 4 || 5 || 6 || 7) { 
                            data.push({
                                subject_name: subject.label,
                                department_id: dept.value,
                                semester_id: sem.value,
                                section_id: sec.value,
                                day_name: day?.value,
                                start_time: startTime?.value,
                                end_time: endTime?.value,
                                faculty_name: faculty?.value,
                                classroom: venue?.value,
                                academic_year: academicYear?.value,
                                course_code: courseCode?.value,
                                status: subjectType.value,
                            });
                        }
                    }
                }
            }
    
            console.log('Final data payload:', data);
    
         
            await axios.post('http://localhost:8080/manual/submit', data);
            setErrorMessage(''); 
            setIsModalOpen(true);
        } catch (error) {
            console.error('Error submitting form:', error);
            setErrorMessage(
                error.response?.data?.message || 'An error occurred during submission.'
            );
            setIsModalOpen(true); 
        }
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setErrorMessage(''); 
    };  

    return (
        <AppLayout
            rId={1}
            title="Manual Entry"
            body={
                <>
                    <div className="manual-container">
                        <center>
                            <h1>Here you can upload the Manual entry</h1>
                        </center>
                        <br />
                        <div className="form-group">
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
              isMulti={true}
            />
                        </div>
                        <div className="form-group">
                        <CustomSelect
        label="SUBJECT TYPE"
        options={subjectTypeOption}
          value={subjectType}
                                onChange={setSubjectType}
        placeholder="SUBJECT TYPE"
        
    
    
      />
                        </div>
                        <CustomSelect
                                label="SUBJECT NAME"
                                placeholder="SUBJECT NAME"
                                value={subject}
                                onChange={setSubject}
                                
                                options={subjectOptions.length > 0 ? subjectOptions : []}
                            />
                        </div>
                       
                        <div className="form-group">
                        <CustomSelect
                                label="COURSE CODE"
                                placeholder="COURSE CODE"
                                value={courseCode}
                                 onChange={setCourseCode}
                                options={courseCodeOptions}
                               
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
                        {subjectType && subjectType.value === 1 && (
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
                                    {errorMessage ? 'Submission Failed!' : 'Submission Successful!'}
                                </Typography>
                                {errorMessage && (
                                    <Typography id="modal-description" variant="body1" color="error" className="modal-description">
                                        {errorMessage}
                                    </Typography>
                                )}
                                {!errorMessage && (
                                    <Typography id="modal-description" variant="body1" className="modal-description">
                                        Your data has been successfully submitted.
                                    </Typography>
                                )}
                                <CustomButton
                                    width="150px"
                                    label="Close"
                                    backgroundColor="#0878d3"
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

export default ManualEntry;