import { Row } from "react-bootstrap";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createAutomation } from "../services/api";
import AutomationForm from "../components/AutomationForm";

export default function NewAutomation() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: "",
    source: "",
    destination: "",
    runtime: "",
    code: "",
  });

  const handleChange = (e) => {
    e.preventDefault();
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    await createAutomation(formData);
    navigate("/automations");
  };

  const handleCancel = (e) => {
    e.preventDefault();
    navigate("/automations");
  };

  return (
    <Row>
      <h2>New Automation</h2>
      <AutomationForm
        handleSubmit={handleSubmit}
        handleChange={handleChange}
        handleCancel={handleCancel}
        handleDelete={null}
        formData={formData}
        submitButtonText="Create"
        cancelButtonText="Cancel"
      />
    </Row>
  );
}
