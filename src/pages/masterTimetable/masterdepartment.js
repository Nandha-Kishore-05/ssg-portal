import React, { useState, useEffect } from 'react'; 
import { useNavigate, useLocation } from 'react-router-dom'; 
import axios from 'axios'; 
import AppLayout from '../../layout/layout';
import CustomCard from '../../components/card';
import './mastertimetable.css';

function Masterdepartment() {
  const navigate = useNavigate(); 
  const location = useLocation(); 
  const { yearId } = location.state || {}; 

  const [cards, setCards] = useState([]); 
  const [deptOptions, setDeptOptions] = useState([]); 


  useEffect(() => {
    axios.get('http://localhost:8080/timetable/options')
      .then(response => {
        setDeptOptions(response.data); 
        setCards(response.data);
      })
      .catch(error => {
        console.error('Error fetching department options:', error);
      });
  }, []);

  const handleCardClick = (semester) => {
    console.log('Clicked Department:', semester);
    navigate('/mastersemester', {
      state: {
        yearId, 
        deptId: semester.value 
      }
    });
  };

  return (
    <AppLayout
      rId={11}
      title={`Department`} 
      body={
        <div className="cards-container">
          {cards.map((semesterObj, index) => (
            <CustomCard 
              key={index} 
              year={semesterObj.label} 
              title={`Semester`} 
              onCardClick={() => handleCardClick(semesterObj)} 
            />
          ))}
        </div>
      }
    />
  );
}

export default Masterdepartment;
