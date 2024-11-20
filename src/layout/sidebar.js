import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import {
  DashboardRounded,
  Event as EventIcon,
  GroupRounded,
  ScienceRounded,
  ScheduleRounded,
  Edit as EditIcon,
} from '@mui/icons-material';
import LogoutIcon from '@mui/icons-material/Logout';
import './style.css';
import { Book } from "@mui/icons-material";
import FolderIcon from '@mui/icons-material/Folder';
import RoomIcon from '@mui/icons-material/Room';
import SchoolIcon from '@mui/icons-material/School';
import PersonIcon from '@mui/icons-material/Person';

function SideBar(props) {
  const [openSubMenuId, setOpenSubMenuId] = useState(null);
  const [menu, setMenu] = useState([]);

  useEffect(() => {
    const fetchResources = async () => {
      const authToken = localStorage.getItem('authToken'); 

      try {
        const response = await axios.get('http://localhost:8080/getResource', {
          headers: {
            Authorization: `${authToken}`, 
          },
        });

      
        const fetchedMenu = response.data.data.map((resource) => ({
          id: resource.id,
          label: resource.name,
          path: resource.path,
        }));
        setMenu(fetchedMenu);
      } catch (error) {
        console.error('Error fetching resources:', error);
      }
    };

    fetchResources();
  }, []);


  const getIcon = (label) => {
    switch (label) {
      case 'Dashboard':
        return <DashboardRounded />;
      case 'Time Table':
        return <EventIcon />;
      case 'Faculty Table':
        return <GroupRounded />;
      case 'Lab Table':
        return <ScienceRounded />;
      case 'Period Allocation':
        return <ScheduleRounded />;
      case 'Manual Entry':
        return <EditIcon />;
        case 'Subject Entry':
          return   <Book />
        case 'Log Out':
          return  <LogoutIcon />
        case 'Master Timetable':
          return  <FolderIcon />
        case 'Venue Table':
            return  <RoomIcon  />
        case 'Student Allocation':
            return  <SchoolIcon  />
        case 'Student Timetable':
              return  <PersonIcon  />
              case 'Lab Entry':
                return <EditIcon />;
      default:
        return null; 
    }
  };

  const toggleSubMenu = (id) => {
    setOpenSubMenuId(openSubMenuId === id ? null : id);
  };

  return (
    <div className="app-sidebar">
      <div className="sidebar-header">
        <h2 style={{ fontSize: '28px', color: 'white' }}>TT PORTAL</h2>
      </div>
      <div className="sidebar-menu">
        {menu.map((item, i) => (
          <div key={i}>
            <Link style={{ textDecoration: 'none' }} to={item.path}>
              <div
                onClick={() => toggleSubMenu(item.id)}
                className={props.rId === item.id ? 'sidebar-menu-item sidebar-selected' : 'sidebar-menu-item'}
              >
                {getIcon(item.label)}
                <h4>{item.label}</h4>
              </div>
            </Link>
          </div>
        ))}
      </div>
    </div>
  );
}

export default SideBar;