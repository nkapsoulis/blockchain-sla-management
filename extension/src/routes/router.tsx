import React from 'react';
import {
  createMemoryRouter,
} from 'react-router-dom';

import AssetTransferService from '../services/AssetTransferService';
import { Authentication, action as authAction } from '../views/Authentication/Authentication';
import { CreateAsset, action as createAction } from '../views/CreateAsset/CreateAsset';
import MyAssets from '../views/MyAssets/MyAssets';
import GetAsset from '../views/GetAsset/GetAsset';
import Landing from '../views/Landing/Landing';
import Navigation from '../views/Navigation/Navigation';
import { TransferAsset, action as transferAction } from '../views/TransferAsset/TransferAsset';

const router = createMemoryRouter([
  {
    path: '/',
    element: <Landing />,
  },
  {
    path: '/auth',
    element: <Authentication />,
    action: authAction,
  },
  {
    path: '/navigation',
    element: <Navigation />,
  },
  {
    path: '/create-asset',
    element: <CreateAsset />,
    action: createAction,
  },
  {
    path: '/assets',
    element: <MyAssets />,
    loader: AssetTransferService.getUserAssets,
  },
  {
    path: '/get-asset',
    element: <GetAsset />,
  },
  {
    path: '/transfer',
    element: <TransferAsset />,
    action: transferAction,
  },
]);

export default router;
