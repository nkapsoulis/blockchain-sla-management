import { getLoginURL, getLogoutURL, getUserURL } from '../utils/urls';

function login(name: string, mnemonic: string) {
  return fetch(getLoginURL(), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      name,
      mnemonic,
    }),
  });
}

function logout() {
  return fetch(getLogoutURL(), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
  });
}

async function getUser() {
  const response = await fetch(getUserURL(), {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
  });

  if (response.status !== 200) {
    return { success: false, message: 'Unauthorized' };
  }

  return response.json();
}

const authService = {
  login,
  logout,
  getUser,
};

export default authService;
