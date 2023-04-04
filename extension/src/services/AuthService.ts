import { getLoginURL, getLogoutURL } from '../utils/urls';

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

const authService = {
  login,
  logout,
};

export default authService;
