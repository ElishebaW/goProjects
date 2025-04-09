"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
const API_BASE_URL = "http://localhost:8080"; // Replace with your backend URL
// Fetch tasks from the backend
function fetchTasks() {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch(`${API_BASE_URL}/tasks`);
        if (!response.ok) {
            throw new Error("Failed to fetch tasks");
        }
        return response.json();
    });
}
// Add a new task
function addTask(name, description) {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch(`${API_BASE_URL}/tasks`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name, description }),
        });
        if (!response.ok) {
            throw new Error("Failed to add task");
        }
    });
}
// Reorganize tasks
function reorganizeTasks() {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch(`${API_BASE_URL}/tasks/reorganize`, {
            method: "POST",
        });
        if (!response.ok) {
            throw new Error("Failed to reorganize tasks");
        }
        return response.json();
    });
}
// Render tasks in the UI
function renderTasks(tasks) {
    const taskList = document.getElementById("taskList");
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
function init() {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const tasks = yield fetchTasks();
            renderTasks(tasks);
        }
        catch (error) {
            console.error(error);
            alert("Failed to load tasks");
        }
    });
}
// Handle form submission
document.getElementById("taskForm").addEventListener("submit", (event) => __awaiter(void 0, void 0, void 0, function* () {
    event.preventDefault();
    const taskName = document.getElementById("taskName").value;
    const taskDescription = document.getElementById("taskDescription").value;
    try {
        yield addTask(taskName, taskDescription);
        const tasks = yield fetchTasks();
        renderTasks(tasks);
        document.getElementById("taskForm").reset();
    }
    catch (error) {
        console.error(error);
        alert("Failed to add task");
    }
}));
// Handle reorganize button click
document.getElementById("reorganizeTasks").addEventListener("click", () => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const tasks = yield reorganizeTasks();
        renderTasks(tasks);
    }
    catch (error) {
        console.error(error);
        alert("Failed to reorganize tasks");
    }
}));
// Load tasks on page load
init();
