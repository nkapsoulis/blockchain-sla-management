interface User {
  name: string,
  pubkey?: string,
  balance?: string,
  clientOf?: string,
  providerOf?: string,
}

export default User;
