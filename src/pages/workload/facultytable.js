import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './Fac.css';

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


  const dayOrder = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
  const timeOrder = [
    "08:45:00 - 09:35:00",
    "09:35:00 - 10:25:00",
    "10:40:00 - 11:30:00",
    "13:45:00 - 14:35:00",
    "14:35:00 - 15:25:00",
    "15:40:00 - 16:30:00"
  ];


  const allDays = Array.from(new Set(schedule.map(item => item.day_name)));
  const sortedDays = allDays.sort((a, b) => {
    return dayOrder.indexOf(a) - dayOrder.indexOf(b);
  });


  const allTimes = Array.from(new Set(schedule.map(item => `${item.start_time} - ${item.end_time}`)));
  const sortedTimes = allTimes.sort((a, b) => {
    return timeOrder.indexOf(a) - timeOrder.indexOf(b);
  });

  return (
    <div className="container">
      <div className="header-k">
        <div className="header-info">
          <h2>Faculty: {props.facultyName}</h2>
        </div>
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
  );
};

export default FacultyTimetable;
