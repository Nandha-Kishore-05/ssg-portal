import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CustomButton from '../../components/button';
import { useParams } from 'react-router-dom';
import AppLayout from '../../layout/layout';

const Timetable = () => {
  const { departmentID } = useParams();
  const [schedule, setSchedule] = useState({});
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);

  useEffect(() => {
    const fetchSchedule = async () => {
      if (!departmentID) {
        console.error('Department ID is required');
        return;
      }

      try {
        const response = await axios.get(`http://localhost:8080/timetable/${departmentID}`);
        const data = response.data;

        console.log('Fetched data:', data);

        const allDays = new Set();
        const allTimes = new Set();

        Object.values(data).forEach(faculty =>
          Object.values(faculty).forEach(subjects =>
            subjects.forEach(subject => {
              allDays.add(subject.day_name);
              allTimes.add(`${subject.start_time} - ${subject.end_time}`);
            })
          )
        );

        const sortedDays = Array.from(allDays).sort((a, b) => {
          const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
          return order.indexOf(a) - order.indexOf(b);
        });

        const sortedTimes = Array.from(allTimes).sort((a, b) => {
          return a.localeCompare(b, undefined, { numeric: true });
        });

        setDays(sortedDays);
        setTimes(sortedTimes);
        setSchedule(data);
      } catch (error) {
        console.error('Error fetching timetable:', error);
      }
    };

    fetchSchedule();
  }, [departmentID]);

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

    // Iterate through each day
    days.forEach(day => {
      // Iterate through each time slot
      times.forEach(time => {
        // Find all schedule entries for the current day and time
        const entries = Object.values(schedule).flatMap(faculty =>
          Object.values(faculty).flatMap(subjects =>
            Array.isArray(subjects) ? subjects.filter(
              item => item.day_name === day && `${item.start_time} - ${item.end_time}` === time
            ) : []
          )
        );

        // Add all found entries to the timetableData array
        entries.forEach(entry => {
          const data = {
            day_name: entry.day_name,
            start_time: entry.start_time,
            end_time: entry.end_time,
            subject_name: entry.subject_name,
            faculty_name: entry.faculty_name,
            classroom: entry.classroom
          };
          timetableData.push(data);
        });
      });
    });

    console.log("Final Timetable Data to Save:", timetableData);

    // Call the function to save the timetable data
    await handleSaveTimetable(timetableData);
  };

  return (
    <AppLayout
      rId={2}
      title="Dashboard"
      body={
        <div style={{ 
          backgroundColor: '#fff', 
          padding: '20px', 
          borderRadius: '8px', 
          boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)', 
          margin: '20px 0'
        }}>
          <div style={{display:'flex',flexDirection:'row',justifyContent:'space-between',marginBottom:'13px'}}>
            <h2 style={{fontSize:'20px',marginTop:'5px'}}>Venue : WW212</h2>
            <CustomButton
              width="150"
              label="Save Timetable"
              onClick={handleSave}
            />
          </div>
          <table style={{ 
            width: '100%', 
            borderCollapse: 'collapse', 
            backgroundColor: '#fff', 
            border: '2px solid #ddd',
            fontSize: '16px',
            minHeight: '600px'
          }}>
            <thead>
              <tr style={{ backgroundColor: '#f4f4f4' }}>
                <th style={{ 
                  border: '2px solid #ddd', 
                  padding: '12px', 
                  textAlign: 'center'
                }}>Day/Time</th>
                {times.map((time, index) => (
                  <th key={index} style={{ 
                    border: '2px solid #ddd', 
                    padding: '12px', 
                    backgroundColor: '#f9f9f9', 
                    textAlign: 'center'
                  }}>
                    {time}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {days.map((day) => (
                <tr key={day}>
                  <td style={{ 
                    border: '2px solid #ddd', 
                    padding: '12px', 
                    fontWeight: 'bold', 
                    textAlign: 'center'
                  }}>{day}</td>
                  {times.map((time, index) => (
                    <td key={index} style={{ 
                      border: '2px solid #ddd', 
                      padding: '7px', 
                      fontSize: '16px',
                      fontWeight:'bold',
                      textAlign: 'center'
                    }}>
                      {Object.values(schedule).flatMap(faculty =>
                        Object.values(faculty).flatMap(subjects =>
                          Array.isArray(subjects) ? subjects.filter(
                            item => item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                          ) : []
                        )
                      ).map((item, idx) => (
                        <div key={idx}>
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
      }
    />
  );
};

export default Timetable;
