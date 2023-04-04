import React from 'react';
import { useNavigate } from 'react-router-dom';

import Button from '../../components/Button/Button';
import Loader from '../../components/Loader/Loader';
import RootAPIService from '../../services/RootAPIService';
import cookies from '../../utils/cookies';
import localStorage from '../../utils/localStorage';
import authService from '../../services/AuthService';

function Navigation() {
  const navigate = useNavigate();
  const [error, setError] = React.useState('');
  const [serverName, setServerName] = React.useState('');
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    RootAPIService.getName()
      .then((name) => {
        setServerName(name);
        setLoading(false);
      })
      .catch((e: unknown) => {
        setError(e as string);
      });
  }, []);

  const onLogOutClick = async () => {
    const res = await authService.logout();
    if (res.status !== 200) {
      const json = await res.json();
      setError(error === '' ? error : `${error}\n${json.message}`);
    }

    Promise.all([cookies.deleteCookie(), localStorage.deleteLocalStorage('url')])
      .then(([deleted, removed]) => {
        if (deleted !== '') {
          setError(deleted);
          return;
        }
        if (!removed) {
          setError(error === '' ? error : `${error}\nLocal storage error`);
          return;
        }
        navigate('/');
      });
  };

  return (
    <div>
      {loading
        ? <Loader />
        : (
          <>
            <h2>
              Welcome at
              {' '}
              {serverName}
            </h2>
            <Button fullWidth onClick={() => navigate('/assets')}>My Assets</Button>
            <Button fullWidth onClick={() => navigate('/get-asset')}>Get Asset</Button>
            <Button fullWidth onClick={() => navigate('/transfer')}>Transfer Asset</Button>
            <Button fullWidth onClick={() => navigate('/create-asset')}>Create Asset</Button>
          </>
        )}
      <Button fullWidth onClick={onLogOutClick}>Logout</Button>
      { error ? <p>{error}</p> : null}
    </div>
  );
}

export default Navigation;
