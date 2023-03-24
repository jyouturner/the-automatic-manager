import "./App.css";
import { Link } from "react-router-dom";
import { Container, Navbar, Nav } from "react-bootstrap";
import { Outlet } from "react-router-dom";

function App() {
  return (
    <div className="app">
      <header>
        <Navbar bg="dark" variant="dark" expand="lg">
          <Container>
            <Navbar.Brand as={Link} to="/">
              TAMI
            </Navbar.Brand>
            <Navbar.Collapse>
              <Nav>
                <Nav.Link as={Link} to="/">
                  Home
                </Nav.Link>
                <Nav.Link as={Link} to="/automations">
                  Automations
                </Nav.Link>
                <Nav.Link as={Link} to="/about">
                  About
                </Nav.Link>
              </Nav>
            </Navbar.Collapse>
          </Container>
        </Navbar>
      </header>
      <Container className="mt-3">
        <Outlet />
      </Container>
    </div>
  );
}

export default App;
