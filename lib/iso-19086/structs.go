package iso19086

type Expression struct {
	Expression string `json:"expression"`
}

type Parameter struct {
	Name        string      `json:"name"`
	ReferenceID string      `json:"referenceId"`
	Unit        string      `json:"unit"`
	Scale       string      `json:"scale,omitempty"`
	Parameter   interface{} `json:"parameter"`
	Parameters  []string    `json:"parameters,omitempty"`
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
	ID       string `json:"id,omitempty"`
	SLO      SLO    `json:"slo"`
	Provider Entity `json:"provider"`
	Client   Entity `json:"client,omitempty"`
}

type SampleData struct {
	IncidentResponseTime   string `json:"incident_response_time"`
	IncidentReportTime     string `json:"incident_report_time"`
	IncidentResolutionTime string `json:"incident_resolution_time"`
}

type Metrics struct {
	ID     string     `json:"id"`
	SLAID  string     `json:"sla_id"`
	Sample SampleData `json:"SAMPLE"`
}
