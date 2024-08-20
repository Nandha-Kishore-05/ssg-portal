import Dashboard from "../pages/dashboard/dashboard";
import Timetable from "../pages/workload/workload";

import SavedTimetable from "../pages/workload/timetable";
import FacultyTimetable from "../pages/workload/facultytable";
import Login from "../auth/login";
import LabTimetable from "../pages/workload/labtable";
import SaveTimetable from "../pages/workload/save";
import GenerateTimetable from "../pages/workload/generate";
import FacTimetable from "../pages/workload/faculty";
import Lab from "../pages/workload/lab";
import PeriodAllocation from "../pages/allocation/periodAllocation";



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
    path: "/timetable",
    element: <GenerateTimetable  />,
  },
 
  {
    path: "/timetable/saved",
    element: <SaveTimetable />,
  },
 
  {
    path: "/timetable/faculty",
    element: <FacTimetable />,
  },
  
  {
    path: "/timetable/lab",
    element: <Lab />,
  },
  {
    path: "/timetable/periodallocation",
    element: <PeriodAllocation />,
  },

 
];

export default routes;
