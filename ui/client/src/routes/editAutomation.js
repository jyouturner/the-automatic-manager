import { Row } from "react-bootstrap";
import { useParams, useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import {
  deleteAutomation,
  getAutomation,
  updateAutomation,
} from "../services/api";
import AutomationForm from "../components/AutomationForm";

export default function EditAutomation() {
  const params = useParams();
  const navigate = useNavigate();

  const [automation, setAutomation] = useState([]);
  useEffect(() => {
    const fetchAutomation = async () => {
      const response = await getAutomation(params.automationId);
      setFormData({ ...response.data.automation });
      setAutomation(response.data.automation);
    };
    fetchAutomation();
  }, []);

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
    await updateAutomation(automation.id, formData);
    navigate("/automations/" + automation.id);
  };

  const handleCancel = (e) => {
    e.preventDefault();
    navigate("/automations");
  };

  const handleDelete = async () => {
    await deleteAutomation(automation.id);
    navigate("/automations");
  };

  return (
    <Row>
      <h2>Edit Automation {automation.name}</h2>
      <AutomationForm
        handleSubmit={handleSubmit}
        handleChange={handleChange}
        handleCancel={handleCancel}
        handleDelete={handleDelete}
        formData={formData}
        submitButtonText="Update"
        cancelButtonText="Cancel"
      />
    </Row>
  );
}
