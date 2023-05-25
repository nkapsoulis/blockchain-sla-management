import AssetTransferService from './AssetTransferService';
import authService from './AuthService';

async function getUserData() {
  return Promise.all([AssetTransferService.getUserAssets(), authService.getUser()]);
}

const UserService = {
  getUserData,
};

export default UserService;
