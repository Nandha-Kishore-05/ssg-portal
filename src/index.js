import React, { useEffect } from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";
import routes from "./routes";
import { GoogleOAuthProvider } from "@react-oauth/google";


const root = ReactDOM.createRoot(document.getElementById("root"));
// const Protected = (props) => {
//   const navigate = useNavigate();
//   const isLogin = localStorage.getItem("DATA#1");

//   useEffect(() => {
//     if (isLogin !== "true") {
//       navigate("/auth/login");
//     }
//   }, [isLogin, navigate]);

//   return isLogin === "true" ? props.com : null;
// };
root.render(
  <GoogleOAuthProvider clientId="760746068225-seabd5fi45talacld9hpa5ovpr2elt51.apps.googleusercontent.com">
  <BrowserRouter>
    <Routes>
      {routes.map((route) => (
        <Route key={route} path={route.path} element={route.element} />
      ))}
    </Routes>
  </BrowserRouter>
  </GoogleOAuthProvider>
);
