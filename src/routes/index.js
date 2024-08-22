import Dashboard from "../pages/dashboard/dashboard";

import Login from "../auth/login";

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
