// Modal.js
import React from 'react';
import Modal from 'react-modal';
import './style.css'; // Ensure you have the CSS for styling the modal

Modal.setAppElement('#root'); // This is important for accessibility

const CustomModal = ({ isOpen, onClose, onConfirm, message }) => {
  return (
    <Modal
      isOpen={isOpen}
      onRequestClose={onClose}
      contentLabel="Confirmation Modal"
      className="modal-content"
      overlayClassName="modal-overlay"
    >
      <h2>{message}</h2>
      <div className="modal-actions">
        <button onClick={onConfirm}>Yes</button>
        <button onClick={onClose}>No</button>
      </div>
    </Modal>
  );
};

export default CustomModal;
