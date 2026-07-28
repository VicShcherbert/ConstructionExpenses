export const fetchProjects = async (userId) => {
  const url = 'http://localhost:8080/get-projects/' + userId;
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Response status: ${response.status}`);
  }
  const result = await response.json();
  return result ?? [];
};

export const fetchProjectInfo = async (project_id) => {
  const projectURL = 'http://localhost:8080/get-project/' + project_id;
  const response = await fetch(projectURL);
  if (!response.ok) {
    throw new Error(`Response status: ${response.status}`);
  }
  const result = await response.json();
  return result;
};

export const fetchExpenses = async (project_id) => {
  const expensesURL =
    'http://localhost:8080/get-project-expenses/' + project_id;
  const response = await fetch(expensesURL);
  if (!response.ok) {
    throw new Error(`Response status: ${response}`);
  }
  const result = await response.json();
  return result;
};
