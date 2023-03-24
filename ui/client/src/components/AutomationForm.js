import { Form, Button, Row, Col, Modal } from "react-bootstrap";
import { useState, useEffect } from "react";
import CodeEditor from "@uiw/react-textarea-code-editor";
import "../App.css";

export default function AutomationForm({
  handleChange,
  handleSubmit,
  handleCancel,
  handleDelete,
  formData,
  submitButtonText,
  cancelButtonText,
}) {
  const [language, setLanguage] = useState("js");
  const [show, setShow] = useState(false);

  const handleClose = (deleting) => {
    setShow(false);
    if (deleting) handleDelete();
  };
  const handleShow = () => setShow(true);

  useEffect(() => {
    let newLanguage = selectLanguage(formData.runtime);
    setLanguage(newLanguage);
  }, [formData.runtime]);

  const handleLanguageChange = (e) => {
    e.preventDefault;
    let runtime = e.target.value;
    let newLanguage = selectLanguage(runtime);
    setLanguage(newLanguage);
    handleChange(e);
  };

  const selectLanguage = (runtime) => {
    let newLanguage;

    switch (runtime) {
      case "node":
        newLanguage = "js";
        break;
      case "python":
        newLanguage = "python";
        break;
      case "ruby":
        newLanguage = "ruby";
        break;
      default:
        newLanguage = "js";
    }

    return newLanguage;
  };

  return (
    <>
      <Form>
        <Form.Group className="mb-3">
          <Form.Label>Name</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter name"
            value={formData.name}
            onChange={handleChange}
            name="name"
          />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Source</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter source"
            value={formData.source}
            onChange={handleChange}
            name="source"
          />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Destination</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter destination"
            value={formData.destination}
            onChange={handleChange}
            name="destination"
          />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Runtime</Form.Label>
          <Form.Select
            value={formData.runtime}
            onChange={handleLanguageChange}
            name="runtime"
          >
            <option>Select runtime</option>
            <option value="node">Node</option>
            <option value="ruby">Ruby</option>
            <option value="python">Python</option>
          </Form.Select>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Code</Form.Label>
          <CodeEditor
            value={formData.code}
            language={language}
            name="code"
            placeholder="Please enter some code..."
            onChange={handleChange}
            padding={15}
            style={{
              fontSize: 12,
              height: "20em",
              backgroundColor: "#333",
              fontFamily:
                "ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace",
            }}
          />
        </Form.Group>
        <Row>
          <Col md={"auto"}>
            <Button variant="primary" type="submit" onClick={handleSubmit}>
              {submitButtonText}
            </Button>
          </Col>
          <Col md={"auto"}>
            <Button variant="secondary" type="button" onClick={handleCancel}>
              {cancelButtonText}
            </Button>
          </Col>
          {handleDelete && (
            <Col md={"auto"}>
              <Button
                variant="danger"
                type="button"
                onClick={() => handleShow()}
              >
                Delete
              </Button>
            </Col>
          )}
        </Row>
      </Form>

      <Modal show={show} onHide={() => handleClose(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Delete Automation</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          Are you sure you want to delete {formData.name}?
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => handleClose(false)}>
            No, take me back!
          </Button>
          <Button variant="danger" onClick={() => handleClose(true)}>
            DELETE
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
}
