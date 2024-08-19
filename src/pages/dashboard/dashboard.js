// src/pages/Dashboard.js
import React from 'react';
import './style.css'; // Import your custom CSS
import AppLayout from '../../layout/layout';
import { BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';

const dailyData = [
  { day: '1', hours: 5 },
  { day: '2', hours: 7 },
  { day: '3', hours: 6 },
  { day: '4', hours: 8 },
  { day: '5', hours: 4 },
  { day: '1', hours: 5 },
  { day: '2', hours: 7 },
  { day: '3', hours: 6 },
  { day: '4', hours: 8 },
  { day: '5', hours: 4 },
  { day: '1', hours: 5 },
  { day: '2', hours: 7 },
  { day: '3', hours: 6 },
  { day: '4', hours: 8 },
  { day: '5', hours: 4 },
  { day: '1', hours: 5 },
  { day: '2', hours: 7 },
  { day: '3', hours: 6 },
  { day: '4', hours: 8 },
  { day: '5', hours: 4 },
  { day: '1', hours: 5 },
  { day: '2', hours: 7 },
  { day: '3', hours: 6 },
  { day: '4', hours: 8 },
  { day: '5', hours: 4 },
  // Add more data as needed
];

const monthlyData = [
  { month: 'Jan', hours: 40 },
  { month: 'Feb', hours: 35 },
  { month: 'Mar', hours: 50 },
  { month: 'Apr', hours: 45 },
  { month: 'May', hours: 55 },
  { month: 'Jun', hours: 60 },
  { month: 'Jul', hours: 70 },
  { month: 'Aug', hours: 65 },
  { month: 'Sep', hours: 50 },
  { month: 'Oct', hours: 45 },
  { month: 'Nov', hours: 40 },
  { month: 'Dec', hours: 55 },
  // Add more data as needed
];

function Dashboard() {
  return (
    <AppLayout
      rId={1}
      title="Dashboard"
      body={
        <>
        <div className='grid-layout'>
          <div className='grid'>
            <div className='working-hours'>
              <div className='header'>
                Daily Working Hours
              </div>
              <BarChart
                width={1350}
                height={300}
                data={dailyData}
                margin={{ top: 40, right: 10, left: 10, bottom: 5 }}
              >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="day" />
                <YAxis tickLine={false} />
                <Tooltip />
                <Bar dataKey="hours" fill="blue" />
              </BarChart>
            </div>
          </div>
          </div>
          <div className='grid-2'>
            <div className='mwh'>
              <div className='header-2'>
                Monthly Working Hours
              </div>
             
            </div>
            <div className='grid-3'>
              <div className='header-2'>Student Log</div>
            </div>
          </div>
        </>
      }
    />
  );
}

export default Dashboard;
