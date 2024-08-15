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
                
//                   options={[
//                     { label: "COMPUTER TECHNOLOGY", value: 1 },
//                     { label: "BIO TECHNOLOGY", value: 2 },
//                   ]}
//                 />
//                 <CustomSelect
                 
//                  placeholder="SEMESTER"
               
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
// import { useParams } from 'react-router-dom';
// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';

// const SavedTimetable = () => {
//   const { departmentID, semesterID } = useParams();
//   const [selectedDepartment, setSelectedDepartment] = useState(departmentID);
//   const [selectedSemester, setSelectedSemester] = useState(semesterID);
//   const [schedule, setSchedule] = useState([]);
//   const [days, setDays] = useState([]);
//   const [times, setTimes] = useState([]);
//   const [venue, setVenue] = useState('');

//   useEffect(() => {
//     const fetchSchedule = async () => {
//       if (!selectedDepartment || !selectedSemester) {
//         console.error('Department ID and Semester ID are required');
//         return;
//       }

//       try {
//         const response = await axios.get(`http://localhost:8080/timetable/saved/${selectedDepartment}/${selectedSemester}`);
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
//   }, [selectedDepartment, selectedSemester]);

//   const handleViewTimetable = () => {
//     // Trigger fetching new timetable based on selected department and semester
//     if (selectedDepartment && selectedSemester) {
//       console.log(`Fetching timetable for Department ${selectedDepartment}, Semester ${selectedSemester}`);
//     } else {
//       console.error('Please select both department and semester.');
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
//             value={selectedDepartment}
//             options={[
//               { label: "COMPUTER TECHNOLOGY", value: 1 },
//               { label: "BIO TECHNOLOGY", value: 2 },
//             ]}
//             onChange={setSelectedDepartment}
//           />
//           <CustomSelect
//             placeholder="SEMESTER"
//             value={selectedSemester}
//             onChange={setSelectedSemester}
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
//               func={handleViewTimetable}
//             />
//           </center>
//           <div style={{ 
//             backgroundColor: '#fff', 
//             padding: '20px', 
//             borderRadius: '8px', 
//             boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)', 
//             margin: '20px 0'
//           }}>
//             <div style={{display:'flex',flexDirection:'row',justifyContent:'space-between',marginBottom:'13px'}}>
//               <h2 style={{fontSize:'20px',marginTop:'5px'}}>Semester : S{selectedSemester}</h2>
//               <h2 style={{fontSize:'20px',marginTop:'5px'}}>Venue: {venue || 'Not Available'}</h2>
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
//                           fontWeight:'bold',
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

// import React, { useState, useEffect } from 'react';
// import axios from 'axios';
// import AppLayout from '../../layout/layout';
// import './workload.css';
// import CustomSelect from '../../components/select';
// import CustomButton from '../../components/button';

// const SavedTimetable = () => {
//   const [selectedDepartment, setSelectedDepartment] = useState(null);
//   const [selectedSemester, setSelectedSemester] = useState(null);
//   const [schedule, setSchedule] = useState([]);
//   const [days, setDays] = useState([]);
//   const [times, setTimes] = useState([]);
//   const [venue, setVenue] = useState('');

//   useEffect(() => {
//     const fetchSchedule = async () => {
//       if (!selectedDepartment?.value || !selectedSemester?.value) {
//         console.error('Department ID and Semester ID are required');
//         return;
//       }

//       try {
//         const response = await axios.get(
//           `http://localhost:8080/timetable/saved/${selectedDepartment.value}/${selectedSemester.value}`
//         );
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

//     if (selectedDepartment && selectedSemester) {
//       fetchSchedule();
//     }
//   }, [selectedDepartment, selectedSemester]);

//   const handleViewTimetable = () => {
//     // Trigger fetching new timetable based on selected department and semester
//     if (selectedDepartment && selectedSemester) {
//       console.log(`Fetching timetable for Department ${selectedDepartment.value}, Semester ${selectedSemester.value}`);
//     } else {
//       console.error('Please select both department and semester.');
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
//             value={selectedDepartment}
//             options={[
//               { label: "COMPUTER TECHNOLOGY", value: 1 },
//               { label: "BIO TECHNOLOGY", value: 2 },
//             ]}
//             onChange={setSelectedDepartment}
//           />
//           <CustomSelect
//             placeholder="SEMESTER"
//             value={selectedSemester}
//             onChange={setSelectedSemester}
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
//               func={handleViewTimetable}
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
//               <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : {selectedSemester?.label}</h2>
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
import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';

const SavedTimetable = ({ departmentID, semesterID }) => {
  const [selectedDepartment, setSelectedDepartment] = useState(
    departmentID ? { label: "COMPUTER TECHNOLOGY", value: departmentID } : null
  );
  const [selectedSemester, setSelectedSemester] = useState(
    semesterID ? { label: `S${semesterID}`, value: semesterID } : null
  );
  const [schedule, setSchedule] = useState([]);
  const [days, setDays] = useState([]);
  const [times, setTimes] = useState([]);
  const [venue, setVenue] = useState('');

  useEffect(() => {
    const fetchSchedule = async () => {
      if (!selectedDepartment?.value || !selectedSemester?.value) {
        console.error('Department ID and Semester ID are required');
        return;
      }

      try {
        const response = await axios.get(
          `http://localhost:8080/timetable/saved/${selectedDepartment.value}/${selectedSemester.value}`
        );
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
  }, [selectedDepartment, selectedSemester]);

  const handleViewTimetable = () => {
    // Trigger fetching new timetable based on selected department and semester
    if (selectedDepartment && selectedSemester) {
      console.log(`Fetching timetable for Department ${selectedDepartment.value}, Semester ${selectedSemester.value}`);
    } else {
      console.error('Please select both department and semester.');
    }
  };

  return (
    <AppLayout
      rId={3}
      title="Venue Table"
      body={
        <>
          <CustomSelect
            placeholder="DEPARTMENT"
            value={selectedDepartment}
            options={[
              { label: "COMPUTER TECHNOLOGY", value: 1 },
              { label: "BIO TECHNOLOGY", value: 2 },
            ]}
            onChange={setSelectedDepartment}
          />
          <CustomSelect
            placeholder="SEMESTER"
            value={selectedSemester}
            onChange={setSelectedSemester}
            options={[
              { label: "S1", value: 1 },
              { label: "S3", value: 3 },
              { label: "S5", value: 5 },
            ]}
          />
          <br />
          <center>
            <CustomButton
              width="150"
              label="View Timetable"
              func={handleViewTimetable}
            />
          </center>
          <div style={{
            backgroundColor: '#fff',
            padding: '20px',
            borderRadius: '8px',
            boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)',
            margin: '20px 0'
          }}>
            <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', marginBottom: '13px' }}>
              <h2 style={{ fontSize: '20px', marginTop: '5px' }}>Semester : {selectedSemester?.label}</h2>
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
      }
    />
  );
};

export default SavedTimetable;
