import React from 'react';
import {
  ActionFunction, Form, useActionData, useNavigate,
} from 'react-router-dom';
import Button from '../../components/Button/Button';
import Input from '../../components/Input/Input';
import APIResponse from '../../models/APIResponse.model';
import authService from '../../services/AuthService';
import { setBaseUrl } from '../../utils/urls';
import crypto from '../../utils/crypto';
import localStorage from '../../utils/localStorage';

export const action: ActionFunction = async ({ request }): Promise<APIResponse> => {
  const formData = await request.formData();
  const name = formData.get('username') as string;
  const mnemonic = formData.get('mnemonic') as string;
  const passphrase = formData.get('passphrase') as string;
  const serverAddress = formData.get('server-address') as string;
  const err = setBaseUrl(serverAddress);
  if (err !== '') {
    return { success: false, message: err };
  }

  const res = await authService.login(name, mnemonic);
  if (res.status === 200) {
    await localStorage.setLocalStorage('mnemonic', crypto.encrypt(mnemonic, passphrase));
    return { success: true };
  }

  const json = await res.json();
  return { success: false, message: json.message };
};

export function Authentication() {
  const navigate = useNavigate();
  const [error, setError] = React.useState('');
  const data = useActionData() as APIResponse;

  React.useEffect(() => {
    if (data) {
      if (data.success) {
        navigate('/navigation');
      } else {
        setError(data.message!);
      }
    }
  }, [data]);

  return (
    <>
      <h1>Log In</h1>
      <Form method="post">
        <Input type="text" id="username" name="Username:" required />
        <Input type="text" id="mnemonic" name="Mnemonic:" required />
        <Input type="password" id="passphrase" name="Passphrase:" required />
        <Input type="url" id="server-address" name="Server address:" required />
        <Button fullWidth>Log In</Button>
      </Form>
      {error !== ''
      && <p>{`Error: ${error}`}</p>}
    </>
  );
}
