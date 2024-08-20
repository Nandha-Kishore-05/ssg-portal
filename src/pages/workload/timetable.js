// import React, { useState, useEffect } from 'react';
// import axios from 'axios';
// import { useParams } from 'react-router-dom';
// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';

// const SavedTimetable = () => {
//   const { departmentID,semesterID  } = useParams();
//   const [schedule, setSchedule] = useState([]);
//   const [days, setDays] = useState([]);
//   const [times, setTimes] = useState([]);
//   const [venue, setVenue] = useState('');
//   const [department,setDepartment] = useState();
//   const [semester,setSemester] = useState();

//   useEffect(() => {
//     const fetchSchedule = async () => {
//       if (!departmentID || !semesterID) {
//         console.error('Department ID and Semester ID are required');
//         return;
//       }

//       try {
//         const response = await axios.get(`http://localhost:8080/timetable/saved/${departmentID}/${semesterID}`);
//         const data = response.data;

//         console.log('Fetched data:', data);

//         const allDays = new Set();
//         const allTimes = new Set();
//         let venueSet = new Set();

//         // Extract unique days, time slots, and venue
//         data.forEach(item => {
//           allDays.add(item.day_name);
//           allTimes.add(`${item.start_time} - ${item.end_time}`);
//           venueSet.add(item.classroom); // Collect unique venues
//         });

//         const sortedDays = Array.from(allDays).sort((a, b) => {
//           const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
//           return order.indexOf(a) - order.indexOf(b);
//         });

//         const sortedTimes = Array.from(allTimes).sort((a, b) => a.localeCompare(b, undefined, { numeric: true }));

//         setDays(sortedDays);
//         setTimes(sortedTimes);
//         setSchedule(data);
//         setVenue(Array.from(venueSet).join(', ')); // Set venue
//       } catch (error) {
//         console.error('Error fetching timetable:', error);
//       }
//     };

//     fetchSchedule();
//   }, [departmentID,semesterID]);

//   return (
//     <AppLayout
//       rId={3}
//       title="Venue Table"
//       body={
//         <>
//         <CustomSelect
                 
//                   placeholder="DEPARTMENT"
//                   value={department}
//                   onChange={ setDepartment}
//                   options={[
//                     { label: "COMPUTER TECHNOLOGY", value: 1 },
//                     { label: "BI0 TECHNOLOGY", value: 2 },
//                   ]}
//                 />
//                 <CustomSelect
                 
//                  placeholder="SEMESTER"
//                value={semester}
//                onChange={ setSemester}
//                  options={[
//                    { label: "S1", value: 1 },
//                    { label: "S3", value: 3 },
//                    { label: "S5", value: 5 },
//                  ]}
//                /><br />
//                <center>
//                 <CustomButton
//               width="150"
//               label="View Timetable"
             
//             />
//             </center>
//         <div style={{ 
//           backgroundColor: '#fff', 
//           padding: '20px', 
//           borderRadius: '8px', 
//           boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)', 
//           margin: '20px 0'
//         }}>
//           <div style={{display:'flex',flexDirection:'row',justifyContent:'space-between',marginBottom:'13px'}}>
//           <h2 style={{fontSize:'20px',marginTop:'5px'}}>Semester : S{semesterID}</h2>
//             <h2 style={{fontSize:'20px',marginTop:'5px'}}>Venue: {venue || 'Not Available'}</h2>
//           </div>
//           <table style={{ 
//             width: '100%', 
//             borderCollapse: 'collapse', 
//             backgroundColor: '#fff', 
//             border: '2px solid #ddd',
//             fontSize: '16px',
//             minHeight: '600px'
//           }}>
//             <thead>
//               <tr style={{ backgroundColor: '#f4f4f4' }}>
//                 <th style={{ 
//                   border: '2px solid #ddd', 
//                   padding: '12px', 
//                   textAlign: 'center'
//                 }}>Day/Time</th>
//                 {times.map((time, index) => (
//                   <th key={index} style={{ 
//                     border: '2px solid #ddd', 
//                     padding: '12px', 
//                     backgroundColor: '#f9f9f9', 
//                     textAlign: 'center'
//                   }}>
//                     {time}
//                   </th>
//                 ))}
//               </tr>
//             </thead>
//             <tbody>
//               {days.map((day) => (
//                 <tr key={day}>
//                   <td style={{ 
//                     border: '2px solid #ddd', 
//                     padding: '12px', 
//                     fontWeight: 'bold', 
//                     textAlign: 'center'
//                   }}>{day}</td>
//                   {times.map((time, index) => {
//                     const classes = schedule.filter(item =>
//                       item.day_name === day && `${item.start_time} - ${item.end_time}` === time
//                     );
//                     return (
//                       <td key={index} style={{ 
//                         border: '2px solid #ddd', 
//                         padding: '7px', 
//                         fontSize: '16px',
//                         fontWeight:'bold',
//                         textAlign: 'center'
//                       }}>
//                         {classes.length > 0 ? (
//                           classes.map((item, idx) => (
//                             <div key={idx}>
//                               <div>{item.subject_name}</div>
//                               <div>{item.faculty_name}</div>
//                             </div>
//                           ))
//                         ) : (
//                           <div>No classes</div>
//                         )}
//                       </td>
//                     );
//                   })}
//                 </tr>
//               ))}
//             </tbody>
//           </table>
//         </div>
//         </>
//       }
//     />
//   );
// };

//  export default SavedTimetable;

// import React, { useState, useEffect } from 'react';
// import axios from 'axios';
// import { useParams, useNavigate } from 'react-router-dom';
// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';

// const SavedTimetable = () => {
//   const navigate = useNavigate();
//   const { departmentID, semesterID } = useParams();
//   const [schedule, setSchedule] = useState([]);
//   const [days, setDays] = useState([]);
//   const [times, setTimes] = useState([]);
//   const [venue, setVenue] = useState('');
//   const [department, setDepartment] = useState('');
//   const [semester, setSemester] = useState('');

//   useEffect(() => {
//     const fetchSchedule = async () => {
//       if (!departmentID || !semesterID) {
//         console.error('Department ID and Semester ID are required');
//         return;
//       }

//       try {
//         const response = await axios.get(`http://localhost:8080/timetable/saved/${departmentID}/${semesterID}`);
//         const data = response.data;

//         console.log('Fetched data:', data);

//         const allDays = new Set();
//         const allTimes = new Set();
//         let venueSet = new Set();

//         // Extract unique days, time slots, and venue
//         data.forEach(item => {
//           allDays.add(item.day_name);
//           allTimes.add(`${item.start_time} - ${item.end_time}`);
//           venueSet.add(item.classroom); // Collect unique venues
//         });

//         const sortedDays = Array.from(allDays).sort((a, b) => {
//           const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
//           return order.indexOf(a) - order.indexOf(b);
//         });

//         const sortedTimes = Array.from(allTimes).sort((a, b) => a.localeCompare(b, undefined, { numeric: true }));

//         setDays(sortedDays);
//         setTimes(sortedTimes);
//         setSchedule(data);
//         setVenue(Array.from(venueSet).join(', ')); // Set venue
//       } catch (error) {
//         console.error('Error fetching timetable:', error);
//       }
//     };

//     fetchSchedule();
//   }, [departmentID, semesterID]);

//   const handleViewTimetable = () => {
//     console.log('Department:', department);
//     console.log('Semester:', semester);
//     if (department && semester) {
//       navigate(`/timetable/saved/${department}/${semester}`);
//     } else {
//       console.error('Please select both department and semester');
//     }
//   };

//   return (
//     <AppLayout
//       rId={3}
//       title="Venue Table"
//       body={
//         <>
//           <CustomSelect
//             placeholder="DEPARTMENT"
//             value={department}
//             onChange={setDepartment}
//             options={[
//               { label: "COMPUTER TECHNOLOGY", value: 1 },
//               { label: "BIO TECHNOLOGY", value: 2 },
//             ]}
//           />
//           <CustomSelect
//             placeholder="SEMESTER"
//             value={semester}
//             onChange={setSemester}
//             options={[
//               { label: "S1", value: 1 },
//               { label: "S3", value: 3 },
//               { label: "S5", value: 5 },
//             ]}
//           />
//           <br />
//           <center>
//             <CustomButton
//               width="150"
//               label="View Timetable"
//               onClick={handleViewTimetable}
//             />
//           </center>
//           <div style={{
//             backgroundColor: '#fff',
//             padding: '20px',
//             borderRadius: '8px',
//             boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)',
//             margin: '20px 0'
//           }}>
//             <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', marginBottom: '13px' }}>
//               <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : S{semesterID}</h2>
//               <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Venue: {venue || 'Not Available'}</h2>
//             </div>
//             <table style={{
//               width: '100%',
//               borderCollapse: 'collapse',
//               backgroundColor: '#fff',
//               border: '2px solid #ddd',
//               fontSize: '16px',
//               minHeight: '600px'
//             }}>
//               <thead>
//                 <tr style={{ backgroundColor: '#f4f4f4' }}>
//                   <th style={{
//                     border: '2px solid #ddd',
//                     padding: '12px',
//                     textAlign: 'center'
//                   }}>Day/Time</th>
//                   {times.map((time, index) => (
//                     <th key={index} style={{
//                       border: '2px solid #ddd',
//                       padding: '12px',
//                       backgroundColor: '#f9f9f9',
//                       textAlign: 'center'
//                     }}>
//                       {time}
//                     </th>
//                   ))}
//                 </tr>
//               </thead>
//               <tbody>
//                 {days.map((day) => (
//                   <tr key={day}>
//                     <td style={{
//                       border: '2px solid #ddd',
//                       padding: '12px',
//                       fontWeight: 'bold',
//                       textAlign: 'center'
//                     }}>{day}</td>
//                     {times.map((time, index) => {
//                       const classes = schedule.filter(item =>
//                         item.day_name === day && `${item.start_time} - ${item.end_time}` === time
//                       );
//                       return (
//                         <td key={index} style={{
//                           border: '2px solid #ddd',
//                           padding: '7px',
//                           fontSize: '16px',
//                           fontWeight: 'bold',
//                           textAlign: 'center'
//                         }}>
//                           {classes.length > 0 ? (
//                             classes.map((item, idx) => (
//                               <div key={idx}>
//                                 <div>{item.subject_name}</div>
//                                 <div>{item.faculty_name}</div>
//                               </div>
//                             ))
//                           ) : (
//                             <div>No classes</div>
//                           )}
//                         </td>
//                       );
//                     })}
//                   </tr>
//                 ))}
//               </tbody>
//             </table>
//           </div>
//         </>
//       }
//     />
//   );
// };

// export default SavedTimetable;

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';

const SavedTimetable = (props) => {
  const navigate = useNavigate();
  
  const [schedule, setSchedule] = useState([]);
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
        const response = await axios.get(`http://localhost:8080/timetable/saved/${props.departmentID}/${props.semesterID}`);
        const data = response.data;

        console.log('Fetched data:', data);

        const allDays = new Set();
        const allTimes = new Set();
        let venueSet = new Set();

        // Extract unique days, time slots, and venue
        data.forEach(item => {
          allDays.add(item.day_name);
          allTimes.add(`${item.start_time} - ${item.end_time}`);
          venueSet.add(item.classroom); // Collect unique venues
        });

        const sortedDays = Array.from(allDays).sort((a, b) => {
          const order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
          return order.indexOf(a) - order.indexOf(b);
        });

        const sortedTimes = Array.from(allTimes).sort((a, b) => a.localeCompare(b, undefined, { numeric: true }));

        setDays(sortedDays);
        setTimes(sortedTimes);
        setSchedule(data);
        setVenue(Array.from(venueSet).join(', ')); // Set venue
      } catch (error) {
        console.error('Error fetching timetable:', error);
      }
    };

    fetchSchedule();
  }, [props.departmentID, props.semesterID]);



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
              <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : S{props.semesterID}</h2>
              <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Venue: {venue || 'Not Available'}</h2>
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
                          fontWeight: 'bold',
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
        </>

 );
              };

export default SavedTimetable;
