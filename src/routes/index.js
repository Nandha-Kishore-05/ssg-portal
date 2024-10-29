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
import Mastertimetable from "../pages/masterTimetable/mastertimetable";
import Masterdepartment from "../pages/masterTimetable/masterdepartment";
import MasterSemester from "../pages/masterTimetable/masterSemester";
import Venue from "../pages/workload/venue/table";
import StudentEntry from "../pages/student/studentEntry";
import StudentTable from "../pages/student/studentTable";
import StudentTimetable from "../pages/student/studentTimetable";
import LabEntry from "../pages/labentry/labentry";





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
    path: "/mastertimetable",
    element: <Mastertimetable />,
  },
  {
    path: "/masterdepartment",
    element: <Masterdepartment />,
  },
  {
    path: "/mastersemester",
    element: <MasterSemester />,
  },
  {
    path: "/venueTable",
    element: <Venue />,
  },
  {
    path: "/studentallocation",
    element: <StudentEntry />,
  },
  {
    path: "/studentTable",
    element: <StudentTable />,
  },
  {
    path: "/labentry",
    element: <LabEntry />,
  },
  {
    path: "/logout",
    element: <Logout />,
  },

];

export default routes;
