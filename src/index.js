import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import routes from "./routes";


const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <BrowserRouter>
    <Routes>
      {routes.map((route) => (
        <Route key={route} path={route.path} element={route.element} />
      ))}
    </Routes>
  </BrowserRouter>
);
