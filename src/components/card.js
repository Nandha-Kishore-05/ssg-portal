
import React from "react";
import DownloadIcon from '@mui/icons-material/Download';
import "./style.css";

function CustomCard({ onAddCard, year, semesterType, onCardClick, onDownload }) {
  return (
    <div className="custom-card" onClick={onCardClick}>
      {year && (
        <DownloadIcon className="download-icon" onClick={onDownload} style={{fontSize:'33px',marginRight:'4px'}} />
      )}
      {year ? (
        <div className="dotted-square">
          
          <span className="year-text">{year} {semesterType}</span>
        </div>
      ) : (
        <div className="dotted-square" onClick={onAddCard}>
          <span className="plus-icon">+</span>
        </div>
      )}
    </div>
  );
}

export default CustomCard;
