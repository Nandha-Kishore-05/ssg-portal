import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';

const GenerateTimetable = () => {
  const navigate = useNavigate();
  const [department, setDepartment] = useState(null);
  const [semester, setSemester] = useState(null);

  const handleViewTimetable = () => {
    if (department && semester) {
      navigate(`/timetable/${department.value}/${semester.value}`);
    } else {
      console.error('Please select both department and semester');
    }
  };
    return (
        <AppLayout
          rId={2}
          title="Time Table"
          body={
            <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
            <CustomSelect
            placeholder="DEPARTMENT"
            value={department}
            onChange={setDepartment}
            options={[
              { label: "COMPUTER TECHNOLOGY", value: 1 },
              { label: "BIO TECHNOLOGY", value: 2 },
            ]}
          />
          <CustomSelect
            placeholder="SEMESTER"
            value={semester}
            onChange={setSemester}
            options={[
              { label: "S1", value: 1 },
              { label: "S3", value: 3 },
              { label: "S5", value: 5 },
            ]}
          />
          <br />
          <center>
            <CustomButton
              width="150"
              label="Generate Timetable"
              onClick={handleViewTimetable}
            />
          </center>
                </div>
          }
          />
        );
};

export default GenerateTimetable;