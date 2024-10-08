
import React from 'react';
import { Modal, Box, Typography, Select, MenuItem, Button } from '@mui/material';
import PropTypes from 'prop-types';

const style = {
  position: 'absolute',
  top: '45%',
  left: '55%',
  transform: 'translate(-50%, -50%)',
  width: 550, 
  bgcolor: 'background.paper',
  boxShadow: 24,
  borderRadius: '12px',
  p: 5, 
};

const AcademicYearModal = ({
  open,
  onClose,
  selectedYear,
  onYearChange,
  selectedSemesterType,
  onSemesterTypeChange,
  academicYears,
  onAddCard,
  title,
  placeholder,
}) => {
  return (
    <Modal
      open={open}
      onClose={onClose}
      aria-labelledby="modal-title"
      aria-describedby="modal-description"
    >
      <Box sx={style}>
 
        <Typography 
          id="modal-title" 
          variant="h5" 
          component="h1" 
          gutterBottom 
          sx={{ textAlign: 'center', fontWeight: 'bold' }}
        >
          {title}
        </Typography>

       
        <Select
          value={selectedYear ? selectedYear.value : ''}
          onChange={onYearChange}
          displayEmpty
          fullWidth
          sx={{ mt: 2, mb: 3 }}
        >
          <MenuItem value="" disabled>
            {placeholder}
          </MenuItem>
          {academicYears.map((year) => (
            <MenuItem key={year.value} value={year.value}>
              {year.label}
            </MenuItem>
          ))}
        </Select>

    
        <Select
          value={selectedSemesterType}
          onChange={onSemesterTypeChange}
          displayEmpty
          fullWidth
          sx={{ mb: 4 }}
        >
          <MenuItem value="" disabled>
            Select Semester Type
          </MenuItem>
          <MenuItem value="odd">ODD</MenuItem>
          <MenuItem value="even">EVEN</MenuItem>
        </Select>

        <Button
          variant="contained"
          fullWidth
          onClick={onAddCard}
          disabled={!selectedYear || !selectedSemesterType}
          sx={{
            padding: '10px 0', 
            backgroundColor: '#1976d2', 
            '&:hover': {
              backgroundColor: '#1565c0',
            },
            fontSize: '16px', 
            fontWeight: 'bold',
          }}
        >
          Submit
        </Button>
      </Box>
    </Modal>
  );
};

AcademicYearModal.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  selectedYear: PropTypes.object,
  onYearChange: PropTypes.func.isRequired,
  selectedSemesterType: PropTypes.string.isRequired,
  onSemesterTypeChange: PropTypes.func.isRequired,
  academicYears: PropTypes.arrayOf(PropTypes.object).isRequired,
  onAddCard: PropTypes.func.isRequired,
  title: PropTypes.string.isRequired,
  placeholder: PropTypes.string.isRequired,
};

export default AcademicYearModal;
