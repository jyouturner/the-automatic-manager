import "bootstrap/dist/css/bootstrap.css";
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import "./index.css";
import App from "./App";
import Automations from "./routes/automations";
import AutomationsIndex from "./routes/automationsIndex";
import Automation from "./routes/automation";
import NewAutomation from "./routes/newAutomation";
import EditAutomation from "./routes/editAutomation";
import About from "./routes/about";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />}>
          <Route path="/automations" element={<Automations />}>
            <Route path="" element={<AutomationsIndex />} />
            <Route path="new" element={<NewAutomation />} />
            <Route path=":automationId" element={<Automation />} />
            <Route path=":automationId/edit" element={<EditAutomation />} />
          </Route>
          <Route path="/about" element={<About />} />
          <Route
            path="*"
            element={
              <main style={{ padding: "1em" }}>
                <p>There&apos;s nothing here!</p>
              </main>
            }
          />
        </Route>
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
