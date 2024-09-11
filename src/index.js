import React, { useEffect } from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";
import routes from "./routes";
import { GoogleOAuthProvider } from "@react-oauth/google";

const root = ReactDOM.createRoot(document.getElementById("root"));

const Protected = (props) => {
  const navigate = useNavigate();
  const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";

  useEffect(() => {
    if (!isLoggedIn) {
      navigate("/");
    }
  }, [isLoggedIn, navigate]);

  return isLoggedIn ? props.com : null;
};

root.render(
  <GoogleOAuthProvider clientId="724153312920-dghbih3atf6c22o1nja1bnrdu1fpribo.apps.googleusercontent.com">
    <BrowserRouter>
      <Routes>
        {routes.map((route, i) => (
          <Route
            key={i}
            path={route.path}
            element={
              route.protected ? (
                <Protected com={route.element} />
              ) : (
                route.element
              )
            }
          />
        ))}
      </Routes>
    </BrowserRouter>
  </GoogleOAuthProvider>
);
