interface Task {
  id: number;
  name: string;
  order: number;
  description: string;
}

const API_BASE_URL = "http://localhost:8080"; // Replace with your backend URL

// Fetch tasks from the backend
async function fetchTasks(): Promise<Task[]> {
  const response = await fetch(`${API_BASE_URL}/tasks`);
  if (!response.ok) {
    throw new Error("Failed to fetch tasks");
  }
  return response.json();
}

// Add a new task
async function addTask(name: string, description: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/tasks`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, description }),
  });
  if (!response.ok) {
    throw new Error("Failed to add task");
  }
}

// Reorganize tasks
async function reorganizeTasks(): Promise<Task[]> {
  const response = await fetch(`${API_BASE_URL}/tasks/reorganize`, {
    method: "POST",
  });
  if (!response.ok) {
    throw new Error("Failed to reorganize tasks");
  }
  return response.json();
}

// Render tasks in the UI
function renderTasks(tasks: Task[]): void {
  const taskList = document.getElementById("taskList")!;
  taskList.innerHTML = "";

  tasks
    .sort((a, b) => a.order - b.order)
    .forEach((task) => {
      const listItem = document.createElement("li");
      listItem.className = "list-group-item d-flex justify-content-between align-items-center";
      listItem.innerHTML = `
        <div>
          <strong>${task.name}</strong>
          <p class="mb-0">${task.description}</p>
        </div>
        <span class="badge bg-primary rounded-pill">Order: ${task.order}</span>
      `;
      taskList.appendChild(listItem);
    });
}

// Initialize the app
async function init() {
  try {
    const tasks = await fetchTasks();
    renderTasks(tasks);
  } catch (error) {
    console.error(error);
    alert("Failed to load tasks");
  }
}

// Handle form submission
document.getElementById("taskForm")!.addEventListener("submit", async (event) => {
  event.preventDefault();

  const taskName = (document.getElementById("taskName") as HTMLInputElement).value;
  const taskDescription = (document.getElementById("taskDescription") as HTMLTextAreaElement).value;

  try {
    await addTask(taskName, taskDescription);
    const tasks = await fetchTasks();
    renderTasks(tasks);
    (document.getElementById("taskForm") as HTMLFormElement).reset();
  } catch (error) {
    console.error(error);
    alert("Failed to add task");
  }
});

// Handle reorganize button click
document.getElementById("reorganizeTasks")!.addEventListener("click", async () => {
  try {
    const tasks = await reorganizeTasks();
    renderTasks(tasks);
  } catch (error) {
    console.error(error);
    alert("Failed to reorganize tasks");
  }
});

// Load tasks on page load
init();