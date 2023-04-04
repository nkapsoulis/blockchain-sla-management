import React from 'react';
import { useNavigate } from 'react-router-dom';

import { JsonViewer } from '@textea/json-viewer';
import Button from '../../components/Button/Button';
import Input from '../../components/Input/Input';
import { Approval, SLA } from '../../models/Asset.model';
import AssetTransferService from '../../services/AssetTransferService';

export default function GetAsset() {
  const [asset, setAsset] = React.useState({} as SLA);
  const [approvals, setApprovals] = React.useState({} as Approval);
  const [error, setError] = React.useState('');
  const navigate = useNavigate();

  const handleSubmit = (e: any) => {
    e.preventDefault();
    const assetID = e.target.elements['asset-id'].value;
    AssetTransferService.getAsset(assetID)
      .then((response) => {
        if (response.success) {
          setAsset(response.asset!);
          setError('');
        } else {
          setAsset({} as SLA);
          setError(error.concat(`${response.message!}\n`));
        }
      });
    AssetTransferService.getAssetApprovals(assetID)
      .then((response) => {
        if (response.success) {
          setApprovals(response.approvals!);
          setError('');
        } else {
          setApprovals({} as Approval);
          setError(error.concat(`${response.message!}\n`));
        }
      });
  };
  return (
    <>
      <h1>Search Asset</h1>
      <form action="GET" onSubmit={handleSubmit}>
        <Input type="text" id="asset-id" name="Asset ID:" placeholder="asset2" required />
        <Button fullWidth>Search</Button>
      </form>
      <Button fullWidth onClick={() => navigate('/navigation')}>Back</Button>

      {error !== '' && (
      <div>
        <h2>{`Error: ${error}`}</h2>
      </div>
      )}

      {Object.keys(asset).length > 0 && (
      <table>
        <tbody>
          <tr>
            <th>SLA</th>
            <th>Approvals</th>
          </tr>
          <tr>
            <td><JsonViewer value={asset} rootName={asset.id} defaultInspectDepth={0} /></td>
            <td><JsonViewer value={approvals} rootName="Approved by" defaultInspectDepth={0} /></td>
          </tr>
        </tbody>
      </table>
      )}
    </>
  );
}
