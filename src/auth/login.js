

import React, { useState } from 'react';
import { GoogleLogin } from '@react-oauth/google';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './style.css';
import CustomButton from '../components/button';

const Login = () => {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [isGoogleSignedIn, setIsGoogleSignedIn] = useState(false); 

  const handleLoginSuccess = async (response) => {
    try {
      const token = response.credential;
      const userInfo = await fetchUserInfo(token);
      const { email } = userInfo;
      setEmail(email); 
  
      // Directly login using backend login API
      await handleBackendLogin(email);
    } catch (error) {
      setErrorMessage('Login Error: ' + error.message);
    }
  };

  const handleLoginError = (error) => {
    setErrorMessage('Login Error: ' + error.message);
  };

  const fetchUserInfo = async (token) => {
    const response = await fetch(`https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=${token}`);
    if (!response.ok) {
      throw new Error('Failed to fetch user info');
    }
    return response.json();
  };

  const handleBackendLogin = async (email) => {
    try {
      const loginResponse = await axios.post('http://localhost:8080/login', { email });
      if (loginResponse.data.auth_token) {
        handleSuccessfulLogin(loginResponse.data);
      } else {
        setErrorMessage('Login failed. Please try again.');
      }
    } catch (error) {
      setErrorMessage('Backend Login Error: ' + error.message);
    }
  };

  const handleSuccessfulLogin = (data) => {
    const { auth_token, user } = data;

    // Save auth token and user info to local storage
    localStorage.setItem('authToken', auth_token);
    localStorage.setItem('userId', user.id);
    localStorage.setItem('userName', user.name);
    localStorage.setItem('isLoggedIn', 'true');

    // Redirect to dashboard
    navigate('/dashboard');
  };

  return (
    <div className="auth-container">
      <div className="auth-form-container sign-in-container">
        <form className="auth-form">
          <h1>Login</h1>
          <div className="social-container">
            {!isGoogleSignedIn ? (
              <>
                <center>
                  <GoogleLogin
                    onSuccess={handleLoginSuccess}
                    onError={handleLoginError}
                    buttonText="Sign in with Google"
                  
                  />
                </center>
                {errorMessage && <p className="error-message">{errorMessage}</p>}
              </>
            ) : (
              <div className="loading-message">
                <p>Logging in...</p>
              </div>
            )}
          </div>

          {!isGoogleSignedIn && (
            <>
              <input type="email" placeholder="Email" className="auth-input" />
              <input type="password" placeholder="Password" className="auth-input" />
              <button className="auth-btn sign-in-btn" type="submit">Login</button>
            </>
          )}
        </form>
      </div>

      <div className="auth-overlay-container">
        <div className="auth-overlay">
          <div className="auth-overlay-panel auth-overlay-right">
            <h1>Hello, Friend!</h1>
            <p>Enter your personal details and start your journey with us</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
