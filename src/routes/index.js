import Dashboard from "../pages/dashboard/dashboard";
import Workload from "../pages/workload/workload";




const routes = [
  
  {
    path: "/dashboard",
    element: <Dashboard />,
  },
  {
    path: "/timetable/:departmentID",
    element: <Workload />,
  },
 
];

export default routes;
