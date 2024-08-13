import Dashboard from "../pages/dashboard/dashboard";
import Timetable from "../pages/workload/workload";

import SavedTimetable from "../pages/workload/timetable";
import FacultyTimetable from "../pages/workload/facultytable";
import Login from "../auth/login";



const routes = [
  {
    path: "/",
    element: <Login />,
  },
  {
    path: "/dashboard",
    element: <Dashboard />,
  },
  {
    path: "/timetable/:departmentID/:semesterID",
    element: <Timetable />,
  },
  {
    path: "/timetable/saved/:departmentID/:semesterID",
    element: <SavedTimetable />,
  },
  {
    path: "/timetable/faculty/:facultyName",
    element: <FacultyTimetable />,
  },
 
];

export default routes;
