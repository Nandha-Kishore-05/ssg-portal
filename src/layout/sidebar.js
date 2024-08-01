import "./style.css";

import { DashboardRounded, LogoutRounded } from "@mui/icons-material";
import { Link } from "react-router-dom";
import { useState } from "react";

function SideBar(props) {
  const [openSubMenuId, setOpenSubMenuId] = useState(null);

  const menu = [
    {
      id: 1,
      icon: <DashboardRounded />,
      label: "Dashboard",
      path: "/dashboard",
    },
    {
        id: 2,
        icon: <DashboardRounded />,
        label: "TimeTable",
        path: "/timetable/1",
      },
    


  
    
  ];

  const toggleSubMenu = (id) => {
    if (openSubMenuId === id) {
      setOpenSubMenuId(null);
    } else {
      setOpenSubMenuId(id);
    }
  };

  // useEffect(() => {
  //   setOpenSubMenuId(props.rId);
  // }, []); 

  return (
    <div className="app-sidebar">
      <div className="sidebar-header">
      
        <h2 style={{fontSize:"28"}}> TT PORTAL </h2>
      </div>

      <div className="sidebar-menu">
        {menu.map((item, i) => (
          <div key={i}>
            <Link
              style={{ textDecoration: "none" }}
              to={
                item.submenu === undefined || item.submenu.length === 0
                  ? item.path
                  : ""
              }
            >
              <div
                onClick={() => toggleSubMenu(item.id)}
                className={
                  props.rId === item.id
                    ? "sidebar-menu-item sidebar-selected"
                    : "sidebar-menu-item"
                }
              >
                {item.icon}
                <h4>{item.label}</h4>
              </div>
            </Link>
            {item.submenu && ( item.id === props.rId || openSubMenuId === item.id) && (
              <div className="submenu">
                {item.submenu.map((subitem, j) => (
                  <Link
                    key={j}
                    style={{ textDecoration: "none" }}
                    to={subitem.path}
                  >
                    <div className="submenu-item">
                      <div className="menu-dot"></div>
                      <h4>{subitem.label}</h4>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

export default SideBar;
