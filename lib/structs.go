package lib

type SLA struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	State      string     `json:"state"`
	Assessment Assessment `json:"assessment"`
	Details    Detail     `json:"details"`
}

type Detail struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Name       string      `json:"Name"`
	Provider   Entity      `json:"provider"`
	Client     Entity      `json:"client"`
	Creation   string      `json:"creation"`
	Guarantees []Guarantee `json:"guarantees"`
	Service    string      `json:"service"`
}

type Assessment struct {
	FirstExecution string `json:"first_execution"`
	LastExecution  string `json:"last_execution"`
}

type Entity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Importance struct {
	Name       string `json:"name"`
	Constraint string `json:"constraint"`
}

type Guarantee struct {
	Name       string       `json:"name"`
	Constraint string       `json:"constraint"`
	Importance []Importance `json:"importance"`
}

type Violation struct {
	ID             string  `json:"id"`
	SLAID          string  `json:"sla_id"`
	GuaranteeID    string  `json:"guarantee_id"`
	Datetime       string  `json:"datetime"`
	Constraint     string  `json:"constraint"`
	Values         []Value `json:"values"`
	ImportanceName string  `json:"importanceName"`
	Importance     int     `json:"importance"`
	AppID          string  `json:"appID"`
}

type Value struct {
	Key      string  `json:"key"`
	Value    float32 `json:"value"`
	Datetime string  `json:"datetime"`
}

type Approval struct {
	ProviderApproved bool `json:"providerApproved"`
	ConsumerApproved bool `json:"consumerApproved"`
}

type User struct {
	Name       string `json:"name"`
	PubKey     string `json:"pubkey"`
	Balance    string `json:"balance"`
	ProviderOf string `json:"providerOf"`
	ClientOf   string `json:"clientOf"`
}

