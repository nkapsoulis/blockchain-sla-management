import CryptoJS from 'crypto-js';

function encrypt(data: string, passphrase: string): string {
  return CryptoJS.AES.encrypt(
    data,
    passphrase,
  ).toString();
}

function decrypt(encrypted: string, passphrase: string) {
  return CryptoJS.AES.decrypt(encrypted, passphrase);
}

const crypto = {
  encrypt,
  decrypt,
};

export default crypto;
