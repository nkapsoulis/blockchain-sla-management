import React from 'react';
import { useLoaderData, useNavigate } from 'react-router-dom';
import { JsonViewer } from '@textea/json-viewer';

import Button from '../../components/Button/Button';
import { SLA } from '../../models/Asset.model';

export default function MyAssets() {
  const [error, setError] = React.useState('');
  const [assets, setAssets] = React.useState<SLA[]>([]);

  const data = useLoaderData() as any;
  React.useEffect(() => {
    if (data.success !== undefined) {
      setError(data.message);
    } else {
      setAssets(data.assets ? data.assets : []);
      setError('');
    }
  }, [data]);

  const navigate = useNavigate();
  return (
    <>
      <h1>Assets</h1>
      {
        error !== ''
        && <p>{`Error: ${error}`}</p>
      }
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
