import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';

const Lab = () => {
  const navigate = useNavigate();
  const [department, setDepartment] = useState(null);
 

  const handleViewTimetable = () => {
    if (department ) {
      navigate(`/timetable/lab/${department.value}`);
    } else {
      console.error('Please select both department and semester');
    }
  };
    return (
        <AppLayout
          rId={5}
          title="Venue Table"
          body={
            <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
            <CustomSelect
            placeholder="DEPARTMENT"
            value={department}
            onChange={setDepartment}
            options={[
              { label: "dce", value: "dce" },
              { label: "BIO TECHNOLOGY", value: 2 },
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

export default Lab;