import { useState, useEffect } from "react";
import { Button } from "react-bootstrap";
import { Link, NavLink } from "react-router-dom";
import { getAutomations } from "../services/api";
import "../App.css";

export default function AutomationsIndex() {
  const [automations, setAutomations] = useState([]);
  useEffect(() => {
    const fetchAutomations = async () => {
      const response = await getAutomations();
      setAutomations(response.data.automations);
    };
    fetchAutomations();
  }, []);

  return (
    <div>
      <div className="padded-element">
        <h2>Automations</h2>
        <Link to="/automations/new">
          <Button variant="primary">Create New Automation</Button>
        </Link>
      </div>
      <div className="padded-element">
        <h5>Existing Automations</h5>
        {automations.map((automation) => (
          <div key={automation.id}>
            <NavLink to={`/automations/${automation.id}`}>
              {automation.name}
            </NavLink>
          </div>
        ))}
      </div>
    </div>
  );
}
