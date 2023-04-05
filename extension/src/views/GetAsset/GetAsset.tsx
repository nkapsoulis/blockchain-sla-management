/* eslint-disable react/jsx-props-no-spreading */

import React from 'react';
import { useNavigate } from 'react-router-dom';

import { JsonViewer } from '@textea/json-viewer';
import Button from '../../components/Button/Button';
import Input from '../../components/Input/Input';
import { Approval, SLA } from '../../models/Asset.model';
import AssetTransferService from '../../services/AssetTransferService';
import Asset from '../../components/Asset/Asset';
import localStorage from '../../utils/localStorage';
import crypto from '../../utils/crypto';

export default function GetAsset() {
  const [asset, setAsset] = React.useState({} as SLA);
  const [approvals, setApprovals] = React.useState({} as Approval);
  const [error, setError] = React.useState('');
  const [showPasswordPrompt, setShowPasswordPrompt] = React.useState(false);
  const navigate = useNavigate();

  const handleApprovals = (assetID: string) => {
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

  const handleGetAsset = (e: any) => {
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

    handleApprovals(assetID);
  };

  const handleApproveAsset = async (e: any) => {
    e.preventDefault();
    const passphrase = e.target.elements.passphrase.value;

    const mnemonicEnc = await localStorage.getLocalStorage('mnemonic');
    const mnemonic = crypto.decrypt(mnemonicEnc, passphrase);

    AssetTransferService.approveAsset(asset.id, mnemonic)
      .then((response) => {
        if (response.success) {
          handleApprovals(asset.id);
        } else {
          setError(error.concat(`${response.message!}\n`));
        }
      });
  };
  return (
    <>
      <h1>Search Asset</h1>
      <form action="GET" onSubmit={handleGetAsset}>
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
            <td><Asset {...asset} /></td>
            <td><JsonViewer value={approvals} rootName="Approved by" defaultInspectDepth={0} /></td>
            <td><Button onClick={() => setShowPasswordPrompt(true)}>Approve</Button></td>
          </tr>
        </tbody>
      </table>
      )}
      { showPasswordPrompt
      && (
      <form action="POST" onSubmit={handleApproveAsset}>
        <Input type="password" id="passphrase" name="Passphrase:" required />
        <Button fullWidth>Confirm</Button>
      </form>
      )}
    </>
  );
}
