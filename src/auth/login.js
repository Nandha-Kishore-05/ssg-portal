import React from 'react';
import { GoogleOAuthProvider, GoogleLogin } from '@react-oauth/google';
import './style.css'; // Import your updated CSS file

function Login() {
  const handleLoginSuccess = (response) => {
    console.log('Login Success:', response);
    // Handle the response and perform necessary actions (e.g., saving token, redirecting)
  };

  const handleLoginError = (error) => {
    console.error('Login Error:', error);
    // Handle errors
  };

  return (
    <>
<div style={{marginTop:250,marginLeft:630}}>
      <div className="login-container">
        <h1>Welcome Back!</h1>
        <div className="login-form">
          <GoogleLogin
            onSuccess={handleLoginSuccess}
            onError={handleLoginError}
            logo_alignment="left"
          />
        </div>
        <p className="create-account">
          Don't have an account? <a href="#">Create one</a>
        </p>
      </div>
      </div>
      </>
  );
}

export default Login;
