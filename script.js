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