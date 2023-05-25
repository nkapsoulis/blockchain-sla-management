import localStorage from './localStorage';

let proto: string;
let hostname: string;
let port: string;

export const getBaseURL = () => {
  const portExpr = port !== '' ? `:${port}` : '';
  return `${proto}//${hostname}${portExpr}`;
};

export const getAssetsURL = () => `${getBaseURL()}/assets`;
export const getAuthURL = () => `${getBaseURL()}/auth`;

export const getTransferAssetURL = () => `${getAssetsURL()}/transfer`;
export const getLoginURL = () => `${getAuthURL()}/login`;
export const getLogoutURL = () => `${getAuthURL()}/logout`;
export const getUserURL = () => `${getAuthURL()}/user`;

export function setBaseUrl(urlString: string): string {
  try {
    // If we don't get an error here, we know it's a
    // valid URL
    const url = new URL(urlString);

    localStorage.setLocalStorage('url', urlString);
    proto = url.protocol;
    hostname = url.hostname;
    port = url.port;

    return '';
  } catch (e: any) {
    if (e instanceof TypeError) {
      return 'Invalid URL';
    }
    return 'Unknown error';
  }
}
