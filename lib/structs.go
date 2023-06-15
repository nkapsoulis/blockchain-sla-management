package lib

type Entity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Approval struct {
	ProviderApproved bool `json:"providerApproved"`
	ConsumerApproved bool `json:"consumerApproved"`
}

type User struct {
	Name          string `json:"name"`
	PubKey        string `json:"pubkey"`
	Balance       string `json:"balance"`
	NftProviderOf string `json:"nftProviderOf"`
	ProviderOf    string `json:"providerOf"`
	ClientOf      string `json:"clientOf"`
}
