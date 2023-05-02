package iso19086

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func ReadSLA(slaData []byte) (*SLA, error) {
	v := new(SLA)
	err := json.Unmarshal(slaData, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func ReadMetric(metricData []byte) (*Metrics, error) {
	m := new(Metrics)
	err := json.Unmarshal(metricData, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Parse(sla *SLA, metric *Metrics) {
	log.Println(sla.ID)
	ParseMetrics(*metric)
	params := ParseParameters(sla.SLO.Parameters)
	ParseExpression(sla.SLO.Expression.Expression, *params)

}

func ParseParameters(params []Parameter) *map[string]string {

	mapping := make(map[string]string)
	for _, p := range params {
		mapping[p.ReferenceID] = p.Parameter.(string)
	}
	return &mapping
}

func ParseMetrics(metrics Metrics) {
	// mapping := make(map[string]string)
	fields := reflect.VisibleFields(reflect.TypeOf(metrics))
	value := reflect.ValueOf(metrics)

	fmt.Println(fields)
	fmt.Println(value)

	// for _, field := range fields {
	// 	f := value.FieldByName(field.Name)
	// 	if f.IsValid() {
	// 		mapping[field.Name] = metrics[field.Name]
	// 	}
	// }
}

func ParseExpression(expr string, params map[string]string) {
	parts := strings.Split(strings.TrimSpace(expr), " ")
	fmt.Println(parts)

	for i, p := range parts {
		param, ok := params[p]
		if ok {
			parts[i] = param
		}
	}
	fmt.Println(parts)
}
