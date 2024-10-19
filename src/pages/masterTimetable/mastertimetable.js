
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import AppLayout from '../../layout/layout';
import CustomCard from '../../components/card';
import AcademicYearModal from '../../components/modal';

import './mastertimetable.css';

function Mastertimetable() {
  const navigate = useNavigate();
  const [openModal, setOpenModal] = useState(false);
  const [selectedYear, setSelectedYear] = useState(null);

  const [cards, setCards] = useState([]);
  const [academicYears, setAcademicYears] = useState([]);

  useEffect(() => {
    axios.get('http://localhost:8080/acdemicYearOptions')
      .then(response => {
        setAcademicYears(response.data);
      })
      .catch(error => {
        console.error('Error fetching academic year options:', error);
      });
  }, []);

  const handleOpenModal = () => {
    setOpenModal(true);
  };

  const handleCloseModal = () => {
    setOpenModal(false);

  };

  const handleYearChange = (event) => {
    const selected = academicYears.find(year => year.value === event.target.value);
    setSelectedYear(selected ? selected : null);
  };


  const handleAddCard = () => {
    if (selectedYear) {
      setCards([...cards, { year: selectedYear }]);
      setSelectedYear(null);
    
      handleCloseModal();
    }
  };

  const handleDownload = (yearId, type) => {
    const downloadUrl = `http://localhost:8080/download/${yearId}`;
    axios.get(downloadUrl, { responseType: 'blob' })
      .then((response) => {
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `timetable_${yearId}_${type}.xlsx`);
        document.body.appendChild(link);
        link.click();
        link.remove();
      })
      .catch((error) => {
        console.error('Error downloading timetable:', error);
      });
  };



  return (
    <AppLayout
      rId={11}
      title="Master Timetable"
      body={
        <>
          <div className="cards-container">
            {cards.map((cardObj, index) => (
              <CustomCard 
                key={index} 
                year={cardObj.year.label}
                semesterType={cardObj.semesterType}
                title={`Academic Year`}
            
                onDownload={() => handleDownload(cardObj.year.value, cardObj.semesterType)} 
              />
            ))}
            <CustomCard onAddCard={handleOpenModal} />
          </div>

          <AcademicYearModal
            open={openModal}
            onClose={handleCloseModal}
            selectedYear={selectedYear}
            onYearChange={handleYearChange}
        
            academicYears={academicYears}
            onAddCard={handleAddCard}
            title="Select an Academic Year and Semester"
            placeholder="Choose an academic year"
          />
        </>
      }
    />
  );
}

export default Mastertimetable;

