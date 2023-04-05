import React from 'react';
import {
  ActionFunctionArgs, useNavigate, useFetcher,
} from 'react-router-dom';

import AssetTransferService from '../../services/AssetTransferService';
import { SLA } from '../../models/Asset.model';
import Button from '../../components/Button/Button';
import APIResponse from '../../models/APIResponse.model';
import Loader from '../../components/Loader/Loader';
import Input from '../../components/Input/Input';

export async function action({ request }: ActionFunctionArgs) {
  const formData = await request.formData();
  const slaObj = Object.fromEntries(formData);
  const actualSLA = JSON.parse(slaObj.sla as string) as SLA;

  return AssetTransferService.createAsset(actualSLA);
}

export function CreateAsset() {
  const navigate = useNavigate();
  const fetcher = useFetcher();

  const [message, setMessage] = React.useState('');

  React.useEffect(() => {
    if (fetcher.data && Object.keys(fetcher.data).length) {
      const data = fetcher.data as APIResponse;
      navigate('/navigation');
      if (!data.success) {
        setMessage(data.message!);
      }
    }
  }, [fetcher.data]);

  return (
    <>
      <div>
        <h1>Create Asset</h1>
        <fetcher.Form method="post">
          <Input type="textarea" id="sla" name="SLA:" required />
          <Button fullWidth>Submit</Button>
        </fetcher.Form>
      </div>
      <Button fullWidth onClick={() => navigate('/navigation')}>Back</Button>
      {
        fetcher.state === 'submitting'
        && <Loader />
      }
      {message && (
      <p>
        Error:
        {' '}
        {message}
      </p>
      )}
    </>
  );
}
