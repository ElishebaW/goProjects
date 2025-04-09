const getUsersBtn = document.getElementById('get-users-btn');
const addUserBtn = document.getElementById('add-user-btn');

getUsersBtn.addEventListener('click', async () => {
  const response = await fetch('/api/users');
  const data = await response.json();
  console.log(data);
});

addUserBtn.addEventListener('click', async () => {
  const name = prompt('Enter user name:');
  const email = prompt('Enter user email:');
  const response = await fetch('/api/user', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email }),
  });
  console.log(response.status);
});

const tasksElement = document.getElementById('tasks');
const addTaskBtn = document.getElementById('add-task-btn');
const startPomodoroBtn = document.getElementById('start-pomodoro-btn');

let timerId;

// Get tasks from server and display in list
fetch('/tasks')
  .then(response => response.json())
  .then(tasks => {
    const taskListHtml = tasks.map(task => `<li>${task.name} (${task.order})</li>`).join('');
    tasksElement.innerHTML = `
      <ul>
        ${taskListHtml}
      </ul>
    `;
  });

// Add new task when button is clicked
addTaskBtn.addEventListener('click', () => {
  const name = prompt('Enter task name:');
  if (name) {
    fetch('/task', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, order: Math.random() * 10 }),
    })
      .then(response => response.json())
      .then(task => {
        const taskListHtml = tasksElement.innerHTML + `<li>${task.name} (${task.order})</li>`;
        tasksElement.innerHTML = taskListHtml;
      });
  }
});

// Start Pomodoro timer when button is clicked
startPomodoroBtn.addEventListener('click', () => {
  timerId = setInterval(() => {
    const minutes = Math.floor((25 - (Date.now() % 1000)) / 60000);
    const seconds = Math.floor((25 - (Date.now() % 1000)) / 1000) % 60;
    document.getElementById('timer').innerText = `${minutes}:${seconds}`;
  }, 10);
});

// Stop Pomodoro timer when context is done
window.addEventListener('contextmenu', event => {
  clearInterval(timerId);
});