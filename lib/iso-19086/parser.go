package iso19086

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
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
	var m Metrics
	err := json.Unmarshal(metricData, &m)
	if err != nil {
		return nil, err
	}
	fmt.Println(m)
	return &m, nil
}

func Parse(slaData, metricData []byte) (bool, error) {
	sla, err := ReadSLA(slaData)
	if err != nil {
		return false, err
	}

	metric, err := ReadMetric(metricData)
	if err != nil {
		return false, err
	}

	log.Println(sla.ID)
	metrics := ParseMetrics(metric)
	params := ParseParameters(sla.SLO.Parameters)
	violated, err := ParseExpression(sla.SLO.Expression.Expression, *params, *metrics)
	if err != nil {
		return false, err
	}

	return violated, nil
}

func ParseParameters(params []Parameter) *map[string]string {

	mapping := make(map[string]string)
	for _, p := range params {
		mapping[p.ReferenceID] = p.Parameter.(string)
	}
	return &mapping
}

func ParseMetrics(metrics *Metrics) *map[string]string {
	mapping := make(map[string]string)
	fields := reflect.VisibleFields(reflect.TypeOf(*metrics))
	value := reflect.ValueOf(metrics).Elem()
	for _, field := range fields {
		if field.Name == "ID" || field.Name == "SLAID" {
			continue
		}
		f := value.FieldByName(field.Name)
		mapping[field.Name] = f.Interface().(string)
	}
	return &mapping
}

func ParseExpression(expr string, params, metric map[string]string) (bool, error) {
	parts := strings.Split(strings.TrimSpace(expr), " ")
	fmt.Println(parts)

	for i, p := range parts {
		param, ok := params[p]
		if ok {
			parts[i] = param
			continue
		}

		param, ok = metric[p]
		if ok {
			parts[i] = param
		}
	}

	p1, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, err
	}

	p2, err := strconv.Atoi(parts[2])
	if err != nil {
		return false, err
	}

	switch parts[1] {
	case "<":
		return !(p1 < p2), nil
	case ">":
		return !(p1 > p2), nil
	default:
		return false, fmt.Errorf("unknown operand %v", parts[1])
	}
}
