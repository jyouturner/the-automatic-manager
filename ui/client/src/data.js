let automations = [
  {
    id: 1,
    name: "Automation 1",
    source: "Jira",
    destination: "Notion",
  },
  {
    id: 2,
    name: "Automation 2",
    source: "Confluence",
    destination: "Notion",
  },
  {
    id: 3,
    name: "Automation 3",
    source: "Notion",
    destination: "Slack",
  },
];

export function getAutomations() {
  return automations;
}

export function getAutomation(id) {
  return automations.find((automation) => automation.id === id);
}
