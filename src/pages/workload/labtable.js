import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './lab.css'; // Import the CSS file

const LabTimetable = (props) => {
  const [schedule, setSchedule] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchTimetable = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/timetable/lab/${props.subjectName}`);
        setSchedule(response.data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (props.subjectName) {
      fetchTimetable();
    }
  }, [props.subjectName]);

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
  const sortedDays = allDays.sort((a, b) => dayOrder.indexOf(a) - dayOrder.indexOf(b));

  // Extract unique time slots and sort based on predefined order
  const allTimes = Array.from(new Set(schedule.map(item => `${item.start_time} - ${item.end_time}`)));
  const sortedTimes = allTimes.sort((a, b) => timeOrder.indexOf(a) - timeOrder.indexOf(b));

  return (
    <div className="container">
      <div className="header-i">
        <h2>Lab Name : {props.subjectName}</h2>
      </div>
      <table className="table">
        <thead>
          <tr>
            <th className="day-time">Day/Time</th>
            {sortedTimes.map((time, index) => (
              <th key={index} className="time">
                {time}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {sortedDays.map(day => (
            <tr key={day}>
              <td className="day">{day}</td>
              {sortedTimes.map((time, index) => {
                const classes = schedule.filter(item =>
                  item.day_name === day && `${item.start_time} - ${item.end_time}` === time
                );
                return (
                  <td key={index} className="subject">
                    {classes.length > 0 ? (
                      classes.map((item, idx) => (
                        <div key={idx}>
                          <div>{item.subject_name}</div>
                          <div>{item.faculty_name}</div>
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
  );
};

export default LabTimetable;
