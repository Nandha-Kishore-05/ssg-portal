import Dashboard from "../pages/dashboard/dashboard";
import Login from "../auth/login";
import SaveTimetable from "../pages/workload/save";
import GenerateTimetable from "../pages/workload/generate";
import FacTimetable from "../pages/workload/faculty";
import Lab from "../pages/workload/lab";
import ManualEntry from "../entry/manualEntry";
import SubjectEntry from "../pages/entry/subjectentry";
import SubjectAllocation from "../pages/allocation/allocation";
import Logout from "../auth/logout";




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
    path: "/timetable/subjectentry",
    element: <SubjectEntry />,
  },

  {
    path: "/manualentry",
    element: <ManualEntry />,
  },
  {
    path: "/subjectallocation",
    element: <SubjectAllocation />,
  },
  {
    path: "/logout",
    element: <Logout />,
  },

];

export default routes;
