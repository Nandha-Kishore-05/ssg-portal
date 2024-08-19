import Dashboard from "../pages/dashboard/dashboard";
import Timetable from "../pages/timetable/workload";

import SavedTimetable from "../pages/timetable/timetable";
import FacultyTimetable from "../pages/timetable/facultytable";
import Login from "../auth/login";
import LabTimetable from "../pages/timetable/labtable";
import SaveTimetable from "../pages/timetable/save";
import GenerateTimetable from "../pages/timetable/generate";
import FacTimetable from "../pages/timetable/faculty";
import Lab from "../pages/timetable/lab";
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
  {
    path: "/timetable/periodallocation",
    element: <PeriodAllocation />,
  },
 
];

export default routes;
