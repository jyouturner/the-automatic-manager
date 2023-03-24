import { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import { Container, Button, Row, Col } from "react-bootstrap";
import { getAutomation } from "../services/api";

export default function Automation() {
  let params = useParams();

  const [automation, setAutomation] = useState([]);
  useEffect(() => {
    const fetchAutomation = async () => {
      const response = await getAutomation(params.automationId);
      setAutomation(response.data.automation);
    };
    fetchAutomation();
  }, [params.automationId]);

  return (
    <Container>
      <h2>{automation.name}</h2>
      <p>
        {automation.source} =&gt; {automation.destination}
      </p>
      <div>
        <Row>
          <Col md={"auto"}>
            <Link to={`/automations/${automation.id}/edit`}>
              <Button variant="primary">Edit</Button>
            </Link>
          </Col>
          <Col md={"auto"}>
            <Link to={`/automations`}>
              <Button variant="secondary">Back</Button>
            </Link>
          </Col>
        </Row>
      </div>
    </Container>
  );
}
