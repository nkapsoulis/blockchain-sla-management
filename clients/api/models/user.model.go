package models

type User struct {
	Name     string `json:"name,omitempty"`
	PrivKey  string `json:"privkey,omitempty"`
	PubKey   string `json:"pubkey,omitempty"`
	Mnemonic string `json:"mnemonic,omitempty"`
}
