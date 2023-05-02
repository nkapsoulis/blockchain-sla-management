package iso19086

type Expression struct {
	Expression string `json:"expression"`
}

type Parameter struct {
	Name        string      `json:"name"`
	ReferenceID string      `json:"referenceId"`
	Unit        string      `json:"unit"`
	Parameter   interface{} `json:"parameter"`
}

type Rule struct {
	Rule        string `json:"rule"`
	Note        string `json:"note"`
	ReferenceID string `json:"referenceId"`
}

type Sample struct {
	Name        string `json:"name"`
	ReferenceID string `json:"referenceId"`
	Timestamp   string `json:"timestamp"`
	Scale       string `json:"scale"`
	Value       string `json:"value"`
	Protocol    string `json:"protocol"`
	Operation   string `json:"operation"`
	Note        string `json:"note"`
}

type UnderlyingMetrics struct {
	Name              string              `json:"name"`
	ReferenceID       string              `json:"referenceId"`
	Unit              string              `json:"unit"`
	Scale             string              `json:"scale"`
	Expression        Expression          `json:"expression,omitempty"`
	Parameters        []Parameter         `json:"parameters,omitempty"`
	UnderlyingMetrics []UnderlyingMetrics `json:"underlyingMetrics,omitempty"`
	Rules             []Rule              `json:"rules,omitempty"`
	Samples           []Sample            `json:"samples,omitempty"`
}

type Entity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SLO struct {
	Name              string              `json:"name"`
	ReferenceID       string              `json:"referenceId"`
	Scale             string              `json:"scale"`
	Expression        Expression          `json:"expression"`
	Parameters        []Parameter         `json:"parameters"`
	UnderlyingMetrics []UnderlyingMetrics `json:"underlyingMetrics,omitempty"`
}

type SLA struct {
	ID       string `json:"id"`
	SLO      SLO    `json:"slo"`
	Provider Entity `json:"provider"`
	Client   Entity `json:"client"`
	State    string `json:"state"`
}

type Metrics struct {
	ID    string `json:"id"`
	SLAID string `json:"sla_id"`
	SIRT  string `json:"sirl,omitempty"`
}
