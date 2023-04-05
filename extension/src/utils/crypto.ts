import CryptoJS from 'crypto-js';

function encrypt(data: string, passphrase: string): string {
  return CryptoJS.AES.encrypt(
    data,
    passphrase,
  ).toString();
}

function decrypt(encrypted: string, passphrase: string): string {
  const decrypted = CryptoJS.AES.decrypt(encrypted, passphrase).toString(CryptoJS.enc.Utf8);
  return decrypted;
}

const crypto = {
  encrypt,
  decrypt,
};

export default crypto;
