import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';

const FacTimetable = () => {
  const navigate = useNavigate();
  const [department, setDepartment] = useState("");
 

  const handleViewTimetable = () => {
    if (department ) {
      navigate(`/timetable/faculty/${department.value}`);
    } else {
      console.error('Please select both department and semester');
    }
  };
    return (
        <AppLayout
          rId={4}
          title="Faculty Table"
          body={
            <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
            <CustomSelect
            placeholder="Faculty Name"
            value={department}
            onChange={setDepartment}
            options={[
              { label: "Dr. Alice Smith", value: "Dr. Alice Smith" },
              { label: "Dr. David Brown", value: "Dr. David Brown" },
            ]}
          />
         
          <br />
          <center>
            <CustomButton
              width="150"
              label="View Timetable"
              onClick={handleViewTimetable}
            />
          </center>
                </div>
          }
          />
        );
};

export default FacTimetable;