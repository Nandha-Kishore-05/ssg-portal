import React, { useState, useEffect } from 'react';
import axios from 'axios';



const FacultyTimetable = (props) => {

  const [schedule, setSchedule] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchTimetable = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/timetable/faculty/${props.facultyName}`);
        setSchedule(response.data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (props.facultyName) {
      fetchTimetable();
    }
  }, [props.facultyName]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;

  // Define the desired order for days and time slots
  const dayOrder = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
  const timeOrder = [
    "08:45:00 - 09:35:00",
    "09:35:00 - 10:25:00",
    "10:40:00 - 11:30:00",
    "13:45:00 - 14:35:00",
    "14:35:00 - 15:25:00",
    "15:40:00 - 16:30:00"
  ];

  // Extract unique days and sort based on predefined order
  const allDays = Array.from(new Set(schedule.map(item => item.day_name)));
  const sortedDays = allDays.sort((a, b) => {
    return dayOrder.indexOf(a) - dayOrder.indexOf(b);
  });

  // Extract unique time slots and sort based on predefined order
  const allTimes = Array.from(new Set(schedule.map(item => `${item.start_time} - ${item.end_time}`)));
  const sortedTimes = allTimes.sort((a, b) => {
    return timeOrder.indexOf(a) - timeOrder.indexOf(b);
  });

  return (
    
        <>
        <div style={{ 
          backgroundColor: '#fff', 
          padding: '20px', 
          borderRadius: '8px', 
          boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)', 
          margin: '20px 0'
        }}>
          <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', marginBottom: '13px' }}>
            <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Faculty: {props.facultyName}</h2>
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
                {sortedTimes.map((time, index) => (
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
              {sortedDays.map(day => (
                <tr key={day}>
                  <td style={{ 
                    border: '2px solid #ddd', 
                    padding: '12px', 
                    fontWeight: 'bold', 
                    textAlign: 'center'
                  }}>{day}</td>
                  {sortedTimes.map((time, index) => {
                    const classes = schedule.filter(item =>
                      item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                    );
                    return (
                      <td key={index} style={{ 
                        border: '2px solid #ddd', 
                        padding: '7px', 
                        fontSize: '16px',
                        fontWeight: 'bold',
                        textAlign: 'center'
                      }}>
                        {classes.length > 0 ? (
                          classes.map((item, idx) => (
                            <div key={idx}>
                              <div>{item.classroom}</div>
                          
                              <div>S{item.semester_id}</div>
                            </div>
                          ))
                        ) : (
                          <div>-</div>
                        )}
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        </>
     
  );
};

export default FacultyTimetable;
