import { Approval, SLA } from '../models/Asset.model';
import APIResponse from '../models/APIResponse.model';

import { getAssetsURL, getTransferAssetURL } from '../utils/urls';

async function getUserAssets(): Promise<Response | APIResponse> {
  const response = await fetch(getAssetsURL(), {
    method: 'GET',
    credentials: 'include',
  });
  if (response.status !== 200) {
    return { success: false, message: 'Unauthorized' };
  }
  return response.json();
}

async function getAsset(assetID: string): Promise<APIResponse> {
  const response = await fetch(`${getAssetsURL()}/${assetID}`, {
    method: 'GET',
    credentials: 'include',

  });
  if (response.status === 200) {
    const json = await response.json();
    return { success: true, asset: json.data as SLA };
  }
  if (response.status === 404) {
    return { success: false, message: `Asset ${assetID} not found` };
  }
  const json = await response.json();
  return { success: false, message: json };
}

async function getAssetApprovals(assetID: string): Promise<APIResponse> {
  const response = await fetch(`${getAssetsURL()}/${assetID}/approvals`, {
    method: 'GET',
    credentials: 'include',

  });
  if (response.status === 200) {
    const json = await response.json();
    return { success: true, approvals: json.data as Approval };
  }
  if (response.status === 404) {
    return { success: false, message: `Asset ${assetID} not found` };
  }
  const json = await response.json();
  return { success: false, message: json };
}

async function createAsset(asset: SLA): Promise<APIResponse> {
  const response = await fetch(getAssetsURL(), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify(asset),
  });
  if (response.status === 201) {
    return { success: true };
  }
  const json = await response.json();
  return { success: false, message: json };
}

async function approveAsset(assetID: string): Promise<APIResponse> {
  const response = await fetch(`${getAssetsURL()}/${assetID}/approve`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
  });
  if (response.status === 200) {
    return { success: true };
  }
  const json = await response.json();
  return { success: false, message: json };
}

async function transferAsset(assetId: string, newOwner: string): Promise<APIResponse> {
  const response = await fetch(getTransferAssetURL(), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ assetId, newOwner }),
  });
  const json = await response.json();
  if (response.status === 200) {
    return { success: true, message: json };
  }
  return { success: false, message: json };
}

const AssetTransferService = {
  getUserAssets,
  getAsset,
  getAssetApprovals,
  approveAsset,
  createAsset,
  transferAsset,
};

export default AssetTransferService;
