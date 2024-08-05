import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';
import AppLayout from '../../layout/layout';

const SavedTimetable = () => {
  const { departmentID } = useParams();
  const [schedule, setSchedule] = useState([]);
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);

  useEffect(() => {
    const fetchSchedule = async () => {
      if (!departmentID) {
        console.error('Department ID is required');
        return;
      }

      try {
        const response = await axios.get(`http://localhost:8080/timetable/saved/${departmentID}`);
        const data = response.data;

        console.log('Fetched data:', data);

        const allDays = new Set();
        const allTimes = new Set();

        // Extract unique days and time slots
        data.forEach(item => {
          allDays.add(item.day_name);
          allTimes.add(`${item.start_time} - ${item.end_time}`);
        });

        const sortedDays = Array.from(allDays).sort((a, b) => {
          const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
          return order.indexOf(a) - order.indexOf(b);
        });

        const sortedTimes = Array.from(allTimes).sort((a, b) => a.localeCompare(b, undefined, { numeric: true }));

        setDays(sortedDays);
        setTimes(sortedTimes);
        setSchedule(data);
      } catch (error) {
        console.error('Error fetching timetable:', error);
      }
    };

    fetchSchedule();
  }, [departmentID]);

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
                  {times.map((time, index) => {
                    const classes = schedule.filter(item =>
                      item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                    );
                    return (
                      <td key={index} style={{ 
                        border: '2px solid #ddd', 
                        padding: '7px', 
                        fontSize: '16px',
                        fontWeight:'bold',
                        textAlign: 'center'
                      }}>
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
        </div>
      }
    />
  );
};

export default SavedTimetable;
