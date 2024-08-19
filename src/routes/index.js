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
    path: "/timetable/:departmentID/:semesterID",
    element: <Timetable  />,
  },
  {
    path: "/timetable/saved",
    element: <SaveTimetable />,
  },
  {
    path: "/timetable/saved/:departmentID/:semesterID",
    element: <SavedTimetable />,
  },
  {
    path: "/timetable/faculty",
    element: <FacTimetable />,
  },
  {
    path: "/timetable/faculty/:facultyName",
    element: <FacultyTimetable />,
  },
  {
    path: "/timetable/lab",
    element: <Lab />,
  },
  {
    path: "/timetable/lab/:subjectName",
    element: <LabTimetable />,
  },
 
];

export default routes;
