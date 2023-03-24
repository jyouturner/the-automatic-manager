import axios from "axios";

const url = "http://localhost:9001/api/automations";

export async function getAutomations() {
  return await axios.get(`${url}`);
}

export async function getAutomation(id) {
  return await axios.get(`${url}/${id}`);
}

export async function createAutomation(automation) {
  return await axios.post(url, automation);
}

export async function updateAutomation(id, automation) {
  return await axios.patch(`${url}/${id}`, automation);
}

export async function deleteAutomation(id) {
  return await axios.delete(`${url}/${id}`);
}
