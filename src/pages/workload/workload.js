import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CustomButton from '../../components/button';

import * as XLSX from 'xlsx'; 

const Timetable = (props) => {

  const [schedule, setSchedule] = useState({});
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);
  const [venue, setVenue] = useState('');

  useEffect(() => {
    const fetchSchedule = async () => {
      if (!props.departmentID || !props.semesterID) {
        console.error('Department ID and Semester ID are required');
        return;
      }
    
      try {
        const response = await axios.get(`http://localhost:8080/timetable/${props.departmentID}/${props.semesterID}`);
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
                classrooms.add(subject.classroom); 
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
      } catch (error) {
        console.error('Error fetching timetable:', error);
      }
    };
    
    fetchSchedule();
  }, [props.departmentID, props.semesterID]);

  const handleSaveTimetable = async (timetableData) => {
    try {
      await axios.post('http://localhost:8080/timetable/save', timetableData);
      alert('Timetable saved successfully!');
    } catch (err) {
      console.error("Error saving timetable:", err);
      alert('Failed to save timetable');
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
          const data = {
            day_name: entry.day_name,
            start_time: entry.start_time,
            end_time: entry.end_time,
            subject_name: entry.subject_name,
            faculty_name: entry.faculty_name,
            classroom: entry.classroom,
            status: entry.status,
            semester_id: entry.semester_id,
            department_id: entry.department_id,
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

  return (
    <div style={{ 
      backgroundColor: '#fff', 
      padding: '20px', 
      borderRadius: '8px', 
     boxShadow: ' 0px 0px 4px 2px rgba(0, 0, 0, 0.1)', 
      margin: '20px 0',

    }}>
      <div style={{display:'flex',flexDirection:'row',justifyContent:'space-between',marginBottom:'13px'}}>
        <div style={{display:'flex',flexDirection:'row'}}>
          <h2 style={{fontSize:'20px',marginTop:'5px'}}>Semester : S{props.semesterID}</h2>
          <h2 style={{fontSize:'20px',marginTop:'5px',marginLeft:'30px'}}>Venue : {venue}</h2>
        </div>
        <div style={{display:'flex',flexDirection:'row',columnGap:10}}>
          <CustomButton
            width="150"
            label="Download Timetable"
            onClick={handleDownload}
          />
          <CustomButton
            width="150"
            label="Save Timetable"
            onClick={handleSave} 
            backgroundColor="red"
          />
        </div>
      </div>
      <table style={{ 
  width: '100%', 
  borderCollapse: 'collapse', 
  backgroundColor: '#ffffff', 
  border: '2px solid #dedede', 
  fontSize: '16px',
  minHeight: '600px',
  marginBottom: '20px',
  boxShadow: '0px 4px 15px rgba(0, 0, 0, 0.1)', 
  borderRadius: '8px',
  overflow: 'hidden'
}}>
  <thead>
    <tr style={{ backgroundColor: '#007bff' }}> {/* Bold blue for the header */}
      <th style={{ 
        border: '2px solid #dedede', 
        padding: '14px', 
        textAlign: 'center',
        backgroundColor: '#343a40', // Dark grey background for main header
        color: '#ffffff', // White text for contrast
        fontWeight: '600', // Slightly bolder text
      }}>Day/Time</th>
      {times.map((time, index) => (
        <th key={index} style={{ 
          border: '2px solid #dedede', 
          padding: '14px', 
          backgroundColor: '#6c757d', // Consistent blue background for time columns
          textAlign: 'center',
          color: '#ffffff', // White text for clarity
          fontWeight: '600', 
        }}>
          {time}
        </th>
      ))}
    </tr>
  </thead>
  <tbody>
    {days.map((day, dayIndex) => (
      <tr key={day} 
      style={{ 
        backgroundColor: '#f8f9fa', // Light background color for rows
        transition: 'transform 0.3s ease, background-color 0.3s ease',
        cursor: 'pointer'
      }} 
     
      >
        <td style={{ 
          border: '2px solid #dedede', 
          padding: '14px', 
          fontWeight: '700', // Bolder text for day labels
          textAlign: 'center',
          backgroundColor: '#6c757d', // Medium grey background for day column
          color: '#ffffff', // White text for contrast
        }}>{day}</td>
        {times.map((time, index) => (
          <td key={index} style={{ 
            border: '2px solid #dedede', 
            padding: '14px', 
            textAlign: 'center',
            color: '#212529', // Darker text for readability
           // backgroundColor: '#e9ecef', // Light grey background for cells
          }}>
            {Object.values(schedule).flatMap(faculty =>
              Object.values(faculty).flatMap(subjects =>
                Array.isArray(subjects) ? subjects.filter(
                  item => item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                ) : []
              )
            ).map((item, idx) => (
              <div key={idx} style={{
                marginBottom: '10px',
                padding: '10px',
                //backgroundColor: '#17a2b8', // Attractive teal background for subject blocks
                color: 'black', // White text for clear visibility
                borderRadius: '5px', // Smooth corners for a modern feel
                fontWeight: '800', // Medium weight for better readability
              }}>
                <div>{item.subject_name}</div>
                <div>{item.faculty_name}</div>
              </div>
            ))}
          </td>
        ))}
      </tr>
    ))}
  </tbody>
</table>

    </div>
  );
};

export default Timetable;
