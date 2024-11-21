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
    const [selectedOption, setSelectedOption] = useState(null);
    const [section, setSection] = useState([]);
    const [sectionOptions, setSectionOptions] = useState([]);

  

    // useEffect(() => {
    //     const fetchSubjectOptions = async () => {
    //         try {
          
    //             if (academicYear && semester && departments && section) {
    //                 for (const sem of semester) {
    //                     for (const dept of departments) {
                    
    //                         const subjectData = {
    //                             department_id: dept.value,
    //                             semester_id: sem.value,
    //                             academic_year_id: academicYear ? academicYear.value : null,
    //                             section_id: section.value,
    //                         };
    
    //                         console.log(subjectData);
    
                      
    //                         const response = await axios.post('http://localhost:8080/subjectoptions', subjectData);
    //                         setSubjectOptions(response.data);
    //                     }
    //                 }
    //             } else {
    //                 console.error('Missing one or more required values.');
    //             }
    //         } catch (error) {
    //             console.error('Error fetching subject options:', error);
    //         }
    //     };
    
     
    //     if (academicYear && semester && departments && section) {
    //         fetchSubjectOptions();
    //     }
    // }, [academicYear, semester, departments, section]);

    useEffect(() => {
        const fetchSubjectptions = async () => {
            try {
                const response = await axios.get('http://localhost:8080/subjectoptions');
                setSubjectOptions(response.data);
            } catch (error) {
                console.error('Error fetching subject  options:', error);
            }
        };
        fetchSubjectptions();
    }, []);
  
    

    

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
        try {
            const data = [];
    
            for (const sem of semester) {
                for (const dept of departments) {
                    for (const sec of section) { // Loop through selected sections
                        if (selectedOption.value === 0) { // Lab Subject
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
                                    status: selectedOption.value,
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
                                    status: selectedOption.value,
                                }
                            );
                        } else if (selectedOption.value === 1) { // Non-Lab Subject
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
                                status: selectedOption.value,
                            });
                        }
                    }
                }
            }
    
            console.log('Final data payload:', data);
    
            // Submit the data to the backend
            await axios.post('http://localhost:8080/manual/submit', data);
    
            setIsModalOpen(true); // Open the modal upon success
        } catch (error) {
            console.error('Error submitting form:', error);
        }
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
                        <CustomSelect
                                label="SUBJECT NAME"
                                placeholder="SUBJECT NAME"
                                value={subject}
                                onChange={setSubject}
                                options={subjectOptions}
                               
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

export default ManualEntry;