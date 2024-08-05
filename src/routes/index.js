import Dashboard from "../pages/dashboard/dashboard";
import Timetable from "../pages/workload/workload";

import SavedTimetable from "../pages/workload/timetable";



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
 
];

export default routes;
