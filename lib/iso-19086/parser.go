package iso19086

import (
	"encoding/json"
	"fmt"
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

	metrics := ParseMetrics(metric)
	violated, err := ParseSLO(sla.SLO, metrics)
	if err != nil {
		return false, err
	}

	return violated, nil
}

func ParseMetrics(metrics *Metrics) *map[string]string {
	mapping := make(map[string]string)
	fields := reflect.VisibleFields(reflect.TypeOf(*metrics))
	value := reflect.ValueOf(metrics).Elem()
	for _, field := range fields {
		switch field.Name {
		case "ID":
		case "SLAID":
			continue
		case "Sample":
			sampleFields := reflect.VisibleFields(reflect.TypeOf(metrics.Sample))
			sampleValue := reflect.ValueOf(&metrics.Sample).Elem()
			for _, sf := range sampleFields {
				ff := sampleValue.FieldByName(sf.Name)
				mapping[sf.Name] = ff.Interface().(string)
			}
		default:
			f := value.FieldByName(field.Name)
			mapping[field.Name] = f.Interface().(string)
		}
	}
	return &mapping
}

func ParseParameters(params []Parameter) map[string]string {

	mapping := make(map[string]string)
	for _, p := range params {
		switch p.ReferenceID {
		case "PBH_List":
			mapping[p.ReferenceID] = strconv.Itoa(len(p.Parameters))
		case "MIRespL":
			mapping[p.ReferenceID] = p.Parameter.(string)
		case "SIRL":
			mapping[p.ReferenceID] = p.Parameter.(string)
		}
	}
	return mapping
}

func ParseSLO(slo SLO, metrics *map[string]string) (bool, error) {
	ums, err := ParseUnderlyingMetrics(slo.UnderlyingMetrics, metrics)
	if err != nil {
		return false, err
	}

	params := ParseParameters(slo.Parameters)

	id := strings.Split(slo.ReferenceID, "_")
	switch id[0] {
	case "IRT":
		SIRT, err := strconv.Atoi((*ums)["SIRT"])
		if err != nil {
			return false, err
		}

		SIRL, err := strconv.Atoi(params["SIRL"])
		if err != nil {
			return false, err
		}

		return !(SIRT < SIRL), nil
	case "IRespT":
		MiRespT, err := strconv.Atoi((*ums)["MiRespT"])
		if err != nil {
			return false, err
		}

		MiRespL, err := strconv.Atoi(params["MiRespL"])
		if err != nil {
			return false, err
		}
		return !(MiRespT < MiRespL), err
	}
	return false, fmt.Errorf("you should have never seen this error")
}

func ParseUnderlyingMetrics(ums []UnderlyingMetrics, metrics *map[string]string) (*map[string]string, error) {
	mapping := make(map[string]string)
	for _, um := range ums {

		if underlyingMetricsExists(um) {

			uum, err := ParseUnderlyingMetrics(um.UnderlyingMetrics, metrics)
			if err != nil {
				return nil, err
			}
			mergeMaps(&mapping, uum)
		}
		fmt.Println(um.ReferenceID)
		switch um.ReferenceID {
		case "PBH":
			params := ParseParameters(um.Parameters)
			mapping["PBH"] = params["PBH_List"]

		case "SIRT":
			fmt.Println(metrics)
			incidentResolutionTime, err := strconv.Atoi((*metrics)["IncidentResolutionTime"])
			if err != nil {
				return nil, err
			}
			incidentReportTime, err := strconv.Atoi((*metrics)["IncidentReportTime"])
			if err != nil {
				return nil, err
			}

			PBH, err := strconv.Atoi(mapping["PBH"])
			if err != nil {
				return nil, err
			}

			mapping["SIRT"] = strconv.Itoa(((incidentResolutionTime - incidentReportTime) / 86400) - PBH)
		case "MiRespT":
			incidentResponseTime, err := strconv.Atoi((*metrics)["incident_response_time"])
			if err != nil {
				return nil, err
			}
			incidentReportTime, err := strconv.Atoi((*metrics)["incident_report_time"])
			if err != nil {
				return nil, err
			}

			PBH, err := strconv.Atoi(mapping["PBH"])
			if err != nil {
				return nil, err
			}
			mapping["SIRT"] = strconv.Itoa(((incidentResponseTime - incidentReportTime) / 3600) - 24*PBH)
		}
	}
	return &mapping, nil
}

func underlyingMetricsExists(st UnderlyingMetrics) bool {
	return st.UnderlyingMetrics != nil
}

func mergeMaps(a, b *map[string]string) {
	for k, v := range *b {
		(*a)[k] = v
	}
}
