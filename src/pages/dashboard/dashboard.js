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
       </>
      }
    />
  );
}

export default Dashboard;
