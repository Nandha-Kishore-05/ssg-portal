import React, { useState, useEffect } from 'react';
import axios from 'axios';

import AppLayout from '../../layout/layout';
import './workload.css';
import CustomSelect from '../../components/select';
import CustomButton from '../../components/button';
import LabTimetable from './labtable';

const Lab = () => {


  const [labOptions, setLabOptions] = useState([]);
  const [selectedLab, setSelectedLab] = useState("");
  const [isOpen,setIsOpen] = useState(false)
  useEffect(() => {
    axios.get('http://localhost:8080/timetable/labOptions')
      .then(response => {
        setLabOptions(response.data);
      })
      .catch(error => {
        console.error('Error fetching lab subject names:', error);
      });
  }, []);

  const handleViewTimetable = () => {
    if (selectedLab) {
      setIsOpen(true);
    } else {
      console.error('Please select a lab options');
    }
  };
    return (
        <AppLayout
          rId={5}
          title="Lab Table"
          body={
            <div style={{backgroundColor:"white",padding: 17,marginTop: 20,borderRadius:"10" }}>
                             <div style={{display:'flex',flexDirection:'row',columnGap:10,alignItems:"center"}}>
            <CustomSelect
            placeholder="Lab Name"
            value={selectedLab}
            onChange={setSelectedLab}
            options={labOptions}
          />
        
          
            <CustomButton
              width="150"
              label="View Timetable"
              onClick={handleViewTimetable}
            />
        
          </div>
          { (selectedLab && isOpen) && 
          <LabTimetable subjectName ={selectedLab.value} />
             
          }
                </div>
          }
          />
        );
};

export default Lab;