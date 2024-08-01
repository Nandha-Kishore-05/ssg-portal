import React from 'react';
import CustomButton from '../../components/button';

const Timetable = ({ days = [], times = [], schedule = {}, classroom, onSave }) => {
    const handleSave = () => {
        const timetableData = [];
      
        days.forEach(day => {
          times.forEach(time => {
            const entry = schedule[day]?.find(item => `${item.start_time} - ${item.end_time}` === time);
            if (entry) {
              console.log("Entry Found:", entry);
      
              const data = {
                day_name: entry.day_name, // Ensure mapping matches data structure
                start_time: entry.start_time,
                end_time: entry.end_time,
                subject_name: entry.subject_name, // Ensure mapping matches data structure
                faculty_name: entry.faculty_name, // Ensure mapping matches data structure
                classroom: entry.classroom
              };
      
              console.log("Formatted Data:", data);
              timetableData.push(data);
            }
          });
        });
      
        console.log("Final Timetable Data to Save:", timetableData);
        onSave(timetableData);
      };
      
  return (
    <div style={{ 
      backgroundColor: '#fff', 
      padding: '20px', 
      borderRadius: '8px', 
      boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)', 
      margin: '20px 0'
    }}>
      <div style={{display:'flex',flexDirection:'row',justifyContent:'space-between',marginBottom:'13px'}}>
        <h2 style={{fontSize:'20px',marginTop:'5px'}}>Venue : {classroom}</h2>
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
                  {schedule[day] && schedule[day].find(
                    item => `${item.start_time} - ${item.end_time}` === time
                  ) ? (
                    <>
                      <div>{schedule[day].find(
                        item => `${item.start_time} - ${item.end_time}` === time
                      ).subject_name}</div>
                      <div>{schedule[day].find(
                        item => `${item.start_time} - ${item.end_time}` === time
                      ).faculty_name}</div>
                    </>
                  ) : ''}
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
