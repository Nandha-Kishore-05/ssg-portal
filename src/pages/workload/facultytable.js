import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { utils, writeFile } from 'xlsx'; 
import './Fac.css';
import CustomButton from '../../components/button';

const FacultyTimetable = (props) => {
  console.log(props)
  const [schedule, setSchedule] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchTimetable = async () => {
      try {
       
        const response = await axios.get(`http://localhost:8080/timetable/faculty/${props.facultyName}/${props.academicYearID}`);
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
  }, [props.facultyName, props.academicYearID]); // Added props.academicYearID to dependency array
  
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;

  const dayOrder = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
  const timeOrder = [
    "08:45:00 - 09:35:00",
    "09:35:00 - 10:25:00",
    "10:40:00 - 11:30:00",
    "11:30:00 - 12:20:00",
    "13:30:00 - 14:20:00",
    "14:20:00 - 15:10:00",
    "15:25:00 - 16:30:00"
  ];

  const allDays = Array.from(new Set(schedule.map(item => item.day_name)));
  const sortedDays = allDays.sort((a, b) => {
    return dayOrder.indexOf(a) - dayOrder.indexOf(b);
  });

  const allTimes = Array.from(new Set(schedule.map(item => `${item.start_time} - ${item.end_time}`)));
  const sortedTimes = allTimes.sort((a, b) => {
    return timeOrder.indexOf(a) - timeOrder.indexOf(b);
  });


  const downloadTimetableAsExcel = () => {
    const wsData = [
      ['Day/Time', ...sortedTimes], 
    ];

    sortedDays.forEach(day => {
      const row = [day];
      sortedTimes.forEach(time => {
        const classes = schedule.filter(item =>
          item.day_name === day && `${item.start_time} - ${item.end_time}` === time
        );
        if (classes.length > 0) {
          row.push(classes.map(item => `S${item.semester_id} - ${item.classroom}`).join('\n'));
        } else {
          row.push('-');
        }
      });
      wsData.push(row);
    });

    const ws = utils.aoa_to_sheet(wsData); 
    const wb = utils.book_new(); 
    utils.book_append_sheet(wb, ws, 'Faculty Timetable'); 

    writeFile(wb, `${props.facultyName}-Timetable.xlsx`); 
  };

  return (
    <div className="container-2">
      <div className="header-k">
        <div className="header-info">
          <h2>Faculty: {props.facultyName}</h2>
        </div>
        <div className="buttons">
          <CustomButton
            width="150"
            label="Download Timetable"
            onClick={downloadTimetableAsExcel} // Call the download function on button click
          />
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
                           {/* <div >{item.subject_name}</div> */}
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
