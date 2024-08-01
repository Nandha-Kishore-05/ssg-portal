import React, { useState, useEffect } from 'react';
import axios from 'axios';
import AppLayout from '../../layout/layout';
import { useParams } from 'react-router-dom';
import Timetable from './timetable';

const sortDays = (days) => {
  const dayOrder = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  return dayOrder.filter(day => days.includes(day));
};

const Workload = () => {
  const { departmentID } = useParams();
  const [data, setData] = useState({ days: [], times: [], schedule: {}, classroom: '', department: '' });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const loadTimetableData = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/timetable/${departmentID}`);
        console.log("Fetched data:", response.data);

        const schedule = response.data;
        const days = sortDays(Object.keys(schedule));
        const times = Array.from(new Set(schedule[days[0]].map(item => `${item.start_time} - ${item.end_time}`)));

        const classroom = schedule[days[0]][0]?.classroom || 'Unknown';
        const department = schedule[days[0]][0]?.department || 'Unknown';
        console.log("Classroom:", classroom);
        setData({ days, times, schedule, classroom, department });
      } catch (err) {
        console.error("Error fetching timetable data:", err);
        setError("Unable to retrieve timetable data");
      } finally {
        setLoading(false);
      }
    };

    loadTimetableData();
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

  if (loading) return <p>Loading...</p>;
  if (error) return <p>{error}</p>;

  return (
    <AppLayout
      rId={2}
      title="Time Table"
      body={
        <Timetable 
          days={data.days}
          times={data.times} 
          schedule={data.schedule} 
          classroom={data.classroom}
          department={data.department}
          onSave={handleSaveTimetable}
        />
      }
    />
  );
};

export default Workload;
