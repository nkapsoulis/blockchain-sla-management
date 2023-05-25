import React from 'react';
import { useLoaderData, useNavigate } from 'react-router-dom';
import { JsonViewer } from '@textea/json-viewer';

import Button from '../../components/Button/Button';
import { SLA } from '../../models/Asset.model';
import User from '../../models/User.model';

export default function MyAssets() {
  const [error, setError] = React.useState('');
  const [assets, setAssets] = React.useState<SLA[]>([]);
  const [userData, setUserData] = React.useState<User>({} as User);

  const data = useLoaderData() as any;
  React.useEffect(() => {
    console.log(data);
    if (data[0].success !== undefined) {
      setError(data.message);
    } else if (data[1].success !== undefined) {
      setError(`${error} ${data[1].message}`);
    } else {
      setAssets(data[0].assets ? data[0].assets : []);
      setUserData(data[1].user);
      setError('');
    }
  }, [data]);

  const navigate = useNavigate();
  return (
    <>
      {
        error !== ''
        && <p>{`Error: ${error}`}</p>
      }
      {userData && (
        <>
          <h1>{userData.name}</h1>
          <h2>
            Balance:
            {' '}
            {userData.balance}
          </h2>
        </>
      )}
      <h3>Assets</h3>
      {assets.length
        ? (
          <table>
            <tbody>
              {assets.map((asset) => (
                <tr key={asset.id}>
                  <JsonViewer value={asset} rootName={asset.id} defaultInspectDepth={0} />
                </tr>
              ))}
            </tbody>
          </table>
        )
        : error === ''
        && (<p>There are no assets to show.</p>)}
      <Button fullWidth onClick={() => navigate('/navigation')}>Back</Button>
    </>
  );
}
