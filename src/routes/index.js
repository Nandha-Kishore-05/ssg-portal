import Dashboard from "../pages/dashboard/dashboard";
import Timetable from "../pages/workload/workload";

import SavedTimetable from "../pages/workload/timetable";
import FacultyTimetable from "../pages/workload/facultytable";



const routes = [
  
  {
    path: "/dashboard",
    element: <Dashboard />,
  },
  {
    path: "/timetable/:departmentID",
    element: <Timetable />,
  },
  {
    path: "/timetable/saved/:departmentID",
    element: <SavedTimetable />,
  },
  {
    path: "/timetable/faculty/:facultyName",
    element: <FacultyTimetable />,
  },
 
];

export default routes;
